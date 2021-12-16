package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/coloring"
	"github.com/nyaosorg/go-readline-ny/simplehistory"
	"github.com/zetamatta/go-shellcommand"
)

func mains() error {
	history := simplehistory.New()

	editor := readline.Editor{
		Prompt:   func() (int, error) { return fmt.Print("$ ") },
		Writer:   colorable.NewColorableStdout(),
		History:  history,
		Coloring: &coloring.VimBatch{},
	}
	for {
		text, err := editor.ReadLine(context.Background())
		if err != nil {
			if err == readline.CtrlC {
				continue
			}
			return err
		}

		process, err := shellcommand.System(text)
		if err != nil {
			return err
		}
		process.Wait()

		history.Add(text)
	}
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
