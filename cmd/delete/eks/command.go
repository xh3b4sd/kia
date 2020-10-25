package eks

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

const (
	name  = "eks"
	short = "Delete kubernetes infrastructure environments for eks."
	long  = `Delete kubernetes infrastructure environments for eks. The deletion process
is mostly straight forward since eks takes care of most of the cloud provider
resources managed in aws. For now there is only one caveat to be aware of. We
use istio gateways and external-dns to register DNS records in Route53. In
order to cleanup the cluster specific DNS records we need to delete the
istio-asset chart first and let external-dns take care of the cleanup
procedure. For now the mechanism is purely time based, which means we just
wait for 5 minutes. This implies the cleanup might fail and we proceed
deleting the cluster regardless, leaving behind Route53 DNS records.
`
)

type Config struct {
	Logger logger.Interface
}

func New(config Config) (*cobra.Command, error) {
	if config.Logger == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	var c *cobra.Command
	{
		f := &flag{}

		r := &runner{
			flag:   f,
			logger: config.Logger,
		}

		c = &cobra.Command{
			Use:   name,
			Short: short,
			Long:  long,
			RunE:  r.Run,
		}

		f.Init(c)
	}

	return c, nil
}
