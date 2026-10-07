package main

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimulationRun(t *testing.T) {
	oldArgs := os.Args
	oldCommandLine := flag.CommandLine
	defer func() {
		os.Args = oldArgs
		flag.CommandLine = oldCommandLine
	}()

	t.Run("with split AIs", func(t *testing.T) {
		os.Args = []string{"cmd", "-ai=fato:fato", "-matches=1", "-format=none", "-logging=false"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		err := run()
		assert.NoError(t, err)
	})

	t.Run("with single AI name", func(t *testing.T) {
		os.Args = []string{"cmd", "-ai=fato", "-matches=1", "-format=md", "-logging=false"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		err := run()
		assert.NoError(t, err)
	})

	t.Run("with board setup file", func(t *testing.T) {
		os.Args = []string{"cmd", "-ai=fato:fato", "-matches=1", "-format=none", "-logging=false", "-setup=../../internal/ai/board-setups.json"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		err := run()
		assert.NoError(t, err)
	})

	t.Run("with split setup files", func(t *testing.T) {
		os.Args = []string{"cmd", "-ai=fato:fato", "-matches=1", "-format=none", "-logging=false", "-setup1=../../internal/ai/board-setups.json", "-setup2=../../internal/ai/board-setups.json"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		err := run()
		assert.NoError(t, err)
	})

	t.Run("with non-existent setup file", func(t *testing.T) {
		os.Args = []string{"cmd", "-ai=fato:fato", "-matches=1", "-format=none", "-logging=false", "-setup=nonexistent.json"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		err := run()
		assert.Error(t, err)
	})

	t.Run("main function execution", func(_ *testing.T) {
		os.Args = []string{"cmd", "-ai=fato:fato", "-matches=1", "-format=none", "-logging=false"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		main()
	})
}
