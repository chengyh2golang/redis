package job

import (
	"fmt"

	"redis/pkg/apis/crd/v1alpha1"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func New(redis *v1alpha1.Redis)  *batchv1.Job {
	return &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Job",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: redis.Name + "-job-" + RandString(8),
			Namespace: redis.Namespace,
			Labels:    map[string]string{"crd.custom.local": redis.Name},
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(redis, schema.GroupVersionKind{
					Group:   v1alpha1.SchemeGroupVersion.Group,
					Version: v1alpha1.SchemeGroupVersion.Version,
					Kind:    "RedisCluster",
				}),
			},
		},
		Spec:batchv1.JobSpec{
			Template:corev1.PodTemplateSpec{

				Spec:corev1.PodSpec{
					RestartPolicy:corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    "redis-trib-create",
							Image: redis.Spec.RedisTribImage,
							ImagePullPolicy:corev1.PullIfNotPresent,
							Command:[]string{
								"/bin/bash",
								"-c",
								"/tmp/generate-script && /tmp/redis-trib-create.sh",
							},
							Env:[]corev1.EnvVar{
								//通过Sprintf把int32转换成了string
								{Name:"CLUSTER_SIZE",Value:fmt.Sprintf("%v",*redis.Spec.Replicas)},
								{Name:"REDISCLUSTER_NAME",Value:redis.Name},
								{Name:"CLUSTER_OP_TYPE",Value:"create"},
								{Name:"NAMESPACE",Value:redis.Namespace},
							},
						},
					},
				},
			},
		},
	}
}
