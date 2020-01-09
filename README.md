# redis
operator-sdk custom redis cluster

#使用方法

## 首先构建出redis-trib镜像
编译pkg/resources/utils/redisoperation目录下的gen_redistrib_file.go：
```
go build -o generate-script gen_redistrib_file.go
```



