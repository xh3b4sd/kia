package eks

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

const (
	name  = "eks"
	short = "Create kubernetes infrastructure environments for eks."
	long  = `Create kubernetes infrastructure environments for eks. The basis for this
type of environment is a remote eks cluster. EKS stands for elastic
kubernetes service. For more information check the eks website.

    https://aws.amazon.com/eks

In order to create and setup the eks cluster we need to properly configure
the kia command line tool. This is done via its config file, tracked on the
local file system. The kia base path must be set. This is the local path of
the kia repository from which general templates are read. Add the following
line to your config file according to your local setup.

    kia: "~/project/xh3b4sd/kia/"

In order to create and setup the eks cluster we need to properly configure
secret data. This is done via a separate private repository containing the
secret data and the red command line tool.

    https://github.com/xh3b4sd/red

The sec base path must be set. This is the local path of the sec repository
from which the red command line tool reads the secret data. Add the following
line to your config file according to your local setup.

    sec: "~/project/xh3b4sd/sec/"

An eks cluster can be created like shown below once the kia config file and
the red command line tool are in place. Cluster creation requires a unique
cluster name. A simple convention could to use the kia prefix and a two digit
number.

    $ kia create eks -c kia02
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:88", "level":"info", "message":"decrypting local secrets", "time":"2020-10-24 20:42:24" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:125", "level":"info", "message":"creating eks cluster", "time":"2020-10-24 20:42:24" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:145", "level":"info", "message":"installing service mesh", "time":"2020-10-24 21:07:52" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:154", "level":"info", "message":"creating infra namespace", "time":"2020-10-24 21:08:30" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:163", "level":"info", "message":"configure istio injection", "time":"2020-10-24 21:08:32" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:172", "level":"info", "message":"installing infra chart", "time":"2020-10-24 21:08:32" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:198", "level":"info", "message":"installing external-dns chart", "time":"2020-10-24 21:08:35" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:235", "level":"info", "message":"installing cert-manager chart", "time":"2020-10-24 21:08:43" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:267", "level":"info", "message":"installing cert-asset chart", "time":"2020-10-24 21:08:54" }
    { "caller":"github.com/xh3b4sd/kia/cmd/create/eks/runner.go:307", "level":"info", "message":"installing istio-asset chart", "time":"2020-10-24 21:09:09" }

After some time the cluster created as shown above would be available
depending on the Route53 hosted zone configured in the secret data
repository. Some api server deployed in the created cluster would be
available like shown below.

    apiserver.kia02.aws.example.com
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
