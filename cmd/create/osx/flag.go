package osx

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/kia/pkg/config"
	"github.com/xh3b4sd/kia/pkg/env"
)

type flag struct {
	Kia string
	Sec string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.Kia, "kia", "k", config.GetKia(os.Getenv(env.KiaBasePath)), "Kia base path on the local file system.")
	cmd.Flags().StringVarP(&f.Sec, "sec", "s", config.GetSec(os.Getenv(env.SecBasePath)), "Sec base path on the local file system.")
}

func (f *flag) Validate() error {
	if f.Kia == "" {
		return tracer.Maskf(invalidFlagError, "-k/--kia must not be empty")
	}

	if f.Sec == "" {
		return tracer.Maskf(invalidFlagError, "-s/--sec must not be empty")
	}

	return nil
}
