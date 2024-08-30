package cmd

import (
	"github.com/morpheuszero/gcloc/gcloc/internal/gcloc"
)

func GlocCommandExecute() {
	impl := gcloc.NewGClocImpl()
	impl.Run()
}
