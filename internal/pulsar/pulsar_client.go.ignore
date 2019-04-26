/*
 *  Copyright (C) 2019 Nalej Group - All Rights Reserved
 */


package pulsar

import (
    "github.com/apache/pulsar/pulsar-client-go/pulsar"
    "github.com/nalej/derrors"
    "time"
)


// Utils to deal with Pulsar clients and configuration

// Create a basic pulsar client.
// params:
//  host where the client must connect with. Example: myhost:6650
//  timeoutSeconds number of seconds before considering a timeout
//  listenerThreads number of threads used in the system
// return:
//  client if everything was correct
//  error if any
func NewClient(host string, timeoutSeconds int, listenerThreads int) (pulsar.Client, derrors.Error) {
    client, err := pulsar.NewClient(pulsar.ClientOptions{
        URL: "pulsar://"+host,
        OperationTimeoutSeconds: time.Duration(timeoutSeconds)*time.Second,
        MessageListenerThreads: listenerThreads,
    })

    if err != nil {
        return nil, derrors.NewInternalError("impossible to create pulsar client", err)
    }

    return client, nil

}