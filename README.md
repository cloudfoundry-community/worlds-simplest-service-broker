# World's Simplest Service Broker

If you have a shared service such as Hadoop where all applications will ultimately bind to the same service with the same credentials then you have found the service broker for you - the World's Simplest Service Broker for Cloud Foundry.

You configure it with a simple environment variable `CREDENTIALS` (the same JSON object that will be returned for all service bindings); and then register it as a service broker.

All services created and all service bindings will be given the same set of credentials. Definitely the simplest thing that could work.

As the admin of a service sharing it via a service broker - see section [Deploy to Cloud Foundry](#deploy-to-cloud-foundry) for setup.

As a user of the broker:

```plain
cf cs myservice some-service-name
cf bs my-app some-service-name
cf restage my-app
```

## Why not "user-provided services"?

Cloud Foundry includes "user-provided services" (see `cf cups` in the CLI) for easy registration of existing external service credentials.

One restriction for "cups" is that it is limited to the Space into which it was registered. For each organization/space, the `cf cups` command needs to be run. That is, when you create a new space, it does not immediately have access to the credentials for the service.  

The other restriction is that "cups" does not currently support tags.  Frameworks such as [Spring Boot](https://github.com/spring-projects/spring-boot) can leverage tags to inject dependency information into your bound applications.

Instead, with the World's Simplest Service Broker you can make the credentials easily and instantly available to all organizations' spaces.

## Build locally

```plain
export BASE_GUID=$(uuid) # or try $(uuidgen) or any GUID that makes you happy
export CREDENTIALS='{"port": "4000", "host": "1.2.3.4"}'
export SERVICE_NAME=myservice
export SERVICE_PLAN_NAME=shared
export TAGS=simple,shared
go run main.go
```

## Deploy to Cloud Foundry

```plain
export SERVICE=myservice
export APPNAME=$SERVICE-broker
cf push $APPNAME --no-start -m 128M -k 256M
cf set-env $APPNAME BASE_GUID $(uuid)
cf set-env $APPNAME CREDENTIALS '{"port": "4000", "host": "1.2.3.4"}'
cf set-env $APPNAME SERVICE_NAME $SERVICE
cf set-env $APPNAME SERVICE_PLAN_NAME shared
cf set-env $APPNAME TAGS simple,shared
cf env $APPNAME
cf start $APPNAME
```

To register the service broker (as an admin user):

```plain
export SERVICE_URL=$(cf app $APPNAME | grep routes: | awk '{print $2}')
cf create-service-broker $SERVICE admin admin https://$SERVICE_URL
cf enable-service-access $SERVICE
```

To change the credentials being offered for bindings:

```plain
export SERVICE=myservice
export APPNAME=$SERVICE-broker
cf set-env $APPNAME CREDENTIALS '{"port": "4000", "host": "1.2.3.4"}'
cf restart $APPNAME
```

Each application will need rebind and restart/restage to get the new credentials.

### Basic Authentication

Set the environment variables `AUTH_USER` and `AUTH_PASSWORD` to enable basic authentication for the broker.

This is useful to prevent unauthorized access to the credentials exposed by the broker (e.g. by somebody doing a `curl -X PUT http://$SERVICE_URL/v2/service_instances/a/service_bindings/b`).

To do so (of course, change `secret_user` and `secret_password` to something more secret):

```plain
cf set-env $APPNAME AUTH_USER secret_user
cf set-env $APPNAME AUTH_PASSWORD secret_password
cf restart $APPNAME
cf update-service-broker $SERVICE secret_user secret_password https://$SERVICE_URL
```

### syslog_drain_url

The broker can advertise a `syslog_drain_url` endpoint with the `$SYSLOG_DRAIN_URL` variable:

```plain
cf set-env $APPNAME SYSLOG_DRAIN_URL 'syslog://1.2.3.4:514'
```

### Dashboard

Each service instance is assigned the same dashboard URL - `/dashboard`.

### Image URL

Adding image url to service broker

```plain
cf set-env $APPNAME IMAGE_URL '<image url>'
```
