package parser

import (
	"context"
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

func (p *PrefixParser) Receive(event cloudevents.Event) {
	// don't block the caller.
	go p.parse(event)
}

func (p *PrefixParser) parse(event cloudevents.Event) {
	switch event.Type() {
	case "botless.slack.message":
		p.parseSlackMessage(event)
	default:
		// ignore
		log.Printf("botless parser ignored event type %q", event.Type())
	}
}

func (p *PrefixParser) parseSlackMessage(event cloudevents.Event) {

	msg := &slack.MessageEvent{}
	if err := event.DataAs(msg); err != nil {
		log.Printf("failed to get slack message event from cloudevent of type %s", event.Type())
		return
	}

	txt := msg.Msg.Text

	cmd := ParseCommand(txt, p.Prefix)
	if cmd != nil {

		ec := event.Context.AsV02()

		log.Printf("got: %+v", msg)

		cmd.Author = msg.Msg.Name
		cmd.Channel = msg.Channel

		cmdEvent := cloudevents.Event{
			Context: cloudevents.EventContextV02{
				Type:       events.Bot.Type("command", cmd.Cmd),
				Source:     ec.Source,
				Time:       ec.Time,
				Extensions: ec.Extensions, // pass all extensions along.
			}.AsV02(),
			Data: cmd,
		}
		if _, err := p.Ce.Send(context.TODO(), cmdEvent); err != nil {
			log.Printf("failed to send botless command %s", err.Error())
		} else {
			log.Printf("sent: %s command: %+v", cmdEvent.Type(), cmd)
		}
	}
}
