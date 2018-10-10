# Exploring running broker on Kubernetes

```plain
kubectl run --generator run-pod/v1 --image cfcommunity/worlds-simplest-service-broker broker
kubectl port-forward broker 3000:3000
```

In another terminal, now connect to broker via local `:3000`:

```plain
$ curl localhost:3000/v2/catalog
{"services":[{"id":"29140B3F-0E69-4C7E-8A35-service-some-service-name","name":"some-service-name","description":"Shared service for some-service-name","bindable":true,"metadata":{"displayName":"some-service-name","imageUrl":"","longDescription":"","providerDisplayName":"","documentationUrl":"","supportUrl":""},"plans":[{"id":"29140B3F-0E69-4C7E-8A35-plan-shared","name":"shared","description":"Shared service for some-service-name","free":true}]}]}
```
