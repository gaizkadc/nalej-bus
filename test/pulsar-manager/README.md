# Using pulsar manager


1. To install the manager, create the `pulsar-manager.*.yaml` files on the target cluster

```
$ kubectl -nnalej create -f pulsar-manager.deployment.yaml
$ kubectl -nnalej create -f pulsar-manager.service.yaml
```

2. Forward the port to the UI

```
$ kubectl -nnalej port-forward service/pulsar-manager 9527:9527
```

3. Connect to http://localhost:9527 using pulsar:pulsar

4. Create a new environment targeting http://broker.nalej:8080

5. Enjoy
