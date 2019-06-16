## Running open-falcon in kubernetes

work on kubernetes 1.14 ,  the `kubectl version` like:
```
Client Version: version.Info{Major:"1", Minor:"14", GitVersion:"v1.14.0", GitCommit:"641856db18352033a0d96dbc99153fa3b27298e5", GitTreeState:"clean", BuildDate:"2019-03-25T15:53:57Z", GoVersion:"go1.12.1", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"14", GitVersion:"v1.14.0", GitCommit:"641856db18352033a0d96dbc99153fa3b27298e5", GitTreeState:"clean", BuildDate:"2019-03-25T15:45:25Z", GoVersion:"go1.12.1", Compiler:"gc", Platform:"linux/amd64"}
```

##### 1. Start mysql in k8s and init the mysql table before the first running

if mysql is already in k8s you can break this step

```
kubectl apply -f mysql.yaml
```

init mysql table before the first running

```
sh init_mysql_data.sh
```

##### 2. Start redis in k8s

if redis is already in k8s you can also break this step

```
kubectl apply -f redis.yaml
```

##### 3. Start falcon-plus modules in one pod

```
kubectl apply -f openfalcon-plus.yaml
```

##### 4. Start falcon-dashboard in k8s

```
kubectl apply -f openfalcon-dashboard.yaml
```

##### 5. browse the dashboard view

```
[tyhall51@192-168-10-21 k8s-example]$ kubectl get svc
NAME                    TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                         AGE
kubernetes              ClusterIP   10.96.0.1       <none>        443/TCP                         105d
mysql                   NodePort    10.110.20.201   <none>        3306:30535/TCP                  25m
open-falcon             NodePort    10.97.12.125    <none>        8433:31952/TCP,8080:31957/TCP   53s
open-falcon-dashboard   NodePort    10.96.119.231   <none>        8081:30191/TCP                  3s
redis                   ClusterIP   10.98.212.126   <none>        6379/TCP                        32m
```

get **open-falcon-dashboard** service localhost port **30191** , then can visit  `http://192.168.10.21:30191` in webrowser 。


[mailto](mailto:studyoo@foxmail.com)

