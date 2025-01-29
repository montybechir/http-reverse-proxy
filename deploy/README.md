## Overview

Deploy using Kubernetes for high availability and reliability.

### Prereqs

1. Create a container registry.

2. Create an AKS cluster via the portal or CLI.

```bash

az aks create --name aks-test --resource-group aks --generate-ssh-keys

```

3. Get credentials for the cluster.

```bash

az aks get-credentials --name aks-test --resource-group aks

```

### Build and push images to registry

```bash
az acr build --image reverseproxy/webapp:a --registry mbnregistry --file build/package/backenda/Dockerfile  .

az acr build --image reverseproxy/webapp:b --registry mbnregistry --file build/package/backendb/Dockerfile  .

az acr build --image reverseproxy/reverseproxy:1 --registry mbnregistry --file build/package/proxy/Dockerfile  .
```

### Create services and deployments

```bash
MB-Machine:http-reverse-proxy montasir$ k apply -f deploy/k8s/webapp.yaml

MB-Machine:http-reverse-proxy montasir$ k apply -f deploy/k8s/proxy-config.yaml

MB-Machine:http-reverse-proxy montasir$ k apply -f deploy/k8s/reverse-proxy.yaml

```

### Attach ACR to K8S

Check events for potential failures.

```bash

MB-Machine:http-reverse-proxy montasir$ k events
...
14s (x2 over 29s)   Warning   Failed                    Pod/webappb-deployment-6d56556f98-rtvxm    Error: ErrImagePull
13s (x2 over 29s)   Warning   Failed                    Pod/webappa-deployment-569469947b-zq599    Failed to pull image "mbnregistry.azurecr.io/reverseproxy/webapp:a": failed to pull and unpack image "mbnregistry.azurecr.io/reverseproxy/webapp:a": failed to resolve reference "mbnregistry.azurecr.io/reverseproxy/webapp:a": failed to authorize: failed to fetch anonymous token: unexpected status from GET request to https://mbnregistry.azurecr.io/oauth2/token?scope=repository%3Areverseproxy%2Fwebapp%3Apull&service=mbnregistry.azurecr.io: 401 Unauthorized
13s (x2 over 30s)   Normal    Pulling                   Pod/webappa-deployment-569469947b-zq599    Pulling image "mbnregistry.azurecr.io/reverseproxy/webapp:a"
13s (x2 over 29s)   Warning   Failed                    Pod/webappa-deployment-569469947b-zq599    Error: ErrImagePull
6s (x2 over 29s)    Warning   Failed                    Pod/webappa-deployment-569469947b-zvld8    Error: ImagePullBackOff
6s (x2 over 29s)    Normal    BackOff                   Pod/webappa-deployment-569469947b-zvld8    Back-off pulling image "mbnregistry.azurecr.io/reverseproxy/webapp:a"
```

Fix:

- Option 1: Update IAM by enabling AcrPull permissions on the aks-test-agentpool user-asssigned managed identity.
- Option 2: Attach the container registry to the cluster.

```bash

az aks update --name aks-test --resource-group aks --attach-acr mbnregistry

```

### Check Statuses

```bash

MB-Machine:http-reverse-proxy montasir$ k get all
NAME                                           READY   STATUS    RESTARTS   AGE
pod/reverseproxy-deployment-78c9c69f9f-9mtsc   1/1     Running   0          93s
pod/reverseproxy-deployment-78c9c69f9f-mtbpp   1/1     Running   0          93s
pod/webappa-deployment-569469947b-6x85j        1/1     Running   0          107s
pod/webappa-deployment-569469947b-gf5pw        1/1     Running   0          107s
pod/webappb-deployment-6d56556f98-5m45m        1/1     Running   0          107s
pod/webappb-deployment-6d56556f98-99h9v        1/1     Running   0          107s

NAME                           TYPE           CLUSTER-IP     EXTERNAL-IP    PORT(S)        AGE
service/kubernetes             ClusterIP      10.0.0.1       <none>         443/TCP        39m
service/reverseproxy-service   LoadBalancer   10.0.131.173   4.174.197.17   80:31425/TCP   93s
service/webappa-service        ClusterIP      10.0.38.115    <none>         60408/TCP      107s
service/webappb-service        ClusterIP      10.0.105.6     <none>         60409/TCP      107s

NAME                                      READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/reverseproxy-deployment   2/2     2            2           94s
deployment.apps/webappa-deployment        2/2     2            2           108s
deployment.apps/webappb-deployment        2/2     2            2           108s

NAME                                                 DESIRED   CURRENT   READY   AGE
replicaset.apps/reverseproxy-deployment-78c9c69f9f   2         2         2       94s
replicaset.apps/webappa-deployment-569469947b        2         2         2       108s
replicaset.apps/webappb-deployment-6d56556f98        2         2         2       108s

```

### Test

```bash

MB-Machine:http-reverse-proxy montasir$ curl http://4.174.197.17 -v
*   Trying 4.174.197.17:80...
* Connected to 4.174.197.17 (4.174.197.17) port 80
> GET / HTTP/1.1
> Host: 4.174.197.17
> User-Agent: curl/8.7.1
> Accept: */*
>
* Request completely sent off
< HTTP/1.1 200 OK
< Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With
< Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
< Access-Control-Allow-Origin:
< Content-Length: 23
< Content-Type: text/plain; charset=utf-8
< Date: Wed, 29 Jan 2025 21:25:11 GMT
<
* Connection #0 to host 4.174.197.17 left intact
Response from Backend AMB-Machine:http-reverse-proxy montasir$ curl http://4.174.197.17 -v
*   Trying 4.174.197.17:80...
* Connected to 4.174.197.17 (4.174.197.17) port 80
> GET / HTTP/1.1
> Host: 4.174.197.17
> User-Agent: curl/8.7.1
> Accept: */*
>
* Request completely sent off
< HTTP/1.1 200 OK
< Access-Control-Allow-Headers: Content-Type, Authorization, X-Requested-With
< Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
< Access-Control-Allow-Origin:
< Content-Length: 23
< Content-Type: text/plain; charset=utf-8
< Date: Wed, 29 Jan 2025 21:25:14 GMT
<
* Connection #0 to host 4.174.197.17 left intact
```

### Clean Up

Delete the resources

```bash
MB-Machine:http-reverse-proxy montasir$ k delete -f deploy/k8s/webapp.yaml
deployment.apps "webappa-deployment" deleted
service "webappa-service" deleted
deployment.apps "webappb-deployment" deleted
service "webappb-service" deleted


MB-Machine:http-reverse-proxy montasir$  k delete -f deploy/k8s/proxy-config.yaml
configmap "proxy-config" deleted
MB-Machine:http-reverse-proxy montasir$

MB-Machine:http-reverse-proxy montasir$ k delete -f deploy/k8s/reverse-proxy.yaml
deployment.apps "reverseproxy-deployment" deleted
service "reverseproxy-service" deleted


```

Stop or delete the cluster.

```bash
MB-Machine:http-reverse-proxy montasir$ az aks delete --name aks-test --resource-group aks --yes --no-wait
```
