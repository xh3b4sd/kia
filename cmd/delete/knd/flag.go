package knd

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/tracer"
)

type flag struct {
	Cluster string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.Cluster, "cluster", "c", "kind", "Cluster ID of the Kind cluster.")
}

func (f *flag) Validate() error {
	if f.Cluster == "" {
		return tracer.Maskf(invalidFlagError, "-c/--cluster must not be empty")
	}

	return nil
}
