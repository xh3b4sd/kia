package knd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/kia/pkg/config"
	"github.com/xh3b4sd/kia/pkg/env"
)

type flag struct {
	KiaPath string
	SecPath string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.KiaPath, "kia", "k", config.GetKia(os.Getenv(env.KiaBasePath)), "Kia base path on the local file system.")
	cmd.Flags().StringVarP(&f.SecPath, "sec", "s", config.GetSec(os.Getenv(env.SecBasePath)), "Sec base path on the local file system.")
}

func (f *flag) Validate() error {
	if f.KiaPath == "" {
		return tracer.Maskf(invalidFlagError, "-k/--kia must not be empty")
	}

	if f.SecPath == "" {
		return tracer.Maskf(invalidFlagError, "-s/--sec must not be empty")
	}

	return nil
}
