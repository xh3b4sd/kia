package knd

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

const (
	name  = "knd"
	short = "Create kubernetes infrastructure environments for kind."
	long  = `Create kubernetes infrastructure environments for kind. The basis for this
type of environment is a local kind cluster. Kind stands for kubernetes in
docker. For more information check the kind repository.

    https://github.com/kubernetes-sigs/kind

In order to create and setup the kind cluster we need to properly configure
the kia command line tool. This is done via its config file, tracked on the
local file system. The kia base path must be set. This is the local path of
the kia repository from which general templates are read. Add the following
line to your config file according to your local setup.

    kia: "~/project/xh3b4sd/kia/"

The sec base path must be set. This is the local path of the sec repository
from which secret data is read. Add the following line to your config file
according to your local setup.

    sec: "~/project/xh3b4sd/sec/"
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
