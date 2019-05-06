/*
 *  Copyright (C) 2019 Nalej Group - All Rights Reserved
 */

package bus

import (
    "context"
    "github.com/nalej/derrors"
)

// Interface for nalej clients.
type NalejClient interface {

    // Build a producer using this client
    // params:
    //  name of the producer
    //  topic this producer will use
    // return:
    //  resulting producer
    //  error if any
    BuildProducer(name string, topic string) (NalejProducer, derrors.Error)

    // Build a consumer using this client
    // params:
    //  name of the producer
    //  topic this producer will use
    //  exclusive sets this consumer as the only one consuming data from the topic
    // return:
    //  resulting consumer
    //  error if any
    BuildConsumer(name string, topic string, exclusive bool) (NalejConsumer, derrors.Error)
}

// Main interfaces to be implemented by any producer or consumer in the Nalej platform
type NalejConsumer interface {
    // Receive a message from a subscribed entry
    // params:
    //  ctx context
    // return:
    //  message payload
    //  error if any
    Receive(ctx context.Context) ([]byte, derrors.Error)

    // Close the consumer
    // params:
    //  ctx context
    // return:
    //  error if any
    Close(ctx context.Context) derrors.Error
}

type NalejProducer interface {

    // Send a new message to the topic of this producer.
    // params:
    //  msg message to be sent
    //  ctx context
    // return:
    //  error if any
    Send(msg []byte, ctx context.Context) derrors.Error

    // Close the producer. This operation must close any connection with brokers with an established connection.
    // params:
    //  ctx context
    // return:
    //  error if any
    Close(ctx context.Context) derrors.Error
}