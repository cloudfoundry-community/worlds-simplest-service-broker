# World's Simplest Service Broker

If you have a shared service such as Hadoop where all applications will ultimately bind to the same service with the same credentials then you have found the service broker for you - the World's Simplest Service Broker for Kubernetes/Service Catalog and for Cloud Foundry.

You configure it with a simple environment variable `CREDENTIALS` (the same JSON object that will be returned for all service bindings); and then register it as a service broker.

All services created and all service bindings will be given the same set of credentials. Definitely the simplest thing that could work.

As a user of the broker in Kubernetes with Service Catalog:

```plain
svcat provision myservice --class some-service-name --plan shared
svcat bind some-service-name
kubectl get secrets some-service-name
```

As a user of the broker in Cloud Foundry:

```plain
cf cs myservice some-service-name
cf bs my-app some-service-name
cf restage my-app
```

## Deploy to Kubernetes & Integrate with Service Catalog

See [README_KUBERNETES.md](README_KUBERNETES.md) for instructions on configuration & deploying the broker to Kubernetes with Helm, plus instructions for integrating and using the broker with Service Catalog.

## Deploy to Cloud Foundry & Integrate with Cloud Foundry

See [README_CLOUDFOUNDRY.md](README_CLOUDFOUNDRY.md) for instructions on configuration & deploying the broker to Cloud Foundry, plus instructions for integrating and using the broker within Cloud Foundry.

## Build and run locally

```plain
export BASE_GUID=$(uuid) # or try $(uuidgen) or any GUID that makes you happy
export CREDENTIALS='{"port": "4000", "host": "1.2.3.4"}'
export SERVICE_NAME=myservice
export SERVICE_PLAN_NAME=shared
export TAGS=simple,shared
go run main.go
```
