package eks

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/xh3b4sd/tracer"

	"github.com/xh3b4sd/kia/pkg/config"
	"github.com/xh3b4sd/kia/pkg/env"
)

type flag struct {
	Cluster string
	KiaPath string
	Region  string
	SecPath string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.Cluster, "cluster", "c", "", "Cluster ID used for AWS and EKS resource naming.")
	cmd.Flags().StringVarP(&f.KiaPath, "kia", "k", config.GetKia(os.Getenv(env.KiaBasePath)), "Kia base path on the local file system.")
	cmd.Flags().StringVarP(&f.Region, "region", "r", "eu-central-1", "Region in which the EKS cluster gets created.")
	cmd.Flags().StringVarP(&f.SecPath, "sec", "s", config.GetSec(os.Getenv(env.SecBasePath)), "Sec base path on the local file system.")
}

func (f *flag) Validate() error {
	if f.Cluster == "" {
		return tracer.Maskf(invalidFlagError, "-c/--cluster must not be empty")
	}

	if f.KiaPath == "" {
		return tracer.Maskf(invalidFlagError, "-k/--kia must not be empty")
	}

	if f.Region == "" {
		return tracer.Maskf(invalidFlagError, "-r/--region must not be empty")
	}

	if f.SecPath == "" {
		return tracer.Maskf(invalidFlagError, "-s/--sec must not be empty")
	}

	return nil
}
