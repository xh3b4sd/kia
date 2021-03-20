package knd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/kia/pkg/config"
	"github.com/xh3b4sd/kia/pkg/env"
	"github.com/xh3b4sd/kia/pkg/file"
)

type flag struct {
	Cluster string
	Image   string
	KiaPath string
	SecPath string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.Cluster, "cluster", "c", "kind", "Cluster ID of the Kind cluster.")
	cmd.Flags().StringVarP(&f.Image, "image", "i", "", "Kind cluster image to use.")
	cmd.Flags().StringVarP(&f.KiaPath, "kia", "k", config.GetKia(os.Getenv(env.KiaBasePath)), "Kia base path on the local file system.")
	cmd.Flags().StringVarP(&f.SecPath, "sec", "s", config.GetSec(os.Getenv(env.SecBasePath)), "Sec base path on the local file system.")
}

func (f *flag) Validate() error {
	{
		if f.Cluster == "" {
			return tracer.Maskf(invalidFlagError, "-c/--cluster must not be empty")
		}
	}

	{
		if f.KiaPath == "" {
			return tracer.Maskf(invalidFlagError, "-k/--kia must not be empty")
		}

		if !file.Exists(f.KiaPath) {
			return tracer.Maskf(invalidFlagError, "-k/--kia path does not exist")
		}
	}

	{
		if f.SecPath == "" {
			return tracer.Maskf(invalidFlagError, "-s/--sec must not be empty")
		}

		if !file.Exists(f.SecPath) {
			return tracer.Maskf(invalidFlagError, "-s/--sec path does not exist")
		}
	}

	return nil
}
