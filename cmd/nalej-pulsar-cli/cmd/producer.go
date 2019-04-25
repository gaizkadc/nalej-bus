/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package cmd

import (
    "fmt"
    "github.com/nalej/nalej-bus/internal/pulsar"
    "github.com/rs/zerolog/log"
    "github.com/spf13/cobra"
    "time"
)

// producerTopic name
var producerTopic string
// producer name
var producerName string

var producerCmd = &cobra.Command{
    Use: "producer",
    Short: "run a producer example",
    Long: "run a producer example",
    Run: func(cmd *cobra.Command, args []string){
      SetupLogging()
      runProducer()
    },
}

func init() {
    producerCmd.Flags().StringVar(&producerTopic, "producerTopic", "public/default/topic", "Topic this producer will publish into")
    producerCmd.Flags().StringVar(&producerName, "producerName", "producer", "Name for this producer")
    RootCmd.AddCommand(producerCmd)
}

func runProducer() {
    client, err := pulsar.NewClient(pulsarAddress,5, 1)
    if err != nil {
        log.Panic().Msg(err.Error())
    }

    producer, err := pulsar.NewPulsarProducer(client, producerName, producerTopic)
    if err != nil {
        log.Panic().Msg(err.Error())
    }
    defer producer.Close()
    defer client.Close()

    counter := 0
    tick := time.Tick(time.Second)
    for {
        select {
            case <- tick:
                msg := fmt.Sprintf("Message number %d", counter)
                err = producer.Send([]byte(msg))
                if err != nil {
                    log.Error().Err(err)
                } else {
                    log.Info().Str("msg", msg).Msg("->")
                }
                counter = counter + 1
        }
    }

}