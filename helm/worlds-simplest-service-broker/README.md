# World's Simplest Service Broker

If you have a shared service such as Hadoop where all applications will ultimately bind to the same service with the same credentials then you have found the service broker for you - the World's Simplest Service Broker for [Service Catalog](https://svc-cat.io/).

## Introduction

The story of using the Helm chart is broken into two user stories:

1. A service operator who wants to expose its location and credentials to all users of a Kubernetes cluster (or specific namespaces)
    As the service operator you will install this Helm chart and register it with your Service Catalog

1. A user of Kubernetes (such as an application developer) who wants to access the exposed service
    As a Kubernetes user/application developer, you will use Service Catalog CLI `svcat` or the CRDs to provision and bind the shared credentials into your application as Kubernetes secrets.

In this walk through we will share an external SMTP email service with the same credentials for anyone who needs to send emails.

## Prerequisites

* Helm
* [Service Catalog](https://svc-cat.io/) - instructions below

## Installation

As a service operator you will need to install the following systems:

* Service Catalog
* Install each instance of the World's Simplest Service Broker for each remote shared service being shared
* Register each broker with your Service Catalog

### Install Service Catalog

If your Kubernetes does not already have the Service Catalog you can install it/upgrade it with:

```commands
helm repo add svc-cat https://svc-catalog-charts.storage.googleapis.com
helm upgrade --install catalog svc-cat/catalog --namespace catalog
```

Confirm with its `svcat` CLI that there are no initial brokers, but everything is working. It may take a minute or two to start up.

```console
$ svcat get brokers
  NAME   NAMESPACE   URL   STATUS
+------+-----------+-----+--------+
```

### Install World's Simplest Service Broker

You can run the World's Simple Service Broker multiple times - one for each set of shared credentials.

```commands
helm plugin install https://github.com/hypnoglow/helm-s3.git
helm repo add starkandwayne s3://helm.starkandwayne.com/charts
helm repo update

helm upgrade --install email starkandwayne/worlds-simplest-service-broker \
    --wait \
    --set "serviceBroker.class=smtp" \
    --set "serviceBroker.plan=shared" \
    --set "serviceBroker.tags=shared\,email\,smtp" \
    --set "serviceBroker.baseGUID=some-guid" \
    --set "serviceBroker.credentials=\{\"host\":\"mail.authsmtp.com\"\,\"port\":2525\,\"username\":\"ac123456\"\,\"password\":\"special-secret\"\}"
```

Provide any JSON object as the `serviceBroker.credentials` value to the Helm chart.

You can confirm that you've configured the credentials and that the broker is running:

```commands
export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/name=worlds-simplest-service-broker,app.kubernetes.io/instance=email" -o jsonpath="{.items[0].metadata.name}")
kubectl logs $POD_NAME
```

The top of the output will look like:

```console
Each provision/deprovision request will support an async GET /last_operation request
main.serviceBindingResponse{
    Credentials: {
        "host":     "mail.authsmtp.com",
        "port":     float64(2525),
        "username": "ac123456",
        "password": "special-secret",
    },
    SyslogDrainURL: "",
}
Running as http://localhost:3000
[martini] listening on :3000 (development)
[martini] Started GET /v2/catalog for 10.16.1.1:58438
[martini] Completed 200 OK in 317.783µs
```

### Register Service Catalog

The World's Simple Service Broker configured above does not require basic auth; but `svcat register` does require a secret with dummy credentials.

```commands
kubectl create secret generic ignore-basic-auth \
--from-literal username=ignoreme \
--from-literal password=ignoreme

svcat register email-worlds-simplest-service-broker \
--url http://email-worlds-simplest-service-broker.default.svc.cluster.local:3000 \
--scope cluster \
--basic-secret ignore-basic-auth
```

The service broker is now available to `svcat` users in all namespaces across the cluster (`--scope cluster`).

```console
$ svcat get brokers
                  NAME                   NAMESPACE                                      URL                                       STATUS
+--------------------------------------+-----------+----------------------------------------------------------------------------+--------+
  email-worlds-simplest-service-broker               http://email-worlds-simplest-service-broker.default.svc.cluster.local:3000   Ready
```

To view the available `smtp` service class:

```console
$ svcat get classes
  NAME   NAMESPACE         DESCRIPTION
+------+-----------+-------------------------+
  smtp               Shared service for smtp
```

We can see the metadata for our service class with the `-o yaml` output:

```console
$ svcat get class smtp -o yaml
spec:
  bindable: true
  bindingRetrievable: false
  clusterServiceBrokerName: email-worlds-simplest-service-broker
  description: Shared service for smtp
  externalID: some-guid-service-smtp
  externalMetadata:
    displayName: smtp
    documentationUrl: ""
    imageUrl: ""
    longDescription: ""
    providerDisplayName: ""
    supportUrl: ""
  externalName: smtp
  planUpdatable: false
  tags:
  - shared
  - email
  - smtp
...
```

## Demonstration of broker with Service Catalog

As an application developer you can now also discover the available service class (`smtp`) and its service plan (`shared`), and provision/bind its credentials into your namespace as a Secret.

```commands
svcat provision demo --class smtp --plan shared
svcat describe instance demo

svcat bind demo
svcat describe binding demo
```

The output of the `svcat describe binding` will include an indication of the secrets created in the namespace:

```console
Secret Data:
  host       17 bytes
  password   14 bytes
  port       4 bytes
  username   8 bytes
```

The service broker credentials, created during the `svcat bind` command, are stored as a Kubernetes Secret in the namespace under the name of the binding `demo`:

```console
$ kubectl get secret demo -o yaml
apiVersion: v1
data:
  host: bWFpbC5hdXRoc210cC5jb20=
  password: c3BlY2lhbC1zZWNyZXQ=
  port: MjUyNQ==
  username: YWMxMjM0NTY=
kind: Secret
...
```

You can now include this secret within your pods/containers and they will be decoded for you and ready to use.
