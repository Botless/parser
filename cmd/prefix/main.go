package main

import (
	"context"
	"github.com/botless/parser/pkg/parser"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/client"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type envConfig struct {
	// Port is server port to be listened.
	Port int `envconfig:"USER_PORT" default:"8080"`

	// Prefix parses commands starting with prefix and creates events.
	Prefix string `envconfig:"PREFIX" default:"/"`
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("[ERROR] Failed to process env var: %s", err)
	}

	c, err := client.NewDefault()
	if err != nil {
		log.Fatalf("Failed to create client: %s", err.Error())
	}

	p := &parser.PrefixParser{
		Ce:     c,
		Prefix: env.Prefix,
	}

	ctx := context.Background()
	log.Fatal(c.StartReceiver(ctx, p.Receive))
}
