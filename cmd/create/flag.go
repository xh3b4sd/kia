package create

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/tracer"
)

type flag struct {
	Github struct {
		Organization string
		Repository   string
	}
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.Github.Organization, "github-organization", "o", "", "Github organization to generate code for.")
	cmd.Flags().StringVarP(&f.Github.Repository, "github-repository", "r", "", "Github repository to generate code for.")
}

func (f *flag) Validate() error {
	{
		if f.Github.Organization == "" {
			return tracer.Maskf(invalidFlagError, "-o/--github-organization must not be empty")
		}
		if f.Github.Repository == "" {
			return tracer.Maskf(invalidFlagError, "-r/--github-repository must not be empty")
		}
	}

	return nil
}
