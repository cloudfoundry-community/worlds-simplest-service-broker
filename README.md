World's Simplest Service Broker
===============================

If you have a shared service such as Hadoop where all applications will ultimately bind to the same service with the same credentials then you have found the service broker for you - the World's Simplest Service Broker for Cloud Foundry.

You configure it with a simple environment variable `CREDENTIALS` (the same JSON object that will be returned for all service bindings); and then register it as a service broker.

All services created and all service bindings will be given the same set of credentials. Definitely the simplest thing that could work.

As the admin of a service sharing it via a service broker - see section [Deploy to Cloud Foundry](#deploy-to-cloud-foundry) for setup.

As a user of the broker:

```
cf cs myservice some-service-name
cf bs my-app some-service-name
cf restage my-app
```

Why not "user-provided services"?
---------------------------------

Cloud Foundry includes "user-provided services" (see `cf cups` in the CLI) for easy registration of existing external service credentials.

The restriction for "cups" is that it is limited to the Space into which it was registered. For each organization/space, the `cf cups` command needs to be run. That is, when you create a new space, it does not immediately have access to the credentials for the service.

Instead, with the World's Simplest Service Broker you can make the credentials easily and instantly available to all organizations' spaces.

Build locally
-------------

```
godep get
export BASE_GUID=$(uuid)
export CREDENTIALS='{"port": "4000", "host": "1.2.3.4"}'
export SERVICE_NAME=myservice
export SERVICE_PLAN_NAME=shared
worlds-simplest-service-broker
```

Deploy to Cloud Foundry
-----------------------

```
export SERVICE=myservice
export APPNAME=$SERVICE-broker
cf push $APPNAME --no-start -m 128M -k 256M
cf set-env $APPNAME BASE_GUID $(uuid)
cf set-env $APPNAME CREDENTIALS '{"port": "4000", "host": "1.2.3.4"}'
cf set-env $APPNAME SERVICE_NAME $SERVICE
cf set-env $APPNAME SERVICE_PLAN_NAME shared
cf env $APPNAME
cf start $APPNAME
```

To register the service broker (as an admin user):

```
export SERVICE_URL=$(cf app $APPNAME | grep urls: | awk '{print $2}')
cf create-service-broker $SERVICE admin admin https://$SERVICE_URL
cf enable-service-access $SERVICE
```

To change the credentials being offered for bindings:

```
export SERVICE=myservice
export APPNAME=$SERVICE-broker
cf set-env $APPNAME CREDENTIALS '{"port": "4000", "host": "1.2.3.4"}'
cf restart $APPNAME
```

Each application will need rebind and restart/restage to get the new credentials.

### Dashboard

Each service instance is assigned the same dashboard URL - `/dashboard`.
