# Deploy with Helm and install into Service Catalog

The World's Simplest Service Broker is packaged as a container image and installable into your Kubernetes cluster as a [Helm](https://helm.sh/) chart.

It is also installable into your Kubernetes [Service Catalog](https://svc-cat.io/).

## Install Service Catalog

```commands
helm repo add svc-cat https://svc-catalog-charts.storage.googleapis.com
helm install svc-cat/catalog --name catalog --namespace catalog
```

Confirm with its `svcat` CLI that there are no initial brokers, but everything is working. It may take a minute or two to start up.

```console
$ svcat get brokers
  NAME   NAMESPACE   URL   STATUS
+------+-----------+-----+--------+
```

## Install World's Simplest Service Broker

You can run the World's Simple Service Broker multiple times - one for each set of shared credentials.

```commands
helm install ./helm --name email --wait \
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

## Register Service Catalog

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
