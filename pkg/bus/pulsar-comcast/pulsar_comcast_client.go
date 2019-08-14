/*
 *  Copyright (C) 2019 Nalej Group - All Rights Reserved
 */


package pulsar_comcast

import (
    "fmt"
    "github.com/Comcast/pulsar-client-go"
    "github.com/nalej/derrors"
    "github.com/nalej/nalej-bus/pkg/bus"
)

// Basic structure to manage pulsar clients
type PulsarClient struct {
    pool *pulsar.ManagedClientPool
    config pulsar.ManagedClientConfig
}

// NewClient creates a new NalejClient
// asyncErr channel where errors will be sent, may be nil
func NewClient(address string, asyncErr chan error) bus.NalejClient {
    pool := pulsar.NewManagedClientPool()

    config := pulsar.ManagedClientConfig{
        ClientConfig:   pulsar.ClientConfig{
            Addr:       fmt.Sprintf("pulsar://%s", address),
            Errs:       asyncErr,
        },
    }

    return PulsarClient{pool: pool, config: config}
}


func (c PulsarClient) BuildProducer(name string, topic string) (bus.NalejProducer, derrors.Error) {
    if c.pool == nil {
        return nil, derrors.NewInvalidArgumentError("missing configuration and pool to build a pulsar producer")
    }

    prod := NewPulsarProducer(c,name, topic)

    if prod == nil {
        return nil, derrors.NewInternalError("impossible to build producer")
    }

    return prod, nil
}

func (c PulsarClient) BuildConsumer(name string, topic string, exclusive bool) (bus.NalejConsumer, derrors.Error) {
    if c.pool == nil {
        return nil, derrors.NewInvalidArgumentError("missing configuration and pool to build a pulsar consumer")
    }

    cons := NewPulsarConsumer(c, name, topic, exclusive)
    if cons == nil {
        return nil, derrors.NewInternalError("impossible to build a consumer")
    }
    return cons, nil
}


