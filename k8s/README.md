# Exploring running broker on Kubernetes

## Namespace

```plain
kubectl create namespace broker-demo
```

## Single Pod

```plain
kubectl run -n broker-demo --generator run-pod/v1 --image cfcommunity/worlds-simplest-service-broker broker
kubectl -n broker-demo port-forward broker 3000:3000
```

In another terminal, now connect to broker via local `:3000`:

```plain
$ curl localhost:3000/v2/catalog
{"services":[{"id":"29140B3F-0E69-4C7E-8A35-service-some-service-name","name":"some-service-name","description":"Shared service for some-service-name","bindable":true,"metadata":{"displayName":"some-service-name","imageUrl":"","longDescription":"","providerDisplayName":"","documentationUrl":"","supportUrl":""},"plans":[{"id":"29140B3F-0E69-4C7E-8A35-plan-shared","name":"shared","description":"Shared service for some-service-name","free":true}]}]}
```

To delete:

```plain
kubectl -n broker-demo delete pod broker
```

## Deployment

```plain
kubectl -n broker-demo run --image cfcommunity/worlds-simplest-service-broker broker --replicas=1 -o yaml --dry-run
kubectl -n broker-demo run --image cfcommunity/worlds-simplest-service-broker broker --replicas=1
```

```plain
$ kubectl get pods -n broker-demo
NAME                      READY   STATUS    RESTARTS   AGE
broker-6c574899fc-9f7gz   1/1     Running   0          1m
$ kubectl get deployment -n broker-demo
NAME     DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
broker   1         1         1            1           1m
```

```plain
kubectl expose deployment -n broker-demo broker --type ClusterIP --port 3000 --target-port 3000
```

```plain
$ kubectl get services -n broker-demo
NAME     TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
broker   ClusterIP   10.106.138.185   <none>        3000/TCP   53s
```

To port foward to the exposed service on port `:3000`:

```plain
kubectl -n broker-demo port-forward service/broker 3000:3000
```

```plain
curl localhost:3000/v2/catalog
```

Explore with `eden`, first setup env vars:

```plain
export SB_BROKER_URL=http://localhost:3000
export SB_BROKER_USERNAME=
export SB_BROKER_PASSWORD=
```

Next interact with broker:

```plain
$ eden catalog
Service Name       Plan Name  Description
some-service-name  shared     Shared service for some-service-name
```

## Deploy from YAML

```plain
kubectl delete namespaces broker-demo
kubectl create namespace broker-demo
kubectl apply -n broker-demo -f k8s/broker-demo.yaml
```

Stop and restart the port-forward tunnel:

```plain
kubectl -n broker-demo port-forward service/broker 3000:3000
```

In the `eden` terminal, confirm it continues to work:

```plain
$ eden catalog
Service Name       Plan Name  Description
some-service-name  shared     Shared service for some-service-name
```
