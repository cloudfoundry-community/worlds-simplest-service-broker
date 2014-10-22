# World's Simplest Service Broker

If you have a shared service such as Hadoop where all applications will ultimately bind to the same service with the same credentials then you have found the service broker for you - the World's Simplest Service Broker.

You configure it with a simple environment variable `CREDENTIALS` (the same JSON object that will be returned for all service bindings); and then register it as a service broker.

## Build locally

```
godep get
export CREDENTIALS='{"port": 4000}'
export SERVICE_NAME=myservice
export SERVICE_PLAN_NAME=shared
worlds-simplest-service-broker
```

## Deploy to Cloud Foundry

```
cf push myservice-broker --no-start
cf set-env myservice-broker CREDENTIALS '{"port": 4000}'
cf set-env myservice-broker SERVICE_NAME myservice
cf set-env myservice-broker SERVICE_PLAN_NAME shared
cf env myservice-broker
cf start myservice-broker
```
