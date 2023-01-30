# cisco
Here the instructions to run the operator that deploys ngnix application

## Description


### Running on the cluster
Within your cluster and under git repo:

1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/cisco_v1_ciscocrd.yaml
```
```sh
$ kubectl get CiscoCRD
NAME              AGE
ciscocrd-sample   24h
```

2. Install Nginx Ingress Controller to publish the application

```sh
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/baremetal/deploy.yaml
kubectl patch deployment ingress-nginx-controller -n ingress-nginx -p '{"spec":{"template":{"spec":{"hostNetwork":true}}}}'
```
check that the IP of node is same as the one of controller

```sh
kubectl get nodes -o wide
NAME                 STATUS   ROLES           AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION      CONTAINER-RUNTIME
kind-control-plane   Ready    control-plane   25h   v1.25.3   172.28.0.2    <none>        Ubuntu 22.04.1 LTS   5.15.0-56-generic   containerd://1.6.9

kubectl get pod -n ingress-nginx -o wide
NAME                                        READY   STATUS      RESTARTS   AGE   IP           NODE                 NOMINATED NODE   READINESS GATES
ingress-nginx-controller-787db7674b-trcdm   1/1     Running     0          25h   172.28.0.2   kind-control-plane   <none>           <none>
```

3. Install cert-manager:

```sh
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml
```
check that all the pods are running

```sh
kubectl -n cert-manager get po
NAME                                      READY   STATUS    RESTARTS   AGE
cert-manager-99bb69456-4v74x              1/1     Running   0          25h
cert-manager-cainjector-ffb4747bb-4tpbr   1/1     Running   0          25h
cert-manager-webhook-545bd5d7d8-44rr6     1/1     Running   0          25h
```

4. Install certificate issuer:

```sh
kubectl apply -f config/samples/letsencrypt-cluster-issuer.yaml
```
Check that's ready :

```sh
kubectl get ClusterIssuer -n cert-manager
NAME                         READY   AGE
letsencrypt-cluster-issuer   True    25h
```

5. Install the certificate:
```sh
kubectl apply -f config/samples/certificate.yaml
```
Important : the DNS name specified in the certificate should be resolvable via DNS Server otherwise this step will fail

Once the certificate has been issued you should see it in Kubernetes secrets.
```sh
kubectl get secrets
```

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```
all the resources should be deployed :

```sh
kubectl get all
NAME                                   READY   STATUS    RESTARTS   AGE
pod/ciscocrd-sample-5799c9b9cd-8tcnk   1/1     Running   0          25h
pod/ciscocrd-sample-5799c9b9cd-9mtlk   1/1     Running   0          25h

NAME                      TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
service/ciscocrd-sample   ClusterIP   10.96.104.87   <none>        80/TCP,443/TCP   25h
service/kubernetes        ClusterIP   10.96.0.1      <none>        443/TCP          25h

NAME                              READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/ciscocrd-sample   2/2     2            2           25h

NAME                                         DESIRED   CURRENT   READY   AGE
replicaset.apps/ciscocrd-sample-5799c9b9cd   2         2         2       25h

kubectl get ingress
NAME                      CLASS   HOSTS             ADDRESS      PORTS     AGE
ciscocrd-sample-ingress   nginx   cisco.local.com   172.28.0.2   80, 443   22h
```
Important : To be able to access the service from outside the cluster the service should be deployed as a Load Balancer, in this case it will have an assigned External_IP adress
