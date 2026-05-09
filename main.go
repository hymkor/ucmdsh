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
	"github.com/nyaosorg/go-readline-ny/completion"
	"github.com/nyaosorg/go-readline-ny/keys"
	"github.com/nyaosorg/go-readline-ny/simplehistory"

	"github.com/hymkor/go-multiline-ny"

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
	var editor multiline.Editor
	history := simplehistory.New()

	editor.SetWriter(colorable.NewColorableStdout())
	editor.SetHistory(history)
	editor.SetPrompt(func(w io.Writer, lnum int) (int, error) {
		if lnum > 0 {
			return fmt.Fprintf(w, "\x1B[0;1;32m%d> \x1B[0m", lnum)
		}
		return io.WriteString(w, "\x1B[0;1;32m$ \x1B[0m")
	},
	)
	editor.BindKey(keys.CtrlI, &completion.CmdCompletionOrList{
		Completion: completion.File{},
	})
	editor.SubmitOnEnterWhen(func(lines []string, _ int) bool {
		quote := false
		cont := false
		for _, line := range lines {
			backslash := false
			for _, c := range line {
				if !backslash {
					if c == '\\' {
						backslash = true
					} else {
						backslash = false
						if c == '"' {
							quote = !quote
						}
					}
				} else {
					backslash = false
				}
			}
			cont = backslash
		}
		return !cont && !quote
	})
	for {
		lines, err := editor.Read(context.Background())
		if err != nil {
			if err == readline.CtrlC {
				continue
			}
			return err
		}
		text4history := strings.Join(lines, "\n")
		text := strings.ReplaceAll(text4history, "\\\n", "")
		cmd, arg := cutField(text)
		cmd = filepath.FromSlash(cmd)
		text = cmd + arg

		process, err := shellcommand.System(text)
		if err != nil {
			return err
		}
		process.Wait()

		history.Add(text4history)
	}
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
