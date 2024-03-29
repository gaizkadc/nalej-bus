/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package events

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-bus-go"
	"github.com/nalej/grpc-connectivity-manager-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/nalej-bus/pkg/bus"
	"github.com/nalej/nalej-bus/pkg/queue"
	"github.com/rs/zerolog/log"
)

const (
	InfrastructureEventsTopic = "nalej/infrastructure/events"
)

type InfrastructureEventsProducer struct {
	producer bus.NalejProducer
}

// Create a new producer for the infrastructure events topic
// params:
//  client to be used
//  name of the producer
// return:
//  built producer
func NewInfrastructureEventsProducer(client bus.NalejClient, name string) (*InfrastructureEventsProducer, derrors.Error) {
	prod, err := client.BuildProducer(name, InfrastructureEventsTopic)
	if err != nil {
		return nil, err
	}
	return &InfrastructureEventsProducer{producer: prod}, nil
}

// Generic function that decides what to send depending on the object type.
func (m InfrastructureEventsProducer) Send(ctx context.Context, msg proto.Message) derrors.Error {

	var wrapper grpc_bus_go.InfrastructureEvents

	switch x := msg.(type) {
	case *grpc_infrastructure_go.UpdateClusterRequest:
		wrapper = grpc_bus_go.InfrastructureEvents{Event: &grpc_bus_go.InfrastructureEvents_UpdateClusterRequest{x}}
	case *grpc_infrastructure_go.SetClusterStatusRequest:
		wrapper = grpc_bus_go.InfrastructureEvents{Event: &grpc_bus_go.InfrastructureEvents_SetClusterStatusRequest{x}}
	case *grpc_connectivity_manager_go.ClusterAlive:
		wrapper = grpc_bus_go.InfrastructureEvents{Event: &grpc_bus_go.InfrastructureEvents_ClusterAlive{x}}
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
type InfrastructureEventsConsumer struct {
	Consumer bus.NalejConsumer
	Config   ConfigInfrastructureEventsConsumer
}

// Struct designed to config a consumer defining what actions to perform depending on the incoming object.
type ConfigInfrastructureEventsConsumer struct {
	// channel to receive update cluster requests
	ChUpdateClusterRequest chan *grpc_infrastructure_go.UpdateClusterRequest
	// channel to receive set cluster status requests
	ChSetClusterStatusRequest chan *grpc_infrastructure_go.SetClusterStatusRequest
	// channel to receive cluster alive requests
	ChClusterAlive chan *grpc_connectivity_manager_go.ClusterAlive
	// object types to be considered for consumption
	ToConsume ConsumableStructsInfrastructureEventsConsumer
}

// Create a new configuration structure for a given channel size
// params:
//  size of channels
// return:
//  instance of a configuration object
func NewConfigInfrastructureEventsConsumer(size int, toConsume ConsumableStructsInfrastructureEventsConsumer) ConfigInfrastructureEventsConsumer {
	chUpdateClusterRequest := make(chan *grpc_infrastructure_go.UpdateClusterRequest, size)
	chSetClusterStatusRequest := make(chan *grpc_infrastructure_go.SetClusterStatusRequest, size)
	chClusterAlive := make(chan *grpc_connectivity_manager_go.ClusterAlive, size)

	return ConfigInfrastructureEventsConsumer{
		ChUpdateClusterRequest:    chUpdateClusterRequest,
		ChSetClusterStatusRequest: chSetClusterStatusRequest,
		ChClusterAlive:            chClusterAlive,
		ToConsume:                 toConsume,
	}
}

// Data struct indicating what data structures available in this topic will be accepted.
type ConsumableStructsInfrastructureEventsConsumer struct {
	// Consume update cluster requests
	UpdateClusterRequest bool
	// Consume set cluster status requests
	SetClusterStatusRequest bool
	// Consume cluster alive requests
	ClusterAliveRequest bool
}

func NewInfrastructureEventsConsumer(client bus.NalejClient, name string, exclusive bool, config ConfigInfrastructureEventsConsumer) (*InfrastructureEventsConsumer, derrors.Error) {
	consumer, err := client.BuildConsumer(name, InfrastructureEventsTopic, exclusive)
	if err != nil {
		return nil, err
	}

	return &InfrastructureEventsConsumer{Consumer: consumer, Config: config}, nil
}

// Consume any of the possible objects that can be sent to this queue and send it to the corresponding channel.
func (c InfrastructureEventsConsumer) Consume(ctx context.Context) derrors.Error {
	msg, err := c.Consumer.Receive(ctx)
	if err != nil {
		return err
	}

	target := &grpc_bus_go.InfrastructureEvents{}

	derr := queue.UnmarshallPbMsg(msg, target)
	if derr != nil {
		return derr
	}

	switch x := target.Event.(type) {
	case *grpc_bus_go.InfrastructureEvents_UpdateClusterRequest:
		if c.Config.ToConsume.UpdateClusterRequest {
			c.Config.ChUpdateClusterRequest <- x.UpdateClusterRequest
		}
	case *grpc_bus_go.InfrastructureEvents_SetClusterStatusRequest:
		if c.Config.ToConsume.SetClusterStatusRequest {
			c.Config.ChSetClusterStatusRequest <- x.SetClusterStatusRequest
		}
	case *grpc_bus_go.InfrastructureEvents_ClusterAlive:
		if c.Config.ToConsume.ClusterAliveRequest {
			c.Config.ChClusterAlive <- x.ClusterAlive
		}
	case nil:
		errMsg := "received nil entry"
		log.Error().Msg(errMsg)
		return derrors.NewInvalidArgumentError(errMsg)
	default:
		errMsg := fmt.Sprintf("unknown object type in %s", InfrastructureEventsTopic)
		log.Error().Interface("type", x).Msg(errMsg)
		return derrors.NewInvalidArgumentError(errMsg)
	}
	return nil
}
