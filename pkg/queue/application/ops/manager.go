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
}

// Struct designed to config a consumer defining what actions to perform depending on the incoming object.
type ConfigApplicationOpsConsumer struct {
    // function to process deployment requests
    fDeploymentReq func(in grpc_conductor_go.DeploymentRequest)()
    // function to process undeploy requests
    fUndeployReq func(in grpc_conductor_go.UndeployRequest)()
}

func NewApplicationOpsConsumer (client bus.NalejClient, name string, exclusive bool) (*ApplicationOpsConsumer, derrors.Error) {
    consumer, err := client.BuildConsumer(name, InfrastructureOpsTopic, exclusive)
    if err != nil {
        return nil, err
    }

    return &ApplicationOpsConsumer{consumer: consumer}, nil

}

/*
// Consume any of the potential objects
func (c ApplicationOpsConsumer) Consume() (proto.Message, interface{}, derrors.Error) {
    msg, err := c.consumer.Receive()
    if err != nil {
        return nil, nil, err
    }

    target := &grpc_bus_go.ApplicationOps{}

    derr := queue.UnmarshallPbMsg(msg, target)
    if derr != nil {
        return nil, nil, derr
    }

    switch x := target.Operation.(type) {
    case *grpc_bus_go.ApplicationOps_DeployRequest:

        //return x.DeployRequest, grpc_conductor_go.DeploymentRequest.ProtoMessage, nil
    case *grpc_bus_go.ApplicationOps_UndeployRequest:
        return x.UndeployRequest, grpc_conductor_go.DeploymentRequest.ProtoMessage, nil
    case nil:
        return nil, nil, derrors.NewInternalError("applicationOpsConsumer was not set in consume method")
    default:
        return nil, nil, derrors.NewInternalError(fmt.Sprintf("Profile.Avatar has unexpected type %T", x))
    }
}
*/