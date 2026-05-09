package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/coloring"
	"github.com/nyaosorg/go-readline-ny/completion"
	"github.com/nyaosorg/go-readline-ny/keys"
	"github.com/nyaosorg/go-readline-ny/simplehistory"

	"github.com/hymkor/go-shellcommand"
)

const spaces = "\t\n\v\f\r "

func cutField(source string) (first, rest string) {
	source = strings.TrimLeft(source, spaces)
	var quote bool
	for i, r := range source {
		if r == '"' {
			quote = !quote
		}
		if !quote && strings.IndexRune(spaces, r) >= 0 {
			return source[:i], source[i:]
		}
	}
	return source, ""
}

func mains() error {
	history := simplehistory.New()

	editor := &readline.Editor{
		PromptWriter: func(w io.Writer) (int, error) {
			return io.WriteString(w, "\x1B[0;1;32m$ \x1B[0m")
		},
		Writer:   colorable.NewColorableStdout(),
		History:  history,
		Coloring: &coloring.VimBatch{},
	}
	editor.BindKey(keys.CtrlI, &completion.CmdCompletionOrList{
		Completion: completion.File{},
	})
	for {
		text, err := editor.ReadLine(context.Background())
		if err != nil {
			if err == readline.CtrlC {
				continue
			}
			return err
		}
		cmd, arg := cutField(text)
		cmd = filepath.FromSlash(cmd)
		text = cmd + arg

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
