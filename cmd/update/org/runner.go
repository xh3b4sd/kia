package org

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/xh3b4sd/kia/pkg/config"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type runner struct {
	flag   *flag
	logger logger.Interface
}

func (r *runner) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := r.flag.Validate()
	if err != nil {
		return tracer.Mask(err)
	}

	err = r.run(ctx, cmd, args)
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (r *runner) run(ctx context.Context, cmd *cobra.Command, args []string) error {
	var err error

	c := config.Select(r.flag.Selected)

	err = config.Validate(c)
	if err != nil {
		return tracer.Mask(err)
	}

	err = config.Write(c)
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}
