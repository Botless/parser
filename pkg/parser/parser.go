package parser

import (
	"strings"

	"github.com/botless/events/pkg/events"
)

func ParseCommand(txt, prefix string) *events.Command {
	txt = strings.TrimSpace(txt)
	if !strings.HasPrefix(txt, prefix) {
		return nil
	}
	cmd := &events.Command{}

	txt = strings.TrimPrefix(txt, prefix)

	// special case.
	if prefix == "(╯°□°)╯︵ " {
		cmd.Cmd = "flip"
		cmd.Args = txt
		return cmd
	}

	cmdEndIndex := strings.Index(txt, " ")
	if cmdEndIndex == -1 {
		// the command has no arguments.
		cmd.Cmd = txt
	} else {
		cmd.Cmd = txt[:cmdEndIndex]
		cmd.Args = txt[cmdEndIndex+1:]
	}
	return cmd
}
