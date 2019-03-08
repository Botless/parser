package slash

import (
	"context"
	"log"
	"strings"

	"github.com/botless/events/pkg/events"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/nlopes/slack"
)

type Slash struct {
	c client.Client
}

func New(c client.Client) *Slash {
	return new(Slash)
}

func (p *Slash) Receive(event cloudevents.Event) {
	// don't block the caller.
	go p.parse(event)
}

func (p *Slash) parse(event cloudevents.Event) {

	switch event.Type() {
	case "botless.slack.message":
		p.parseSlackMessage(event)
	default:
		// ignore
		log.Printf("botless parser ignored event type %q", event.Type())
	}
}

func (p *Slash) parseSlackMessage(event cloudevents.Event) {

	msg := &slack.MessageEvent{}
	if err := event.DataAs(msg); err != nil {
		log.Printf("failed to get slack message event from cloudevent of type %s", event.Type())
		return
	}

	txt := msg.Msg.Text

	cmd := parseCommand(txt)
	if cmd != nil {

		ec := event.Context.AsV02()

		cmd.Author = msg.Msg.Username
		cmd.Channel = msg.Channel

		cmdEvent := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				Type:       events.Bot.Type("command"),
				Source:     ec.Source,
				Time:       ec.Time,
				Extensions: ec.Extensions, // pass all extensions along.
			}.AsV02(),
			Data: cmd,
		}
		if _, err := p.c.Send(context.TODO(), cmdEvent); err != nil {
			log.Printf("failed to send botless command %s", err.Error())
		} else {
			log.Printf("sent: %s command: %+v", cmdEvent.Type(), cmd)
		}
	}
}

func parseCommand(txt string) *events.Command {
	txt = strings.TrimSpace(txt)
	if !strings.HasPrefix(txt, "/") {
		return nil
	}
	cmd := &events.Command{}

	txt = strings.TrimPrefix(txt, "/")

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
