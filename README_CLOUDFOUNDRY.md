# Deploy into Cloud Foundry and integrate with Cloud Foundry

You can both deploy the broker to Cloud Foundry, and then turn around and offer it as a service broker to all other users of that Cloud Foundry (or a different set of Cloud Foundries):

```plain
export SERVICE=myservice
export APPNAME=$SERVICE-broker
cf push $APPNAME --no-start -m 128M -k 256M
cf set-env $APPNAME BASE_GUID $(uuid) # or try $(uuidgen) or any GUID that makes you happy
cf set-env $APPNAME CREDENTIALS '{"port": "4000", "host": "1.2.3.4"}'
cf set-env $APPNAME SERVICE_NAME $SERVICE
cf set-env $APPNAME SERVICE_PLAN_NAME shared
cf set-env $APPNAME TAGS simple,shared
cf set-env $APPNAME AUTH_USER broker
cf set-env $APPNAME AUTH_PASSWORD broker
cf env $APPNAME
cf start $APPNAME
```

To register the service broker (as an admin user):

```plain
export SERVICE_URL=$(cf app $APPNAME | grep routes: | awk '{print $2}')
cf create-service-broker $SERVICE broker broker https://$SERVICE_URL
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

## Basic Authentication

Set the environment variables `AUTH_USER` and `AUTH_PASSWORD` to enable basic authentication for the broker.

This is useful to prevent unauthorized access to the credentials exposed by the broker (e.g. by somebody doing a `curl -X PUT http://$SERVICE_URL/v2/service_instances/a/service_bindings/b`).

To do so (of course, change `secret_user` and `secret_password` to something more secret):

```plain
cf set-env $APPNAME AUTH_USER secret_user
cf set-env $APPNAME AUTH_PASSWORD secret_password
cf restart $APPNAME
cf update-service-broker $SERVICE broker broker https://$SERVICE_URL
```

## syslog_drain_url

The broker can advertise a `syslog_drain_url` endpoint with the `$SYSLOG_DRAIN_URL` variable:

```plain
cf set-env $APPNAME SYSLOG_DRAIN_URL 'syslog://1.2.3.4:514'
```

## Dashboard

Each service instance is assigned the same dashboard URL - `/dashboard`.

## Image URL

Adding image URL to service broker for marketplace advertisements.

```plain
cf set-env $APPNAME IMAGE_URL '<image url>'
```

