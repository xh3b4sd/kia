package eks

import (
	"context"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
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

	{
		r.logger.Log(ctx, "level", "info", "message", "deleting istio-asset chart")

		out, err = exec.Command("helm", "delete", "istio-eks", "--namespace", "istio-system").CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	// We use istio gateways and external-dns to register DNS records in
	// Route53. In order to cleanup the cluster specific DNS records we need to
	// delete the istio-eks chart first and let external-dns take care of the
	// cleanup procedure. For now the mechanism is purely time based, which
	// means we just wait for 5 minutes. This implies the cleanup might fail and
	// we proceed deleting the cluster regardless, leaving behind Route53 DNS
	// records.
	{
		r.logger.Log(ctx, "level", "info", "message", "waiting for cleanup")

		time.Sleep(5 * time.Minute)
	}

	{
		r.logger.Log(ctx, "level", "info", "message", "deleting cert-asset chart")

		out, err = exec.Command("helm", "delete", "cert-asset", "--namespace", "cert-manager").CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	{
		r.logger.Log(ctx, "level", "info", "message", "deleting eks cluster")

		out, err = exec.Command("eksctl", "delete", "cluster", "--name", r.flag.Cluster).CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	return nil
}
