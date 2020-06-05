package main

import (
	"flag"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/ndjordjevic/go-eventhub/internal/protogen/api"
)

// NOTE: Can test with demo servers.
// nats-pub -s demo.nats.io <subject> <msg>
// nats-pub -s demo.nats.io:4443 <subject> <msg> (TLS version)

func usage() {
	log.Printf("Usage: nats-pub [-s server] [-creds file] <subject>\n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var userCreds = flag.String("creds", "", "User Credentials File")
	var showHelp = flag.Bool("h", false, "Show help message")
	var reply = flag.String("reply", "", "Sets a specific reply subject")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("args != 1")
		showUsageAndExit(1)
	}

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Publisher")}

	// Use UserCredentials
	if *userCreds != "" {
		opts = append(opts, nats.UserCredentials(*userCreds))
	}

	// Connect to NATS
	nc, err := nats.Connect(*urls, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	i := &api.Instrument{
		User:        "ndjordjevic",
		MessageType: api.Instrument_UPDATE,
		InstrumentPayload: []*api.Instrument_Payload{
			{
				Id:   1,
				Isin: "BMW1",
			},
		},
	}

	msg, err := proto.Marshal(i)
	if err != nil {
		log.Fatalln("Failed to encode instrument:", err)
	}

	subj := args[0]

	if reply != nil && *reply != "" {
		nc.PublishRequest(subj, *reply, msg)
	} else {
		nc.Publish(subj, msg)
	}

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%v'\n", subj, i)
	}
}
