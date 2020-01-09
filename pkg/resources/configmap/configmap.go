package configmap

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"redis/pkg/apis/crd/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	RedisConfigKey          = "redis.conf"
	FixIPKey = "fix-ip.sh"
	//RedisConfigRelativePath = "redis.conf"
)

var redisConfig = `cluster-enabled yes
cluster-config-file /data/nodes.conf
cluster-node-timeout 10000
protected-mode no
daemonize no
pidfile /var/run/redis.pid
port 6379
tcp-backlog 511
bind 0.0.0.0
timeout 3600
tcp-keepalive 1
loglevel verbose
logfile /data/redis.log
databases 16
save 900 1
save 300 10
save 60 10000
stop-writes-on-bgsave-error yes
rdbcompression yes
rdbchecksum yes
dbfilename dump.rdb
dir /data
#requirepass yl123456
appendonly yes
appendfilename "appendonly.aof"
appendfsync everysec
no-appendfsync-on-rewrite no
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
lua-time-limit 20000
slowlog-log-slower-than 10000
slowlog-max-len 128
#rename-command FLUSHALL  ""
latency-monitor-threshold 0
notify-keyspace-events ""
hash-max-ziplist-entries 512
hash-max-ziplist-value 64
list-max-ziplist-entries 512
list-max-ziplist-value 64
set-max-intset-entries 512
zset-max-ziplist-entries 128
zset-max-ziplist-value 64
hll-sparse-max-bytes 3000
activerehashing yes
client-output-buffer-limit normal 0 0 0
client-output-buffer-limit slave 256mb 64mb 60
client-output-buffer-limit pubsub 32mb 8mb 60
hz 10
aof-rewrite-incremental-fsync yes
`


var fixIPConfig = `#!/bin/sh
    CLUSTER_CONFIG="/data/nodes.conf"
    if [ -f ${CLUSTER_CONFIG} ]; then
      if [ -z "${POD_IP}" ]; then
        echo "Unable to determine Pod IP address!"
        exit 1
      fi
      echo "Updating my IP to ${POD_IP} in ${CLUSTER_CONFIG}"
      sed -i.bak -e '/myself/ s/[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}\.[0-9]\{1,3\}/'${POD_IP}'/' ${CLUSTER_CONFIG}
    fi
    exec "$@"
`


func New(redis *v1alpha1.Redis) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		TypeMeta:metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta:metav1.ObjectMeta{
			Name:redis.Name,
			Namespace:redis.Namespace,
			OwnerReferences:[]metav1.OwnerReference{
				*metav1.NewControllerRef(redis,schema.GroupVersionKind{
					Group:v1alpha1.SchemeGroupVersion.Group,
					Version:v1alpha1.SchemeGroupVersion.Version,
					Kind:"Redis",
				}),
			},
			Labels:map[string]string{"crd.custom.local": redis.Name},
		},
		Data: map[string]string{
			RedisConfigKey:redisConfig,
			FixIPKey:fixIPConfig,
		},
	}
}

