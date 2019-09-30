/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package ops

import (
    "context"
    "fmt"
    "github.com/golang/protobuf/proto"
    "github.com/nalej/derrors"
    "github.com/nalej/grpc-application-network-go"
    "github.com/nalej/grpc-bus-go"
    "github.com/nalej/grpc-network-go"
    "github.com/nalej/nalej-bus/pkg/bus"
    "github.com/nalej/nalej-bus/pkg/queue"
    "github.com/rs/zerolog/log"
)

const (
    NetworkOpsTopic="nalej/network/ops"
)

type NetworkOpsProducer struct {
    producer bus.NalejProducer
}

// Create a new producer for the application operations topic
// params:
//  client to be used
//  name of the producer
// return:
//  built producer
func NewNetworkOpsProducer (client bus.NalejClient, name string) (*NetworkOpsProducer, derrors.Error) {
    prod, err := client.BuildProducer(name, NetworkOpsTopic)
    if err != nil {
        return nil, err
    }
    return &NetworkOpsProducer{producer: prod}, nil
}

// Generic function that decides what to send depending on the object type.
func (m NetworkOpsProducer) Send(ctx context.Context, msg proto.Message) derrors.Error {

    var  wrapper grpc_bus_go.NetworkOps

    switch x := msg.(type) {
    case *grpc_network_go.AuthorizeMemberRequest:
        wrapper = grpc_bus_go.NetworkOps{Operation: &grpc_bus_go.NetworkOps_AuthorizeMemberRequest{x}}
    case *grpc_network_go.DisauthorizeMemberRequest:
        wrapper = grpc_bus_go.NetworkOps{Operation: &grpc_bus_go.NetworkOps_DisauthorizeMemberRequest{x}}
    case *grpc_network_go.AddDNSEntryRequest:
        wrapper = grpc_bus_go.NetworkOps{Operation: &grpc_bus_go.NetworkOps_AddDnsEntryRequest{x}}
    case *grpc_network_go.DeleteDNSEntryRequest:
        wrapper = grpc_bus_go.NetworkOps{Operation: &grpc_bus_go.NetworkOps_DeleteDnsEntryRequest{x}}
    case *grpc_network_go.InboundServiceProxy:
        wrapper = grpc_bus_go.NetworkOps{Operation: &grpc_bus_go.NetworkOps_InboundAppServiceProxy{x}}
    case *grpc_network_go.OutboundService:
        wrapper = grpc_bus_go.NetworkOps{Operation: &grpc_bus_go.NetworkOps_OutboundAppService{x}}
    case *grpc_application_network_go.AddConnectionRequest:
        wrapper = grpc_bus_go.NetworkOps{Operation: &grpc_bus_go.NetworkOps_AddConnectionRequest{x}}
    case *grpc_application_network_go.RemoveConnectionRequest:
        wrapper = grpc_bus_go.NetworkOps{Operation: &grpc_bus_go.NetworkOps_RemoveConnectionRequest{x}}
    case *grpc_network_go.AuthorizeZTConnectionRequest:
        wrapper = grpc_bus_go.NetworkOps{Operation: &grpc_bus_go.NetworkOps_AuthorizeZtConnection{x}}
    case *grpc_network_go.RegisterZTConnectionRequest:
        wrapper = grpc_bus_go.NetworkOps{Operation: &grpc_bus_go.NetworkOps_RegisterZtConnection{x}}
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

// Network consumer
type NetworkOpsConsumer struct {
    Consumer bus.NalejConsumer
    Config ConfigNetworkOpsConsumer
}

// Struct designed to config a consumer defining what actions to perform depending on the incoming object.
type ConfigNetworkOpsConsumer struct {
    // channel to receive authorize requests
    ChAuthorizeMembersRequest chan *grpc_network_go.AuthorizeMemberRequest
    // channel to receive disauthorize requests
    ChDisauthorizeMembersRequest chan *grpc_network_go.DisauthorizeMemberRequest
    // channel to add DNS entries
    ChAddDNSEntryRequest chan *grpc_network_go.AddDNSEntryRequest
    // channel to delete DNS entries
    ChDeleteDNSEntryRequest chan *grpc_network_go.DeleteDNSEntryRequest
    // channel to receive InboundAppServiceProxy
    ChInboundServiceProxy chan *grpc_network_go.InboundServiceProxy
    // channel to receive OutboundAppService
    ChOutboundService chan *grpc_network_go.OutboundService
    // channel to receive AddConnectionRequest
    ChAddConnectionRequest chan *grpc_application_network_go.AddConnectionRequest
    // channel to receive RemoveConnectionRequest
    ChRemoveConnectionRequest chan *grpc_application_network_go.RemoveConnectionRequest
    // channel to receive AuthorizeZTConnectionRequest
    ChAuthorizeZTConnection chan *grpc_network_go.AuthorizeZTConnectionRequest
    // channel to receive RegisterZTConnectionRequest
    ChRegisterZTConnection chan *grpc_network_go.RegisterZTConnectionRequest
    // object types to be considered for consumption
    ToConsume ConsumableStructsNetworkOpsConsumer
}

func NewConfigNetworksOpsConsumer(size int, toConsume ConsumableStructsNetworkOpsConsumer) ConfigNetworkOpsConsumer {
    chAuthorize := make(chan *grpc_network_go.AuthorizeMemberRequest,size)
    chDisauthorize := make(chan *grpc_network_go.DisauthorizeMemberRequest,size)
    chAddDNSEntry := make(chan *grpc_network_go.AddDNSEntryRequest, size)
    chDeleteDNSEntry := make(chan *grpc_network_go.DeleteDNSEntryRequest, size)
    chInboundServiceProxy := make(chan *grpc_network_go.InboundServiceProxy, size)
    chOutboundServiceProxy := make(chan *grpc_network_go.OutboundService, size)
    chAddConnection := make (chan * grpc_application_network_go.AddConnectionRequest, size)
    chRemoveConnection := make (chan * grpc_application_network_go.RemoveConnectionRequest, size)
    chAuthorizeZTConnection := make(chan * grpc_network_go.AuthorizeZTConnectionRequest, size)
    chRegisterZTConnection := make (chan * grpc_network_go.RegisterZTConnectionRequest, size)



    return ConfigNetworkOpsConsumer{ChAuthorizeMembersRequest: chAuthorize,
        ChDisauthorizeMembersRequest: chDisauthorize,
        ChAddDNSEntryRequest:       chAddDNSEntry,
        ChDeleteDNSEntryRequest:    chDeleteDNSEntry,
        ChInboundServiceProxy:      chInboundServiceProxy,
        ChOutboundService:          chOutboundServiceProxy,
        ChAddConnectionRequest:     chAddConnection,
        ChRemoveConnectionRequest:  chRemoveConnection,
        ChAuthorizeZTConnection:    chAuthorizeZTConnection,
        ChRegisterZTConnection:     chRegisterZTConnection,
        ToConsume: toConsume}
}

// Data struct indicating what data structures available in this topic will be accepted.
type ConsumableStructsNetworkOpsConsumer struct {
    // Consume authorize members
    AuthorizeMember bool
    // Consume disauthorize members
    DisauthorizeMember bool
    // Consume add dns entries
    AddDNSEntry bool
    // Consume delete dns entries
    DeleteDNSEntry bool
    // Consume inbound service proxy
    InboundServiceProxy bool
    // Consume outbound service
    OutboundService bool
    // Consume add connection
    AddConnection bool
    // Consume remove connection
    RemoveConnection bool
    // Consume authorize zt connection
    AuthorizeZTConnection bool
    // Consume register zt connection
    RegisterZTConnecion bool
}


func NewNetworkOpsConsumer (client bus.NalejClient, name string, exclusive bool, config ConfigNetworkOpsConsumer) (*NetworkOpsConsumer, derrors.Error) {
    consumer, err := client.BuildConsumer(name, NetworkOpsTopic, exclusive)
    if err != nil {
        return nil, err
    }

    return &NetworkOpsConsumer{Consumer: consumer,Config: config}, nil
}


// Consume any of the possible objects that can be sent to this queue and send it to the corresponding channel.
func (c NetworkOpsConsumer) Consume(ctx context.Context) derrors.Error{
    msg, err := c.Consumer.Receive(ctx)
    if err != nil {
        return err
    }

    target := &grpc_bus_go.NetworkOps{}

    derr := queue.UnmarshallPbMsg(msg, target)
    if derr != nil {
        return derr
    }

    switch x := target.Operation.(type) {
    case *grpc_bus_go.NetworkOps_AuthorizeMemberRequest:
        if  c.Config.ToConsume.AuthorizeMember {
            c.Config.ChAuthorizeMembersRequest <- x.AuthorizeMemberRequest
        }
    case *grpc_bus_go.NetworkOps_DisauthorizeMemberRequest:
        if  c.Config.ToConsume.DisauthorizeMember {
            c.Config.ChDisauthorizeMembersRequest <- x.DisauthorizeMemberRequest
        }
    case *grpc_bus_go.NetworkOps_AddDnsEntryRequest:
        if c.Config.ToConsume.AddDNSEntry {
            c.Config.ChAddDNSEntryRequest <- x.AddDnsEntryRequest
        }
    case *grpc_bus_go.NetworkOps_DeleteDnsEntryRequest:
        if c.Config.ToConsume.DeleteDNSEntry {
            c.Config.ChDeleteDNSEntryRequest <- x.DeleteDnsEntryRequest
        }
    case *grpc_bus_go.NetworkOps_InboundAppServiceProxy:
        if c.Config.ToConsume.InboundServiceProxy {
            c.Config.ChInboundServiceProxy <- x.InboundAppServiceProxy
        }
    case *grpc_bus_go.NetworkOps_OutboundAppService:
        if c.Config.ToConsume.OutboundService {
            c.Config.ChOutboundService <- x.OutboundAppService
        }
    case *grpc_bus_go.NetworkOps_AddConnectionRequest:
        if c.Config.ToConsume.AddConnection {
            c.Config.ChAddConnectionRequest <- x.AddConnectionRequest
        }
    case *grpc_bus_go.NetworkOps_RemoveConnectionRequest:
        if c.Config.ToConsume.RemoveConnection {
            c.Config.ChRemoveConnectionRequest <- x.RemoveConnectionRequest
        }
    case *grpc_bus_go.NetworkOps_AuthorizeZtConnection:
        if c.Config.ToConsume.AuthorizeZTConnection {
            c.Config.ChAuthorizeZTConnection <- x.AuthorizeZtConnection
        }
    case *grpc_bus_go.NetworkOps_RegisterZtConnection:
        if c.Config.ToConsume.RegisterZTConnecion {
            c.Config.ChRegisterZTConnection <- x.RegisterZtConnection
        }
    case nil:
        errMsg := "received nil entry"
        log.Error().Msg(errMsg)
        return derrors.NewInvalidArgumentError(errMsg)
    default:
        errMsg := fmt.Sprintf("unknown object type in %s",NetworkOpsTopic)
        log.Error().Interface("type",x).Msg(errMsg)
        return derrors.NewInvalidArgumentError(errMsg)
    }
    return nil
}