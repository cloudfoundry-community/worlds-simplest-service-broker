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

### Why not "user-provided services"?

Cloud Foundry includes "user-provided services" (see `cf cups` in the CLI) for easy registration of existing external service credentials.

One restriction for "cups" is that it is limited to the Space into which it was registered. For each organization/space, the `cf cups` command needs to be run. That is, when you create a new space, it does not immediately have access to the credentials for the service.

The other restriction is that "cups" does not currently support tags.  Frameworks such as [Spring Boot](https://github.com/spring-projects/spring-boot) can leverage tags to inject dependency information into your bound applications.

Instead, with the World's Simplest Service Broker you can make the credentials easily and instantly available to all organizations' spaces.

## Build and run locally

```plain
export BASE_GUID=$(uuid) # or try $(uuidgen) or any GUID that makes you happy
export CREDENTIALS='{"port": "4000", "host": "1.2.3.4"}'
export SERVICE_NAME=myservice
export SERVICE_PLAN_NAME=shared
export TAGS=simple,shared
go run main.go
```
