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

package pulsar_comcast

import (
    "context"
    "github.com/Comcast/pulsar-client-go"
    "github.com/nalej/derrors"
    "github.com/nalej/nalej-bus/pkg/bus"
    "github.com/rs/zerolog/log"
    "time"
)

const (
    // Number of seconds to wait before we consider a timeout for a receive operation
    NalejPulsarSendTimeout = 1
)

type PulsarProducer struct {
    producer *pulsar.ManagedProducer

}


func NewPulsarProducer(client PulsarClient, name string, topic string) bus.NalejProducer {
    // Create a basic config object
    cfg := pulsar.ManagedProducerConfig{
        Name: name,
        Topic: topic,
        NewProducerTimeout: time.Second,
        InitialReconnectDelay: time.Second,
        MaxReconnectDelay: time.Minute,
        ManagedClientConfig: client.config,
    }
    // create producer
    producer := pulsar.NewManagedProducer(client.pool, cfg)

    return PulsarProducer{producer: producer}
}

func(p PulsarProducer) Send(ctx context.Context, msg []byte) derrors.Error {
    _, err := p.producer.Send(ctx, msg)

    if err != nil {
        log.Error().Err(err).Msg("Producer send error")
        return derrors.NewInternalError("impossible to send message", err)
    }
    return nil

}

// Close the producer. This operation must close any connection with brokers with an established connection.
// return:
//  error if any
func (p PulsarProducer) Close(ctx context.Context) derrors.Error {
    err := p.producer.Close(ctx)
    if err != nil {
        return derrors.NewInternalError("impossible to close producer", err)
    }
    return nil
}