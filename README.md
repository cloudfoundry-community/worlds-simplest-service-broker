World's Simplest Service Broker
===============================

If you have a shared service such as Hadoop where all applications will ultimately bind to the same service with the same credentials then you have found the service broker for you - the World's Simplest Service Broker.

You configure it with a simple environment variable `CREDENTIALS` (the same JSON object that will be returned for all service bindings); and then register it as a service broker.

All services created and all service bindings will be given the same set of credentials. Definitely the simplest thing that could work.

As the admin of a service sharing it via a service broker - see section [Deploy to Cloud Foundry](#deploy-to-cloud-foundry) for setup.

As a user of the broker:

```
cf cs myservice some-service-name
cf bs my-app some-service-name
cf restage my-app
```

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
cf push $APPNAME --no-start
cf set-env $APPNAME BASE_GUID $(uuid)
cf set-env $APPNAME CREDENTIALS '{"port": "4000", "host": "1.2.3.4"}'
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

To change the credentials being offered for bindings:

```
export SERVICE=myservice
export APPNAME=$SERVICE-broker
cf set-env $APPNAME CREDENTIALS '{"port": "4000", "host": "1.2.3.4"}'
cf restart $APPNAME
```

Each application will need rebind and restart/restage to get the new credentials.
