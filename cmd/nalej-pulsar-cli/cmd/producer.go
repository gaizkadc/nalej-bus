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
    "fmt"
    "github.com/nalej/nalej-bus/pkg/bus/pulsar-comcast"
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
    producerCmd.Flags().StringVar(&producerTopic, "topic", "public/default/topic", "Topic this producer will publish into")
    producerCmd.Flags().StringVar(&producerName, "name", "producer", "Name for this producer")
    RootCmd.AddCommand(producerCmd)
}

func runProducer() {
    client := pulsar_comcast.NewClient(pulsarAddress, nil)

    producer, error := client.BuildProducer(producerName, producerTopic)
    if error != nil {
        log.Panic().Err(error).Msg("Impossible to build producer")
    }

    ctx,cancel := context.WithTimeout(context.Background(), time.Second * 5)
    defer producer.Close(ctx)
    cancel()

    counter := 0
    tick := time.Tick(time.Second)
    for {
        select {
            case <- tick:
                msg := fmt.Sprintf("Message number %d", counter)
                ctx, cancel = context.WithTimeout(context.Background(), time.Second * 5)
                err := producer.Send(ctx,[]byte(msg))
                cancel()
                if err != nil {
                    log.Error().Err(err).Msg("")
                } else {
                    log.Info().Str("msg", msg).Msg("->")
                }
                counter = counter + 1
        }
    }

}