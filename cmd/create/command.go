package create

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

const (
	name  = "create"
	short = "Generate a grpc workflow for e.g. golang code generation."
	long  = `Generate a grpc workflow for e.g. golang code generation. The workflow
generated here works in a setup of two Github repositories. Call them
apischema and gocode. The workflow generated with the following command is
added to the apischema repository.

    workflow generate grpc -o xh3b4sd -r gocode

In order to make the workflow function correctly a deploy key is generated
and distributed as follows. The public and private key files are added to the
apischema repository. The public key is added as deploy key with write access
to the gocode repository. Deploy keys and GPG password are generated with the
following command.

    red generate keys -d .github/asset/xh3b4sd/gocode

Generating the deploy keys also generates a GPG password which is used to
decrypt the encrypted private key within the build container of the workflow.
The GPG password needs to be added to the apischema secrets using the
following name.

    RED_GPG_PASS_XH3B4SD_GOCODE

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
