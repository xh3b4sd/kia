package knd

import (
	"os/exec"
	"strings"

	"github.com/xh3b4sd/tracer"
)

type path struct {
	Binary []string
}

func (p *path) Validate() error {
	var m []string

	for _, b := range p.Binary {
		_, err := exec.LookPath(b)
		if err != nil {
			m = append(m, b)
		}
	}

	if len(m) != 0 {
		return tracer.Maskf(binaryNotFoundError, strings.Join(m, ", "))
	}

	return nil
}
