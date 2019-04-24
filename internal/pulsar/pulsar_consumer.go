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


// Encapsulate the pulsar consumer types
type PulsarConsumerType int

const (
    PulsarExclusiveConsumer PulsarConsumerType = iota
    PulsarSharedConsumer
    PulsarFailoverConsumer
)

func(p PulsarConsumerType) translate() pulsar.SubscriptionType {
    switch p {
    case PulsarExclusiveConsumer:
        return pulsar.Exclusive
    case PulsarSharedConsumer:
        return pulsar.Shared
    case PulsarFailoverConsumer:
        return pulsar.Failover
    }
    return -1
}



// Wrapper of a pulsar consumer using Nalej specification
type PulsarConsumer struct {
    consumer pulsar.Consumer
}

// It returns a simple pulsar client for the given topic
// params:
//  client
//  name of this subscription
//  topic to subscribe to
//  type of consumer: exclusive, shared or failover. See pulsar documentation for more details.
//
func NewPulsarConsumer(client pulsar.Client, name string, topic string, consumerType PulsarConsumerType) (bus.NalejConsumer, derrors.Error) {
    if client == nil {
        return nil, derrors.NewFailedPreconditionError("received nil pulsar client")
    }
    internalConsumer, err := client.Subscribe(pulsar.ConsumerOptions{
        SubscriptionName: name,
        Topic: topic,
        Type: consumerType.translate(),
    })
    if err != nil {
        return nil, derrors.NewInternalError("impossible to create consumer", err)
    }
    return PulsarConsumer{consumer: internalConsumer}, nil
}


func (c PulsarConsumer) Receive() ([]byte, derrors.Error) {
    msg, err := c.consumer.Receive(context.Background())
    if err != nil {
        return nil, derrors.NewInternalError("failed receiving message", err)
    }
    // By default we Acknowledge this message
    c.consumer.Ack(msg)
    return msg.Payload(), nil
}

func (c PulsarConsumer) Close() derrors.Error {
    err := c.consumer.Close()
    if err != nil {
        return derrors.NewInternalError("impossible to close consumer", err)
    }
    return nil
}
