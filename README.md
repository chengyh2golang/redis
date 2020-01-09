# redis
operator-sdk custom redis cluster

# 使用方法
## 准备redis-trib镜像

* 准备redis-trib需要用到的基础镜像
```
FROM docker.io/centos:7

ADD ./redis-trib-rpms.tar.gz /tmp
COPY ./redis-trib.repo /etc/yum.repos.d/

RUN yum install -y redis-trib bind-utils expect && rm -rf /tmp/redis-trib-rpms && yum clean all
```

说明：
redis-trib-rpms.tar.gz是安装redis-trib缓存下来的包，通过createrepo做成了一个本地yum源。基础镜像用的repo文件内容：
redis-trib.repo文件内容：
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

* 构建出redis-trib镜像
编译pkg/resources/utils/redisoperation目录下的gen_redistrib_file.go：
```
go build -o generate-script gen_redistrib_file.go
```

然后基于redis-trib的基础镜像，制作redis-trib的应用镜像，该镜像主要被job调用，用于创建、扩容和缩容redis集群。
```
FROM redis-trib-base:1.0
COPY generate-script /tmp
```




