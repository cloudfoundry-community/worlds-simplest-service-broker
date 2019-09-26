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

See [helm/worlds-simplest-service-broker/README.md](helm/worlds-simplest-service-broker/README.md) for instructions on configuration & deploying the broker to Kubernetes with Helm, plus instructions for integrating and using the broker with Service Catalog.

## Deploy to Cloud Foundry & Integrate with Cloud Foundry

See [README_CLOUDFOUNDRY.md](README_CLOUDFOUNDRY.md) for instructions on configuration & deploying the broker to Cloud Foundry, plus instructions for integrating and using the broker within Cloud Foundry.

## Build and run locally

```shell
export BASE_GUID=$(uuid)
# or
export BASE_GUID=$(uuidgen)

export CREDENTIALS='{"port": "4000", "host": "1.2.3.4"}'
export SERVICE_NAME=myservice
export SERVICE_PLAN_NAME=shared
export TAGS=simple,shared
export AUTH_USER=broker
export AUTH_PASSWORD=broker
go run cmd/worlds-simplest-service-broker/main.go
```

## Docker

Below are sections on building and running with OCI/Docker.

### Cloud Native Buildpacks

You can compile this project and produce an OCI/Docker image with [Cloud Native Buildpacks](https://buildpacks.io/) [`pack` CLI](https://buildpacks.io/docs/install-pack/):

```plain
pack build cfcommunity/worlds-simplest-service-broker --builder cloudfoundry/cnb:tiny
```

NOTE: the [CI pipeline](https://ci2.starkandwayne.com/teams/starkandwayne/pipelines/worlds-simplest-service-broker/jobs/latest-image/builds/5) uses this method to create OCI/docker images.

Whilst the resulting OCI is larger (`docker images` says `39MB` vs `19MB` for `docker build`), the image includes a lot of metadata about the tools used to build the application, and supports the `pack rebase` command to allow operators to update the base `run` image over time without ever needing the OCI to be rebuilt or redistributed. Read https://buildpacks.io/ for more.

### Docker Build

The project also contains a `Dockerfile` for `docker build`:

```plain
docker build -t cfcommunity/worlds-simplest-service-broker .
```

### Docker Run

Either create the OCI/Docker image as above, or pull down from Docker Hub, and then run using the same env vars from above:

```plain
docker run -e BASE_GUID=$BASE_GUID \
    -e CREDENTIALS=$CREDENTIALS \
    -e SERVICE_NAME=$SERVICE_NAME \
    -e SERVICE_PLAN_NAME=$SERVICE_PLAN_NAME \
    -e TAGS=$TAGS \
    -e AUTH_USER=broker -e AUTH_PASSWORD=broker \
    -p 3000:3000 cfcommunity/worlds-simplest-service-broker
```
