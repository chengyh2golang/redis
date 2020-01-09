# redis
operator-sdk custom redis cluster

# 使用方法
## 准备redis-trib镜像
1. 准备redis-trib需要用到的基础镜像，基础镜像的Dockerfile内容：
    ```
    FROM docker.io/centos:7
    
    ADD ./redis-trib-rpms.tar.gz /tmp
    COPY ./redis-trib.repo /etc/yum.repos.d/
    
    RUN yum install -y redis-trib bind-utils expect && rm -rf /tmp/redis-trib-rpms && yum clean all
    ```
    
    说明：
    - redis-trib-rpms.tar.gz是安装redis-trib缓存下来的包，通过createrepo做成了一个本地yum源。
    - 基础镜像用到的repo文件redis-trib.repo内容：
    ```
    [redis-trib]
    name=redis-trib
    baseurl=file:///tmp/redis-trib-rpms
    enabled=1
    gpgcheck=0
    
    [epel]
    name=epel
    baseurl=https://mirrors.tuna.tsinghua.edu.cn/epel/7/x86_64
    #mirrorlist=https://mirrors.fedoraproject.org/metalink?repo=epel-7&arch=
    enabled=1
    gpgcheck=0
    ```

2. 构建redis-trib基础镜像
    ```
   docker build -t redis-trib-base:1.0 .
   ```
3. 编译redis-trib需要用到的generate-script脚本文件：
    ```
    go build -o generate-script pkg/resources/utils/redisoperation/gen_redistrib_file.go
    ```
4. 准备redis-trib镜像的Dockerfile
    ```
    FROM redis-trib-base:1.0
    COPY generate-script /tmp
    ```
5. 构建redis-trib镜像
    ```
   docker build -t redis-trib:1.0 .
   ```

    
## 构建operator镜像
```
operator-sdk build redis-operator:1.0
```

## 部署operator
1. 应用crd：
    ```
   kubectl apply -f deploy/crds/crd.custom.local_redis_crd.yaml
   ```
2. 创建serviceAccount：
    ```
   kubectl apply -f deploy/service_account.yaml
   ```
3. 创建role：
    ```
   kubectl apply -f deploy/role.yaml
   ```
4. 创建role_binding：
    ```
   kubectl apply -f deploy/role_binding.yaml
   ```
5. 应用operator：
    ```
   kubectl apply -f deploy/operator.yaml
   ```
6. 创建redis集群：
    ```
   kubectl apply -f deploy/crds/crd.custom.local_v1alpha1_redis_cr.yaml
   ```
   





