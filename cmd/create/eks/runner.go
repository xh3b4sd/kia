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
	var err error
	var out []byte

	secrets := map[string]string{}
	{
		r.logger.Log(ctx, "level", "info", "message", "decrypting local secrets")

		out, err = exec.Command("red", "decrypt", "-i", mustAbs(r.flag.Sec), "-o", "-", "-s").CombinedOutput()
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

		out, err = exec.Command("eksctl", "create", "cluster", "--config-file", mustAbs(r.flag.Kia, "env/eks/eksctl.yaml")).CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	{
		r.logger.Log(ctx, "level", "info", "message", "installing service mesh")

		out, err = exec.Command("istioctl", "install", "-f", mustAbs(r.flag.Kia, "env/eks/istio.yaml")).CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	{
		r.logger.Log(ctx, "level", "info", "message", "creating infra namespace")

		out, err = exec.Command("kubectl", "create", "namespace", "infra").CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	{
		r.logger.Log(ctx, "level", "info", "message", "configure istio injection")

		out, err = exec.Command("kubectl", "label", "namespace", "infra", "istio-injection=enabled").CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	{
		r.logger.Log(ctx, "level", "info", "message", "installing infra chart")

		out, err = exec.Command(
			"helm",
			"install",
			"infra",
			mustAbs(r.flag.Kia, "env/def/infra/"),
			"--namespace", "infra",
			"--set", "dockerconfigjson="+mustDockerAuth(secrets),
		).CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	// TODO create Hosted Zone in Route53.
	//
	//     https://github.com/kubernetes-sigs/external-dns/blob/master/docs/tutorials/aws.md#set-up-a-hosted-zone
	//
	// Where to find chart documentation.
	//
	//     https://artifacthub.io/packages/helm/bitnami/external-dns
	//
	{
		r.logger.Log(ctx, "level", "info", "message", "installing external-dns chart")

		out, err = exec.Command("kubectl", "create", "namespace", "external-dns").CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}

		out, err = exec.Command("helm", "repo", "add", "bitnami", "https://charts.bitnami.com/bitnami").CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}

		out, err = exec.Command("helm", "repo", "update").CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}

		out, err = exec.Command(
			"helm",
			"install",
			"external-dns",
			"bitnami/external-dns",
			"--namespace", "external-dns",
			"--version", "v3.4.9",
			"--set", "aws.credentials.accessKey="+secrets["aws.accessid"],
			"--set", "aws.credentials.secretKey="+secrets["aws.secretid"],
			"--set", "aws.region=eu-central-1",
			"--set", "domainFilters={aws.venturemark.co}",
			"--set", "provider=aws",
			"--set", "sources={istio-gateway}",
		).CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	{
		r.logger.Log(ctx, "level", "info", "message", "installing cert-manager chart")

		out, err = exec.Command("kubectl", "create", "namespace", "cert-manager").CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}

		out, err = exec.Command("helm", "repo", "add", "jetstack", "https://charts.jetstack.io").CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}

		out, err = exec.Command("helm", "repo", "update").CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}

		out, err = exec.Command(
			"helm",
			"install",
			"cert-manager",
			"jetstack/cert-manager",
			"--namespace", "cert-manager",
			"--version", "v1.0.3",
			"--set", "installCRDs=true",
		).CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	{
		r.logger.Log(ctx, "level", "info", "message", "installing cert-asset chart")

		out, err = exec.Command(
			"helm",
			"install",
			"cert-asset",
			mustAbs(r.flag.Kia, "env/eks/cert-asset/"),
			"--namespace", "cert-manager",
			"--set", "aws.accessid="+secrets["aws.accessid"],
			"--set", "aws.secretid="+secrets["aws.secretid"],
		).CombinedOutput()
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

func mustDockerAuth(s map[string]string) string {
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
