package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/containerssh/auditlog/codec/binary"
	"github.com/containerssh/auditlog/message"
)

func main() {
	file := ""
	flag.StringVar(&file, "file", "", "File to process")
	flag.Parse()
	if flag.NArg() != 0 {
		flag.Usage()
		os.Exit(1)
	}

	if file == "" {
		flag.Usage()
		os.Exit(1)
	}

	fh, err := os.Open(file)
	if err != nil {
		log.Fatalf("failed to open audit log file %s (%v)", file, err)
	}

	decoder := binary.NewDecoder()
	messages, errors := decoder.Decode(fh)
	for {
		var msg *message.Message
		select {
		case msg = <-messages:
			if msg == nil {
				break
			}

			var data []byte
			data, err = json.Marshal(msg)
			if err != nil {
				structuredError := map[string]string{
					"error": fmt.Sprintf("JSON encoding error: (%v)", err),
				}
				data, _ = json.Marshal(structuredError)
				_, _ = os.Stdout.Write(data)
				_, _ = os.Stdout.Write([]byte("\n"))
			} else if data != nil {
				_, _ = os.Stdout.Write(data)
				_, _ = os.Stdout.Write([]byte("\n"))
			} else {
				break
			}
		case channelError := <-errors:
			if channelError != nil {
				structuredError := map[string]string{
					"error": channelError.Error(),
				}
				data, _ := json.Marshal(structuredError)
				_, _ = os.Stdout.Write(data)
				_, _ = os.Stdout.Write([]byte("\n"))
			} else {
				break
			}
		}
	}
}
