/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */


package cmd

import (
    "context"
    "github.com/nalej/nalej-bus/pkg/bus/pulsar-comcast"
    "github.com/rs/zerolog/log"
    "github.com/spf13/cobra"
    "time"
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
    client := pulsar_comcast.NewClient(pulsarAddress, nil)

    consumer,error := client.BuildConsumer(consumerName, topicConsumer, true)
    if error!=nil{
        log.Panic().Err(error).Msg("impossible to build consumer")
    }
    ctx,cancel := context.WithTimeout(context.Background(), time.Second * 5)
    defer consumer.Close(ctx)
    cancel()

    for {
        ctx,cancel := context.WithTimeout(context.Background(), time.Second * 10)
        msg, err := consumer.Receive(ctx)
        cancel()
        if err != nil {
            log.Error().Err(err)
        } else {
            log.Info().Str("received",string(msg)).Msg("<-")
        }

    }

}

