/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package ops

import (
    "context"
    "github.com/golang/protobuf/proto"
    "github.com/nalej/derrors"
    "github.com/nalej/grpc-bus-go"
    "github.com/nalej/grpc-conductor-go"
    "github.com/nalej/nalej-bus/pkg/bus"
    "github.com/nalej/nalej-bus/pkg/queue"
    "github.com/rs/zerolog/log"
)

const (
    ApplicationOpsTopic="nalej/application/ops"
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
    prod, err := client.BuildProducer(name, ApplicationOpsTopic)
    if err != nil {
        return nil, err
    }
    return &ApplicationOpsProducer{producer: prod}, nil
}

// Generic function that decides what to send depending on the object type.
func (m ApplicationOpsProducer) Send(ctx context.Context, msg proto.Message) derrors.Error {

    var  wrapper grpc_bus_go.ApplicationOps

    switch x := msg.(type) {
    case *grpc_conductor_go.DeploymentRequest:
        wrapper = grpc_bus_go.ApplicationOps{ Operation: &grpc_bus_go.ApplicationOps_DeployRequest{x}}
    case *grpc_conductor_go.UndeployRequest:
        wrapper = grpc_bus_go.ApplicationOps{ Operation: &grpc_bus_go.ApplicationOps_UndeployRequest{x}}
    case *grpc_conductor_go.DrainClusterRequest:
        wrapper = grpc_bus_go.ApplicationOps{ Operation: &grpc_bus_go.ApplicationOps_DrainRequest{x}}
    default:
        return derrors.NewInvalidArgumentError("invalid proto message type")
    }

    payload, err := queue.MarshallPbMsg(&wrapper)
    if err != nil {
        return err
    }

    err = m.producer.Send(ctx, payload)
    if err != nil {
        return err
    }

    return nil
}



// Application consumer
type ApplicationOpsConsumer struct {
    Consumer bus.NalejConsumer
    Config ConfigApplicationOpsConsumer
}

// Struct designed to config a consumer defining what actions to perform depending on the incoming object.
type ConfigApplicationOpsConsumer struct {
    // channel to receive deployment requests
    ChDeploymentRequest chan *grpc_conductor_go.DeploymentRequest
    // channel to receive undeploy requests
    ChUndeployRequest chan *grpc_conductor_go.UndeployRequest
    // channel to receive drain cluster requests
    ChDrainClusterRequest chan *grpc_conductor_go.DrainClusterRequest
    // object types to be considered for consumption
    ToConsume ConsumableStructsApplicationOpsConsumer
}

// Create a new configuration structure for a given channel size
// params:
//  size of channels
// return:
//  instance of a configuration object
func NewConfigApplicationOpsConsumer(size int, toConsume ConsumableStructsApplicationOpsConsumer) ConfigApplicationOpsConsumer {
    chDeploymentRequest := make(chan *grpc_conductor_go.DeploymentRequest, size)
    chUndeployRequest := make(chan *grpc_conductor_go.UndeployRequest, size)
    chDrainClusterRequest := make(chan *grpc_conductor_go.DrainClusterRequest, size)
    return ConfigApplicationOpsConsumer{
        ChDeploymentRequest: chDeploymentRequest,
        ChUndeployRequest: chUndeployRequest,
        ChDrainClusterRequest: chDrainClusterRequest,
        ToConsume: toConsume,
    }
}

// Data struct indicating what data structures available in this topic will be accepted.
type ConsumableStructsApplicationOpsConsumer struct {
    // Consume deploy requests
    DeployRequest bool
    // Consume undeploy requests
    UndeployRequest bool
    // Drain request
    DrainClusterRequest bool
}

func NewApplicationOpsConsumer (client bus.NalejClient, name string, exclusive bool, config ConfigApplicationOpsConsumer) (*ApplicationOpsConsumer, derrors.Error) {
    consumer, err := client.BuildConsumer(name, ApplicationOpsTopic, exclusive)
    if err != nil {
        return nil, err
    }

    return &ApplicationOpsConsumer{Consumer: consumer,Config: config}, nil
}


// Consume any of the possible objects that can be sent to this queue and send it to the corresponding channel.
func (c ApplicationOpsConsumer) Consume(ctx context.Context) derrors.Error{
    msg, err := c.Consumer.Receive(ctx)
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
        if  c.Config.ToConsume.DeployRequest {
            c.Config.ChDeploymentRequest <- x.DeployRequest
        }
    case *grpc_bus_go.ApplicationOps_UndeployRequest:
        if c.Config.ToConsume.UndeployRequest {
            c.Config.ChUndeployRequest <- x.UndeployRequest
        }
    case *grpc_bus_go.ApplicationOps_DrainRequest:
        if c.Config.ToConsume.DrainClusterRequest {
            c.Config.ChDrainClusterRequest <- x.DrainRequest
        }
    case nil:
        errMsg := "received nil entry"
        log.Error().Msg(errMsg)
        return derrors.NewInvalidArgumentError(errMsg)
    default:
        errMsg := "unknown object type in infrastructure ops"
        log.Error().Interface("type",x).Msg(errMsg)
        return derrors.NewInvalidArgumentError(errMsg)
    }
    return nil
}