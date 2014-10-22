# World's Simplest Service Broker

If you have a shared service such as Hadoop where all applications will ultimately bind to the same service with the same credentials then you have found the service broker for you - the World's Simplest Service Broker.

You configure it with a simple environment variable `CREDENTIALS` (the same JSON object that will be returned for all service bindings); and then register it as a service broker.

## Build locally

```
godep get
export BASE_GUID=$(uuid)
export CREDENTIALS='{"port": 4000}'
export SERVICE_NAME=myservice
export SERVICE_PLAN_NAME=shared
worlds-simplest-service-broker
```

## Deploy to Cloud Foundry

```
export SERVICE=myservice
export APPNAME=$SERVICE-broker
cf push $APPNAME --no-start
cf set-env $APPNAME BASE_GUID $(uuid)
cf set-env $APPNAME CREDENTIALS '{"port": 4000}'
cf set-env $APPNAME SERVICE_NAME $SERVICE
cf set-env $APPNAME SERVICE_PLAN_NAME shared
cf env $APPNAME
cf start $APPNAME
```

To register the service broker (as an admin user):

```
cf create-service-broker $SERVICE admin admin https://$APPNAME.gotapaas.com
cf enable-service-access $SERVICE
```
