package osx

import (
	"context"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/xh3b4sd/logger"
	"github.com/xh3b4sd/tracer"
)

type runner struct {
	logger logger.Interface
}

func (r *runner) Run(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	err := r.run(ctx, cmd, args)
	if err != nil {
		return tracer.Mask(err)
	}

	return nil
}

func (r *runner) run(ctx context.Context, cmd *cobra.Command, args []string) error {
	var err error
	var out []byte

	{
		r.logger.Log(ctx, "level", "info", "message", "deleting kind cluster")

		out, err = exec.Command("kind", "delete", "cluster").CombinedOutput()
		if err != nil {
			return tracer.Maskf(executionFailedError, "%s", out)
		}
	}

	return nil
}
