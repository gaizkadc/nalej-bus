/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */


package cmd

import (
    "github.com/nalej/nalej-bus/internal/pulsar-comcast"
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
    consumerCmd.Flags().StringVar(&topicConsumer, "topic", "public/default/topic", "Topic this consumer will publish into")
    consumerCmd.Flags().StringVar(&consumerName, "name", "consumer", "Name for this consumer")
    RootCmd.AddCommand(consumerCmd)
}

func runConsumer () {
    client := pulsar_comcast.NewClient(pulsarAddress)

    consumer,error := client.BuildConsumer(consumerName, topicConsumer, true)
    if error!=nil{
        log.Panic().Err(error).Msg("impossible to build consumer")
    }

    defer consumer.Close()

    for {
        msg, err := consumer.Receive()
        if err != nil {
            log.Error().Err(err)
        } else {
            log.Info().Str("received",string(msg)).Msg("<-")
        }

    }



}

