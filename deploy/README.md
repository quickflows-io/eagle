## deploy

Mainly store some deployment-related configuration files and scripts

### Deploy the Go application

> see: https://eddycjy.com/posts/kubernetes/2020-05-03-deployment/


## monitor

It includes monitoring of machines (nodes or containers), monitoring of applications, monitoring of databases, etc.

Using `docker-compose` can be deployed locally with one click, the configuration is as follows:

```yaml

```

## configure etcd

create a namespace

`$ kubectl create namespace etcd`

create the service 

`$ kubectl apply -f etcd-service.yaml -n etcd`

create the cluster(statefulSet)

`$ cat etcd.yml.tmpl | etcd_config.bash | kubectl apply -n etcd -f -`

Verify the cluster's health

`$ kubectl exec -it etcd-0 -n etcd etcdctl cluster-health`

The cluster is exposed through minikube's IP

```bash
$ IP=$(minikube ip)
$ PORT=$(kubectl get services -o jsonpath="{.spec.ports[].nodePort}" etcd-client -n etcd)
$ etcdctl --endpoints http://${IP}:${PORT} get foo
```

Destroy the services

```bash
$ kubectl delete services,statefulsets --all -n etcd
```

Make a Web UI for the Etcd cluster

```
$ kubectl apply -f etcd-ui-configmap.yaml
$ kubectl apply -f etcd-ui.yaml
```

> Ref: https://github.com/kevinyan815/LearningKubernetes/tree/master/e3w

> refer to： 
> https://mp.weixin.qq.com/s/AkIvkW22dvqcdFXkiTpv8Q
> https://github.com/kevinyan815/LearningKubernetes/tree/master/etcd