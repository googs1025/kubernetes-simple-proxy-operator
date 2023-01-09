## kubernetes-simple-proxy-operator 
## 简易型集群网关operator

### 项目思路与功能
设计背景：集群的网关通常都是采用nginx-controller部署的方式，对自用小集群难免存在部署步骤复杂且资源消耗大等问题。
本项目在此问题上，基于k8s原有的ingress资源上进行简易适配，做出一个有反向代理功能的controller应用。调用方可在cluster中部署与启动相关配置即可使用。
![](https://github.com/googs1025/kubernetes-simple-proxy-operator/blob/main/images/%E6%B5%81%E7%A8%8B%E5%9B%BE%20(2).jpg?raw=true)
思路：当应用启动后，会启动一个controller与proxy反向代理服务，controller会监听ingress资源(annotation需要限定)，并执行相应的业务逻辑。

### 项目部署
1. 进入目录(上传到有集群的机器上 ex:云服务器)

2. 编译镜像
```bigquery
docker run --rm -it -v /root/k8s-operator-proxy:/app -w /app -e GOPROXY=https://goproxy.cn -e CGO_ENABLED=0  golang:1.18.7-alpine3.15 go build -o ./myproxyoperator .
```
3. 部署
```bigquery
# 进入yaml目录
[root@VM-0-16-centos yaml]# kubectl apply -f .
deployment.apps/myproxy-controller created
service/myproxy-svc created
serviceaccount/myproxy-sa created
clusterrole.rbac.authorization.k8s.io/myproxy-clusterrole created
clusterrolebinding.rbac.authorization.k8s.io/myproxy-ClusterRoleBinding created
```
4. 查看应用
```bigquery
[root@VM-0-16-centos ~]# kubectl get pods | grep myproxy
myproxy-controller-74946757f6-x2244    1/1     Running   0                 46m
[root@VM-0-16-centos ~]# kubectl get svc  | grep myproxy
myproxy-svc           NodePort    10.103.46.106   <none>        80:31180/TCP   46m
[root@VM-0-16-centos ~]# kubectl logs myproxy-controller-74946757f6-x2244
I0109 15:32:38.400032       1 init_k8s_config.go:16] run in the cluster
[DBG] "1673278358" sourceOpts="[]" msg="options applied" dst="&{openBalance:false weights:[] addresses:[] tlsConfig:<nil> timeout:0}"
[DBG] "1673278358" dst="&{openBalance:false weights:[] addresses:[] tlsConfig:<nil> timeout:0}" sourceOpts="[]" msg="options applied"
[DBG] "1673278358" dst="&{openBalance:false weights:[] addresses:[] tlsConfig:<nil> timeout:0}" sourceOpts="[]" msg="options applied"
I0109 15:32:38.434310       1 main.go:84] proxy start!!
I0109 15:32:38.434403       1 main.go:76] controller start!!
```
### 项目测试
重要：默认配置文件的修改与注释。
```
server: # server端口
  port: 80
ingress: # 内部默认的ingress配置，用户也可以自己使用kubectl apply -f ingress.yaml任意添加其他ingress资源
  - apiVersion: networking.k8s.io/v1
    kind: Ingress
    metadata:
      name: ingress-myservicea
      annotations:
        # 自定义配置，可以先忽略
        myproxy.ingress.kubernetes.io/add-request-header: name=shenyi;age=19
        myproxy.ingress.kubernetes.io/add-response-header: ret=ok
        myproxy.ingress.kubernetes.io/rewrite-target: /$1
        kubernetes.io/ingress.class: myproxy # annotation 很重要!! controller默认会监听此标签的ingress才进行业务逻辑操作。
    spec:
      rules: # 可以自行配置
        - host: test.jtthink.com
          http:
            paths:
              - path: /baidu/{param:.*}  # 只修改前缀即可
                backend:
                  service:
                    name: www.baidu.com  # 服务名
                    port:
                      number: 80
              - path: /cccccc/{param:.*}
                backend:
                  service:
                    name: cloud.tencent.com
                    port:
                      number: 80
```

#### 增加或删除ingress资源示例
例如：集群中如要添加此ingress资源并加入controller
```bigquery
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-ingress11111
  namespace: default
  annotations:
    kubernetes.io/ingress.class: "myproxy" # 注意这里一定要写死这个annotation
spec:
  rules:
    - host: test.xxxxx.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: nginx
                port:
                  number: 80
```

```bigquery
# 创建
[root@VM-0-16-centos yaml]# kubectl apply -f try_ingress1111.yaml
ingress.networking.k8s.io/my-ingress11111 created
# 可进入container查看
[root@VM-0-16-centos yaml]# kubectl exec -it myproxy-controller-74946757f6-x2244 -- sh
/app # vi app.yaml
```
可以发现如下结果：
```bigquery
Ingress:
- apiVersion: networking.k8s.io/v1
  kind: Ingress
  metadata:
  ...... # 这里都是一样的
- apiVersion: networking.k8s.io/v1
  kind: Ingress
  metadata:
    annotations:
      kubectl.kubernetes.io/last-applied-configuration: |
        {"apiVersion":"networking.k8s.io/v1","kind":"Ingress","metadata":{"annotations":{"kubernetes.io/ingress.class":"jtthink"},"name":"my-ingress11111","namespace":"default"},"spec":{"rules":[{"host":"test.jtthink.com111qqqq","http":{"paths":[{"backend":{"service":{"name":"nginx","port":{"number":80}}},"path":"/","pathType":"Prefix"}]}}]}}
      kubernetes.io/ingress.class: jtthink
    creationTimestamp: "2023-01-09T16:29:15Z"
    generation: 1
    managedFields:
    - apiVersion: networking.k8s.io/v1
      fieldsType: FieldsV1
      fieldsV1:
        f:metadata:
          f:annotations:
            .: {}
            f:kubectl.kubernetes.io/last-applied-configuration: {}
            f:kubernetes.io/ingress.class: {}
        f:spec:
          f:rules: {}
      manager: kubectl-client-side-apply
      operation: Update
      time: "2023-01-09T16:29:15Z"
     name: my-ingress11111 # 这里多加了该配置。
     namespace: default
     resourceVersion: "11180566"
     uid: eeb714bd-b268-47ba-b780-cb664c2ced08
     spec:
       rules:
         - host: test.jtthink.com111qqqq
           http:
           ......          
```
同样，如果删除ingress资源，查看container内app.yaml配置，会发现已经删除。
