package org

import (
	"github.com/spf13/cobra"
	"github.com/xh3b4sd/tracer"
)

type flag struct {
	Selected string
}

func (f *flag) Init(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&f.Selected, "selected", "", "", "Select the given organization for current use.")
}

func (f *flag) Validate() error {
	if f.Selected == "" {
		return tracer.Maskf(invalidFlagError, "--selected must not be empty")
	}

	return nil
}
