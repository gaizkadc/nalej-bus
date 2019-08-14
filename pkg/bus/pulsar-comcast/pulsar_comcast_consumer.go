/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
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
    NalejPulsarReceiveACKTimeout = 2
)

type PulsarConsumer struct {
    consumer *pulsar.ManagedConsumer
}

func NewPulsarConsumer(client PulsarClient, name string, topic string, exclusive bool) bus.NalejConsumer {
    config := pulsar.ManagedConsumerConfig{
        Name: name,
        Topic: topic,
        Exclusive: exclusive,
        ManagedClientConfig: client.config,
    }
    consumer := pulsar.NewManagedConsumer(client.pool, config)
    return PulsarConsumer{consumer: consumer}
}

// Receive a message from a subscribed entry
// return:
//  message payload
//  error if any
func (c PulsarConsumer) Receive(ctx context.Context) ([]byte, derrors.Error) {
    msg, err := c.consumer.Receive(ctx)
    if err != nil {
        log.Error().Err(err).Msg("Consumer receive error")
        return nil, derrors.NewInternalError("failed receiving message", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), NalejPulsarReceiveACKTimeout * time.Second)
    defer cancel()
    err = c.consumer.Ack(ctx, msg)
    if err != nil {
        return nil, derrors.NewInternalError("impossible to acknowledge message", err)
    }
    return msg.Payload, nil
}

// Close the consumer
// return:
//  error if any
func (c PulsarConsumer) Close(ctx context.Context) derrors.Error {
    err := c.consumer.Close(ctx)
    if err != nil {
        return derrors.NewInternalError("impossible to close consumer", err)
    }
    return nil

}
