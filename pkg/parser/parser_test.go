package parser

import (
	"testing"

	"github.com/botless/events/pkg/events"
	"github.com/google/go-cmp/cmp"
)

func Test_ParseCommand_slash(t *testing.T) {
	testCases := map[string]struct {
		cmd  string
		want *events.Command
	}{
		"not a command": {
			cmd: "Hello, World!",
		},
		"cmd no args": {
			cmd: "/hello",
			want: &events.Command{
				Cmd: "hello",
			},
		},
		"cmd with args": {
			cmd: "/hello world",
			want: &events.Command{
				Cmd:  "hello",
				Args: "world",
			},
		},
		"cmd with more args": {
			cmd: "/hello world, I made it",
			want: &events.Command{
				Cmd:  "hello",
				Args: "world, I made it",
			},
		},
		"cmd with extra whitespace and args": {
			cmd: " \n \t    /hello world, I made it        ",
			want: &events.Command{
				Cmd:  "hello",
				Args: "world, I made it",
			},
		},
		"not a command, has whitespace": {
			cmd: "\t    Hello, World!",
		},
	}
	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {

			got := ParseCommand(tc.cmd, "/")

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("unexpected (-want, +got) = %v", diff)
			}
		})
	}
}

func Test_ParseCommand_hash(t *testing.T) {
	testCases := map[string]struct {
		cmd  string
		want *events.Command
	}{
		"not a command": {
			cmd: "Hello, World!",
		},
		"cmd no args": {
			cmd: "#hello",
			want: &events.Command{
				Cmd: "hello",
			},
		},
		"cmd with args": {
			cmd: "#hello world",
			want: &events.Command{
				Cmd:  "hello",
				Args: "world",
			},
		},
		"cmd with more args": {
			cmd: "#hello world, I made it",
			want: &events.Command{
				Cmd:  "hello",
				Args: "world, I made it",
			},
		},
		"cmd with extra whitespace and args": {
			cmd: " \n \t    #hello world, I made it        ",
			want: &events.Command{
				Cmd:  "hello",
				Args: "world, I made it",
			},
		},
		"not a command, has whitespace": {
			cmd: "\t    Hello, World!",
		},
	}
	for n, tc := range testCases {
		t.Run(n, func(t *testing.T) {

			got := ParseCommand(tc.cmd, "#")

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("unexpected (-want, +got) = %v", diff)
			}
		})
	}

}
