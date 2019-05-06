/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package ops

import (
    "github.com/golang/protobuf/proto"
    "github.com/nalej/derrors"
    "github.com/nalej/grpc-bus-go"
    "github.com/nalej/grpc-conductor-go"
    "github.com/nalej/nalej-bus/pkg/bus"
    "github.com/nalej/nalej-bus/pkg/queue"
    "github.com/rs/zerolog/log"
)

const (
    InfrastructureOpsTopic="nalej/infrastructure/ops"
)

type ApplicationOpsProducer struct {
    producer bus.NalejProducer
}

// Create a new producer for the application operations topic
// params:
//  client to be used
//  name of the producer
// return:
//  built producer
func NewApplicationOpsProducer (client bus.NalejClient, name string) (*ApplicationOpsProducer, derrors.Error) {
    prod, err := client.BuildProducer(name, InfrastructureOpsTopic)
    if err != nil {
        return nil, err
    }
    return &ApplicationOpsProducer{producer: prod}, nil
}

// Generic function that decides what to send depending on the object type.
func (m ApplicationOpsProducer) Send(msg proto.Message) derrors.Error {

    var  wrapper grpc_bus_go.ApplicationOps

    switch x := msg.(type) {
    case *grpc_conductor_go.DeploymentRequest:
        wrapper = grpc_bus_go.ApplicationOps{ Operation: &grpc_bus_go.ApplicationOps_DeployRequest{x}}
    case *grpc_conductor_go.UndeployRequest:
        wrapper = grpc_bus_go.ApplicationOps{ Operation: &grpc_bus_go.ApplicationOps_UndeployRequest{x}}
    default:
        return derrors.NewInvalidArgumentError("invalid proto message type")
    }

    payload, err := queue.MarshallPbMsg(&wrapper)
    if err != nil {
        return err
    }

    err = m.producer.Send(payload)
    if err != nil {
        return err
    }

    return nil
}



// Application consumer
type ApplicationOpsConsumer struct {
    consumer bus.NalejConsumer
    config ConfigApplicationOpsConsumer
}

// Struct designed to config a consumer defining what actions to perform depending on the incoming object.
type ConfigApplicationOpsConsumer struct {
    // channel to receive deployment requests
    chDeploymentRequest chan *grpc_conductor_go.DeploymentRequest
    // channel to receive undeploy requests
    chUndeployRequest chan *grpc_conductor_go.UndeployRequest
}

func NewApplicationOpsConsumer (client bus.NalejClient, name string, exclusive bool, config ConfigApplicationOpsConsumer) (*ApplicationOpsConsumer, derrors.Error) {
    consumer, err := client.BuildConsumer(name, InfrastructureOpsTopic, exclusive)
    if err != nil {
        return nil, err
    }

    return &ApplicationOpsConsumer{consumer: consumer}, nil
}


// Consume any of the possible objects that can be sent to this queue and send it to the corresponding channel.
func (c ApplicationOpsConsumer) Consume() derrors.Error{
    msg, err := c.consumer.Receive()
    if err != nil {
        return err
    }

    target := &grpc_bus_go.ApplicationOps{}

    derr := queue.UnmarshallPbMsg(msg, target)
    if derr != nil {
        return derr
    }

    switch x := target.Operation.(type) {
    case *grpc_bus_go.ApplicationOps_DeployRequest:
        c.config.chDeploymentRequest <- x.DeployRequest
    case *grpc_bus_go.ApplicationOps_UndeployRequest:
        c.config.chUndeployRequest <- x.UndeployRequest
    case nil:
        errMsg := "received nil entry"
        log.Error().Msg(errMsg)
        return derrors.NewInvalidArgumentError(errMsg)
    default:
        errMsg := "unknown object type in infrastructure ops"
        log.Error().Interface("type",x).Msg(erMsg)
        return derrors.NewInvalidArgumentError(errMsg)
    }
    return nil
}