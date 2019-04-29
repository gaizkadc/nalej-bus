/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package pulsar_comcast

import (
    "context"
    "github.com/Comcast/pulsar-client-go"
    "github.com/nalej/derrors"
    "github.com/nalej/nalej-bus/pkg/bus"
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

// Send a new message to the topic of this producer.
// params:
//  msg message to be sent
// return:
//  error if any
func(p PulsarProducer) Send(msg []byte) derrors.Error {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second * NalejPulsarSendTimeout)
    defer cancel()

    _, err := p.producer.Send(ctx, msg)

    if err != nil {
        return derrors.NewInternalError("impossible to send message", err)
    }
    return nil

}

// Close the producer. This operation must close any connection with brokers with an established connection.
// return:
//  error if any
func (p PulsarProducer) Close() derrors.Error {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()
    err := p.producer.Close(ctx)
    if err != nil {
        return derrors.NewInternalError("impossible to close producer", err)
    }
    return nil
}