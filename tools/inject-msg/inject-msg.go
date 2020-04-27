package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/openware/rango/pkg/upstream"
	"github.com/rs/zerolog/log"
)

func main() {
	addr := "amqp://guest:guest@localhost:5672/"
	mq := upstream.NewAMQPSession("", addr)

	for !mq.IsReady {
		log.Info().Msg("Not connected, waiting")
		time.Sleep(time.Second * 1)
	}

	for {
		file, err := os.Open("msg.txt")
		if err != nil {
			panic(err.Error())
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				panic(err.Error())
			}

			msg := strings.Split(scanner.Text(), " ")

			if err := mq.Push(upstream.PeatioRangerEventsEx, msg[0], []byte(msg[1])); err != nil {
				fmt.Printf("Push failed: %s\n", err)
			} else {
				log.Info().Msgf("Pushed on %s msg: %s", msg[0], msg[1])
			}

		}
		file.Close()
		log.Info().Msg("Waiting 5 seconds")
		time.Sleep(time.Second * 5)
	}
}