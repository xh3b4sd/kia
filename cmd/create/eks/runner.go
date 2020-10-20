package eks

import (
	"context"
	"encoding/json"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/kia/pkg/docker"
)

type runner struct {
	flag   *flag
	logger logger.Interface
}

func (r *runner) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := r.flag.Validate()
	if err != nil {
		return tracer.Mask(err)
	}

	err = r.run(ctx, cmd, args)
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (r *runner) run(ctx context.Context, cmd *cobra.Command, args []string) error {
	secrets := map[string]string{}
	{
		r.logger.Log(ctx, "level", "info", "message", "decrypting local secrets")

		out, err := exec.Command("red", "decrypt", "-i", mustAbs(r.flag.Sec), "-o", "-", "-s").CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}

		err = json.Unmarshal(out, &secrets)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	{
		r.logger.Log(ctx, "level", "info", "message", "creating eks cluster")

		out, err := exec.Command("eksctl", "create", "cluster", "--config-file", mustAbs(r.flag.Kia, "env/eks/eksctl.yaml")).CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	{
		r.logger.Log(ctx, "level", "info", "message", "installing service mesh")

		out, err := exec.Command("istioctl", "install", "-f", mustAbs(r.flag.Kia, "env/osx/istio.yaml")).CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	{
		r.logger.Log(ctx, "level", "info", "message", "creating infra namespace")

		out, err := exec.Command("kubectl", "create", "namespace", "infra").CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	{
		r.logger.Log(ctx, "level", "info", "message", "configure istio injection")

		out, err := exec.Command("kubectl", "label", "namespace", "infra", "istio-injection=enabled").CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	{
		r.logger.Log(ctx, "level", "info", "message", "installing infra chart")

		out, err := exec.Command("helm", "-n", "infra", "install", "infra", mustAbs(r.flag.Kia, "env/def/infra/"), "--set", "dockerconfigjson="+mustAuth(secrets)).CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	return nil
}

func mustAbs(p ...string) string {
	u, err := user.Current()
	if err != nil {
		panic(err)
	}

	return filepath.Join(u.HomeDir, strings.TrimPrefix(filepath.Join(p...), "~/"))
}

func mustAuth(s map[string]string) string {
	var err error

	var a docker.AuthEncoder
	{
		c := docker.AuthConfig{
			Secrets: s,
		}

		a, err = docker.NewAuth(c)
		if err != nil {
			panic(err)
		}
	}

	enc, err := a.Encode()
	if err != nil {
		panic(err)
	}

	return enc
}
