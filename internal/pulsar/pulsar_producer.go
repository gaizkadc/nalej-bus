/*
 *  Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package pulsar

import (
    "context"
    "github.com/apache/pulsar/pulsar-client-go/pulsar"
    "github.com/nalej/derrors"
    "github.com/nalej/nalej-bus/internal/bus"
)

// Nalej producer object is a wrapper of an Apache Pulsar client.
type PulsarProducer struct {
    producer pulsar.Producer
}

// It returns a simple pulsar client for the given topic
// params:
//  client
//  name of this producer
//  topic
func NewPulsarProducer(client pulsar.Client, name string, topic string) (bus.NalejProducer, derrors.Error) {
    if client == nil {
        return nil, derrors.NewFailedPreconditionError("received nil pulsar client")
    }
    internalProducer,err := client.CreateProducer(
        pulsar.ProducerOptions{
        Topic: topic,
        Name: name,
    })

    if err != nil {
        return nil, derrors.NewInternalError("impossible to create producer", err)
    }

    return &PulsarProducer{producer:internalProducer}, nil
}

func(prod PulsarProducer) Send(msg []byte) derrors.Error {
    err := prod.producer.Send(context.Background(), pulsar.ProducerMessage{Payload: msg})
    if err != nil {
        return derrors.NewInternalError("impossible to send message", err)
    }
    return nil
}

func(prod PulsarProducer) Close() derrors.Error {
    err := prod.producer.Close()
    if err != nil {
        return derrors.NewInternalError("impossible to close producer", err)
    }
    return nil
}