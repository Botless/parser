package parser

import (
	"fmt"
	"github.com/botless/events/pkg/events"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/nlopes/slack"
	"log"
)

type PrefixParser struct {
	Ce     client.Client
	Prefix string
}

func (p *PrefixParser) Receive(event cloudevents.Event, resp *cloudevents.EventResponse) {
	switch event.Type() {
	case "botless.slack.message":
		if cmd, err := p.parseSlackMessage(event); err != nil {
			log.Printf("failed to parse message: %s", err)
		} else if cmd != nil {
			resp.RespondWith(200, cmd)
		}

	default:
		// ignore
		log.Printf("botless parser ignored event type %q", event.Type())
	}
}

func (p *PrefixParser) parseSlackMessage(event cloudevents.Event) (*cloudevents.Event, error) {
	msg := &slack.MessageEvent{}
	if err := event.DataAs(msg); err != nil {
		return nil, fmt.Errorf("failed to get slack message event from cloudevent of type %s", event.Type())
	}

	txt := msg.Msg.Text

	cmd := ParseCommand(txt, p.Prefix)
	if cmd != nil {
		ec := event.Context.AsV02()

		cmd.Author = msg.Msg.Name
		cmd.Channel = msg.Channel

		return &cloudevents.Event{
			Context: cloudevents.EventContextV02{
				Type:       events.Bot.Type("command", cmd.Cmd),
				Source:     ec.Source,
				Time:       ec.Time,
				Extensions: ec.Extensions, // pass all extensions along.
			}.AsV02(),
			Data: cmd,
		}, nil
	}
	return nil, nil
}
