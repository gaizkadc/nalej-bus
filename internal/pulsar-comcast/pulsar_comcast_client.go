/*
 *  Copyright (C) 2019 Nalej Group - All Rights Reserved
 */


package pulsar_comcast

import (
    "fmt"
    "github.com/Comcast/pulsar-client-go"
)

// Basic structure to manage pulsar clients
type PulsarClient struct {
    pool *pulsar.ManagedClientPool
    config pulsar.ManagedClientConfig
}

func NewClient(address string) PulsarClient {
    pool := pulsar.NewManagedClientPool()

    config := pulsar.ManagedClientConfig{
        ClientConfig: pulsar.ClientConfig{
            Addr:        fmt.Sprintf("pulsar://%s", address),
        },
    }

    return PulsarClient{pool: pool, config: config}
}

