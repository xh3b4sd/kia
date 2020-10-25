package eks

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/tracer"
)

type flag struct {
	Cluster string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.Cluster, "cluster", "c", "", "Cluster ID used for AWS and EKS resource naming.")
}

func (f *flag) Validate() error {
	if f.Cluster == "" {
		return tracer.Maskf(invalidFlagError, "-c/--cluster must not be empty")
	}

	return nil
}
