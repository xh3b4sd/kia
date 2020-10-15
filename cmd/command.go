package cmd

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/kia/cmd/completion"
	"github.com/xh3b4sd/kia/cmd/create"
	"github.com/xh3b4sd/kia/cmd/version"
	"github.com/xh3b4sd/kia/pkg/project"
)

var (
	name        = project.Name()
	description = project.Description()
)

type Config struct {
	Logger logger.Interface
}

func New(config Config) (*cobra.Command, error) {
	if config.Logger == nil {
		return nil, tracer.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	var err error

	var completionCmd *cobra.Command
	{
		c := completion.Config{
			Logger: config.Logger,
		}

		completionCmd, err = completion.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var createCmd *cobra.Command
	{
		c := create.Config{
			Logger: config.Logger,
		}

		createCmd, err = create.New(c)
		if err != nil {
			return nil, tracer.Mask(err)
		}
	}

	var versionCmd *cobra.Command
	{
		c := version.Config{
			Logger: config.Logger,
		}

		versionCmd, err = version.New(c)
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
			Short: description,
			Long:  description,
			RunE:  r.Run,
			// We slience errors because we do not want to see spf13/cobra printing.
			// The errors returned by the commands will be propagated to the main.go
			// anyway, where we have custom error printing for the command line
			// tool.
			SilenceErrors: true,
			SilenceUsage:  true,
		}

		c.AddCommand(completionCmd)
		c.AddCommand(createCmd)
		c.AddCommand(versionCmd)
	}

	return c, nil
}
