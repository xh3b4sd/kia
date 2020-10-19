package create

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/kia/cmd/create/eks"
	"github.com/xh3b4sd/kia/cmd/create/osx"
)

const (
	name  = "create"
	short = "Create kubernetes infrastructure environments for e.g. eks and osx."
	long  = "Create kubernetes infrastructure environments for e.g. eks and osx."
)

type Config struct {
	Logger logger.Interface
}

func New(config Config) (*cobra.Command, error) {
	if config.Logger == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	var err error

	var eksCmd *cobra.Command
	{
		c := eks.Config{
			Logger: config.Logger,
		}

		eksCmd, err = eks.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var osxCmd *cobra.Command
	{
		c := osx.Config{
			Logger: config.Logger,
		}

		osxCmd, err = osx.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var c *cobra.Command
	{
		r := &runner{
			logger: config.Logger,
		}

		c = &cobra.Command{
			Use:   name,
			Short: short,
			Long:  long,
			RunE:  r.Run,
		}

		c.AddCommand(eksCmd)
		c.AddCommand(osxCmd)
	}

	return c, nil
}
