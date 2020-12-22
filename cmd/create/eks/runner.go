package eks

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/xh3b4sd/budget"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/kia/pkg/docker"
)

type runner struct {
	flag   *flag
	path   *path
	logger logger.Interface
}

func (r *runner) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := r.flag.Validate()
	if err != nil {
		return tracer.Mask(err)
	}

	err = r.path.Validate()
	if err != nil {
		return tracer.Mask(err)
	}

	err = r.run(ctx, cmd, args)
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (r *runner) eksData() interface{} {
	type AWS struct {
		Region string
	}

	type Cluster struct {
		Name string
	}

	type Data struct {
		AWS     AWS
		Cluster Cluster
	}

	return Data{
		AWS: AWS{
			Region: r.flag.Region,
		},
		Cluster: Cluster{
			Name: r.flag.Cluster,
		},
	}
}

func (r *runner) run(ctx context.Context, cmd *cobra.Command, args []string) error {
	var err error
	var out []byte

	var b budget.Interface
	{
		c := budget.ConstantConfig{
			Budget:   10,
			Duration: 10 * time.Second,
		}

		b, err = budget.NewConstant(c)
		if err != nil {
			panic(err)
		}
	}

	secrets := map[string]string{}
	{
		r.logger.Log(ctx, "level", "info", "message", "decrypting local secrets")

		out, err = exec.Command("red", "decrypt", "-i", mustAbs(r.flag.SecPath), "-o", "-", "-s").CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}

		err = json.Unmarshal(out, &secrets)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	var eksYaml string
	{
		p := "env/eks/eksctl.yaml"

		f, err := ioutil.ReadFile(mustAbs(r.flag.KiaPath, p))
		if err != nil {
			return tracer.Mask(err)
		}

		t, err := template.New(p).Parse(string(f))
		if err != nil {
			return tracer.Mask(err)
		}

		var b bytes.Buffer
		err = t.ExecuteTemplate(&b, p, r.eksData())
		if err != nil {
			return tracer.Mask(err)
		}

		eksYaml = b.String()
	}

	{
		r.logger.Log(ctx, "level", "info", "message", "creating eks cluster")

		cmd := exec.Command("eksctl", "create", "cluster", "-f", "-")

		in, err := cmd.StdinPipe()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}

		_, err = io.WriteString(in, eksYaml)
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
		in.Close()

		out, err := cmd.CombinedOutput()

		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	{
		r.logger.Log(ctx, "level", "info", "message", "installing service mesh")

		out, err = exec.Command("istioctl", "install", "-f", mustAbs(r.flag.KiaPath, "env/eks/istio.yaml")).CombinedOutput()
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
			mustAbs(r.flag.KiaPath, "env/def/infra/"),
			"--namespace", "infra",
			"--set", "dockerconfigjson="+mustDockerAuth(secrets),
		).CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	// Note that the following code requires a created Hosted Zone in Route53.
	// This should usually be done during the process of domain registration and
	// AWS account creation. A registrar might be Namecheap. A registered domain
	// might be managed in Cloudflare. Some subdomain can then be delegated from
	// Cloudflare to Route53 as such that NS records of AWS nameservers are
	// configured in Cloudflare. For more information about external-dns chart
	// we use below check the following resource.
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
			"--set", "aws.region="+r.flag.Region,
			"--set", "domainFilters={"+secrets["aws.hostedzone"]+"}",
			"--set", "policy=sync",
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

		var action string
		var budget int

		o := func() error {
			defer func() { budget++ }()

			if budget == 0 {
				action = "install"
			} else {
				action = "upgrade"
			}

			out, err = exec.Command(
				"helm",
				action,
				"cert-asset",
				mustAbs(r.flag.KiaPath, "env/eks/cert-asset/"),
				"--namespace", "cert-manager",
				"--set", "aws.accessid="+secrets["aws.accessid"],
				"--set", "aws.region="+r.flag.Region,
				"--set", "aws.secretid="+secrets["aws.secretid"],
			).CombinedOutput()
			if err != nil {
				return tracer.Maskf(executionFailedError, "%s", out)
			}

			return nil
		}

		err := b.Execute(o)
		if err != nil {
			return tracer.Mask(err)
		}
	}

	// Istio assets like gateways and their corresponding certificates have to
	// be installed once cert-manager is properly setup.
	{
		r.logger.Log(ctx, "level", "info", "message", "installing istio-asset chart")

		out, err = exec.Command(
			"helm",
			"install",
			"istio-asset",
			mustAbs(r.flag.KiaPath, "env/eks/istio-asset/"),
			"--namespace", "istio-system",
			"--set", "cluster.name="+r.flag.Cluster,
			"--set", "cluster.zone="+secrets["aws.hostedzone"],
		).CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	{
		r.logger.Log(ctx, "level", "info", "message", "installing flux toolkit")

		os.Setenv("GITHUB_TOKEN", secrets["github.flux.token"])

		out, err = exec.Command(
			"flux",
			"bootstrap",
			"github",
			"--owner",
			"venturemark",
			"--repository",
			"flux",
			"--token-auth",
			"true",
		).CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}

		os.Unsetenv("GITHUB_TOKEN")
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
