/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */


package cmd

import (
    "github.com/nalej/nalej-bus/internal/pulsar"
    "github.com/rs/zerolog/log"
    "github.com/spf13/cobra"
)

// producerTopic name
var topicConsumer string
// consumer name
var consumerName string
// consumer subscription
var subscription string

var consumerCmd = &cobra.Command{
    Use: "consumer",
    Short: "run a consumer example",
    Long: "run a consumer example",
    Run: func(cmd *cobra.Command, args []string){
        SetupLogging()
        runConsumer()
    },
}

func init() {
    producerCmd.Flags().StringVar(&topicConsumer, "topicConsumer", "public/default/topic", "Topic this consumer will publish into")
    producerCmd.Flags().StringVar(&consumerName, "consumerName", "consumer", "Name for this consumer")
    RootCmd.AddCommand(consumerCmd)
}

func runConsumer () {
    client, err := pulsar.NewClient(pulsarAddress,5, 1)
    if err != nil {
        log.Panic().Msg(err.Error())
    }


    consumer,err := pulsar.NewPulsarConsumer(client, consumerName, topicConsumer, pulsar.PulsarExclusiveConsumer)
    if err != nil {
        log.Panic().Msg(err.Error())
    }

    for {
        msg, err := consumer.Receive()
        if err != nil {
            log.Error().Err(err)
        } else {
            log.Info().Bytes("received",msg).Msg("<-")
        }

    }

    defer consumer.Close()
    defer client.Close()


}

