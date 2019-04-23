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



