package delete

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/kia/cmd/delete/eks"
	"github.com/xh3b4sd/kia/cmd/delete/knd"
)

const (
	name  = "delete"
	short = "Delete kubernetes infrastructure environments for e.g. eks and knd."
	long  = "Delete kubernetes infrastructure environments for e.g. eks and knd."
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

	var kndCmd *cobra.Command
	{
		c := knd.Config{
			Logger: config.Logger,
		}

		kndCmd, err = knd.New(c)
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
		c.AddCommand(kndCmd)
	}

	return c, nil
}
