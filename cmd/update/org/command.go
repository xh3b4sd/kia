package org

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

const (
	name  = "org"
	short = "Update the current organization name."
	long  = `Update the current organization name. Managing clusters for different
organizations requires kia to know where its own assets are and where to find
secret data. The latter is managed via the red command line tool. See
https://github.com/xh3b4sd/red for more information. Below is shown the
expected config file location on your file system, including the required
structure and its associated values.

    $ cat ~/.config/kia/config.yaml
    kia: "~/projects/xh3b4sd/kia"
    org:
      list:
        - org: "xh3b4sd"
          sec: "~/projects/xh3b4sd/sec"
        - org: "yourorg"
          sec: "~/projects/yourorg/sec"
      selected: "xh3b4sd"

Given the example config file above the organization used by kia can be
changed as shown below.

    kia update org --selected yourorg
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
