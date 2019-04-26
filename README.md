# Nalej-bus
The Nalej bus is a common events driven message bus for the nalej management cluster.

## Pulsar

Apache Pulsar (http://pulsar.apache.org) is the underlying message queue technology for the Nalej-Bus.
The current K8s solution defined is based on () and deploys the following elements:

* Zookeeper
    * Stateful set with the Zookeeper image
    * The corresponding config map with the Zookeeper configuration
    * A pod disruption budget ensuring at least two running instances
    * Cluster-oriented service
* Apache Book for Pulsar
    * Stateful set running three instances for Apache Book with Pulsar configuration
    * A config map with the corresponding configuration
    * Cluster-oriented service 
* Apache Pulsar broker
    * Pulsar broker with three replicas
    * A config map with the corresponding configuration
    * Cluster-oriented service
* Cluster metadata initialization job
* Cluster topics initialization
    * Create tenants
    * Create namespaces
    
### Caveats

Currently Azure does not permit to set independent volumes for daemon sets. We have modified the proposed
Bookie and Zookeeper solutions to use stateful sets instead. In the case of the Book deployment, a PodDisruptionBudget
should be added to enforce pods to not be deployed into the same physical host.

    
## Topics

The complete installation of the pulsar component generates the following topics:
* infrastructure/ops
* infrastructure/events
* application/ops
* application/events
* network-app/ops
* network-app/events

## Nalej-bus client (Comcast driver)

The Nalej bus client encapsulates basi operations to work with the Nalej bus.
This version uses the Comcast connector https://github.com/Comcast/pulsar-client-go
and it is fully coded in Golang. 

### Executable example

A running example can be compiled to test the production/consumption for an
existing Pular queue. By default the producer send messages to the 
/public/default/topic and the consumer subscribes to it.


```bash
Example of a Nalej client for a Pulsar message queue

Usage:
  nalej-pulsar-cli [command]

Available Commands:
  consumer    run a consumer example
  help        Help about any command
  producer    run a producer example

Flags:
      --consoleLogging         Pretty print logging
      --debug                  enable debugLevel mode
  -h, --help                   help for nalej-pulsar-cli
      --pulsarAddress string   Address for the pulsar server (default "localhost:6650")
      --version                version for nalej-pulsar-cli
```

To simply run a producer use:
```bash
> ./bin/nalej-pulsar-cli producer --consoleLogging
1556264704 INF -> msg="Message number 0"
1556264705 INF -> msg="Message number 1"
1556264706 INF -> msg="Message number 2"
```

In a second terminal run the consumer:
```bash
./bin/nalej-pulsar-cli consumer --consoleLogging
1556264772 INF <- received="Message number 0"
1556264772 INF <- received="Message number 1"
```

### Nalej-bus client example

The following example describes a simple sender/consumer scenario.
```go
import "github.com/nalej/nalej-bus/internal/pulsar-comcast"

// Create a client with the Pulsar address
client := pulsar_comcast.NewClient("localhost:6550")

// Create a producer
producer := pulsar_comcast.NewPulsarProducer(client, producerName, producerTopic)
defer producer.Close()

// Create a consumer. This one is exclusively subscribe to a topic.
consumer := pulsar_comcast.NewPulsarConsumer(client, consumerName, topicConsumer, true)
defer consumer.Close()

// ...

// Wait for messages
go func(){
    for {
        msg, err := consumer.Receive()
        if err != nil {
            log.Error().Err(err)
        } else {
            log.Info().Str("received",string(msg)).Msg("<-")
        }    
    }
}

//...

// Send something
err := producer.Send([]byte("this is a message))
if err != nil {
    log.Error().Err(err).Msg("")
}
```


## Nalej-bus client (Official Pulsar driver)

**This solution is not going to be used right now due to cross compilation issues
Use the comcast driver instead**

The Nalej bus client encapsulates basic operations to work with the nalej bus. The current
version offers an Apache Pulsar wrapper with basic functionality. In order to compile this
client you need the C++ pulsar client library installed in your local environment. For MacOS
run: ```brew install libpulsar```. For other operative system check the instructions at 
https://pulsar.apache.org/docs/en/develop-cpp/

### Nalej-bus client example
The following example describes how to create a Nalej-bus client for a pulsar backend.

```go
import "github.com/nalej-bus/pulsar 

// Create client, only indicate host and port. No protocol is required.
// Additional parameters indicate the timeout in seconds and the number of threads
// used for listening.
client, err := pulsar.NewClient("localhost:6650",5,1)
if err != nil {
    log.Panic("impossible to create client")
}

// Create a producer
prod,err := NewPulsarProducer(client, "prod1", "public/default/topic")
if err != nil {
    log.Panic("impossible to create a producer")
}

// Send something to the queue
msg:="this is a test message"
err = prod.send([]byte(msg))
if err != nil {
    log.Panic("error when sending message")
}

// Create a consumer
cons, err := NewPulsarConsumer(client, "cons1", "public/default/topic", PulsarExclusiveConsumer)

// Block execution flow until we receive something
rec_msg, err := cons.receive()
if err != nil {
    log.Panic("Error when receiving data")
}

log.Info().Msg(string(rec_msg))

// Close connections
prod.Close()
cons.Close()
client.Close()
```

### Considerations
* By default the Nalej Pulsar client acknowledges received messages.
* If using shared subscriptions, the name of the subscription must be the same. 
See https://pulsar.apache.org/docs/en/concepts-messaging/#subscription-modes for more details


## Integration tests

The following set of variables have go be set in order to proceed with integration tests.

| Variable  | Example Value | Description |
| ------------- | ------------- |------------- |
| RUN_INTEGRATION_TEST  | true | Run integration tests |
| IT_PULSAR_ADDRESS | localhost:6650 | Address of an available pulsar message queue |
