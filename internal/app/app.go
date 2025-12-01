package app

import (
	"github.com/andreyxaxa/pkg/grep"
)

// Run init and run all components
func Run() error {
	p := grep.NewParams()

	if err := p.Start(); err != nil {
		return err
	}

	return nil
}
