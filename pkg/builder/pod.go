package builder

import (
	"context"
	"strconv"
	"strings"

	taskv1alpha1 "github.com/cleverhu/task-cicd/pkg/apis/task/v1alpha1"

	"github.com/google/go-containerregistry/pkg/name"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	PodPrefix           = "task-pod-"
	TaskOrderAnnotation = "taskorder"
)

type PodBuilder struct {
	task *taskv1alpha1.Task // 任务对象
	client.Client
}

//构造函数
func NewPodBuilder(client client.Client, task *taskv1alpha1.Task) *PodBuilder {
	return &PodBuilder{task: task, Client: client}
}

//构建 创建出 对应的POD
func (pb *PodBuilder) Build(ctx context.Context) error {
	pod := &corev1.Pod{}
	err := pb.Get(ctx, types.NamespacedName{Namespace: pb.task.Namespace, Name: PodPrefix + pb.task.Name}, pod)
	if apierrors.IsNotFound(err) {
		pod.Spec.RestartPolicy = corev1.RestartPolicyNever //从不 重启
		pb.setPodMetadata(pod)
		pb.setPodVolumes(pod)
		pb.setInitContainer(pod)
		containers := []corev1.Container{} // 容器切片
		for i, step := range pb.task.Spec.Steps {
			container, err := pb.getContainer(i, step)
			if err != nil {
				return err
			}
			containers = append(containers, container)
		}
		pod.Spec.Containers = containers
		//设置owner
		controller := true
		pod.OwnerReferences = append(pod.OwnerReferences,
			metav1.OwnerReference{
				APIVersion: pb.task.APIVersion,
				Kind:       pb.task.Kind,
				Name:       pb.task.Name,
				UID:        pb.task.UID,
				Controller: &controller,
			})
		return pb.Create(ctx, pod)
	}
	if err != nil {
		return err
	}
	return pb.UpdatePod(pod)
}

func (pb *PodBuilder) UpdatePod(pod *corev1.Pod) error {
	if pod.Status.Phase == corev1.PodRunning {
		if order, exist := pod.Annotations[TaskOrderAnnotation]; exist {
			idx, _ := strconv.Atoi(order)
			if len(pod.Spec.Containers) > idx {
				if idx >= 0 {
					if pod.Status.ContainerStatuses[idx].State.Terminated != nil {
						if pod.Status.ContainerStatuses[idx].State.Terminated.ExitCode != 0 {
							pod.Status.Phase = corev1.PodFailed
							return pb.Client.Status().Update(context.Background(), pod)
						}
						pod.Annotations[TaskOrderAnnotation] = strconv.Itoa(idx + 1)
						return pb.Client.Update(context.Background(), pod)
					}
				} else {
					if len(pod.Spec.Containers) > 0 {
						pod.Annotations[TaskOrderAnnotation] = "0"
						return pb.Client.Update(context.Background(), pod)
					}
				}
			}
		} else {
			if len(pod.Spec.Containers) > 0 {
				pod.Annotations[TaskOrderAnnotation] = "0"
				return pb.Client.Update(context.Background(), pod)
			}
		}

	}

	return nil
}

func (pb *PodBuilder) setPodVolumes(pod *corev1.Pod) {
	pod.Spec.Volumes = []corev1.Volume{
		corev1.Volume{
			Name: "entrypoint-volume",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		},
		corev1.Volume{
			Name: "podinfo",
			VolumeSource: corev1.VolumeSource{
				DownwardAPI: &corev1.DownwardAPIVolumeSource{
					Items: []corev1.DownwardAPIVolumeFile{
						corev1.DownwardAPIVolumeFile{
							Path: "order",
							FieldRef: &corev1.ObjectFieldSelector{
								FieldPath: "metadata.annotations['" + TaskOrderAnnotation + "']",
							},
						},
					},
				},
			},
		},
	}
}

func (pb *PodBuilder) setInitContainer(pod *corev1.Pod) {
	pod.Spec.InitContainers = []corev1.Container{
		corev1.Container{
			Name:            "init-container",
			Image:           "shenyisyn/entrypoint:v1.1",
			ImagePullPolicy: corev1.PullIfNotPresent,
			Command:         []string{"cp", "/app/entrypoint", "/entrypoint/bin"},
			VolumeMounts: []corev1.VolumeMount{
				corev1.VolumeMount{
					Name:      "entrypoint-volume",
					ReadOnly:  false,
					MountPath: "/entrypoint/bin",
				},
			},
		},
	}
}

func (pb *PodBuilder) getContainer(index int, taskStep taskv1alpha1.TaskStep) (corev1.Container, error) {
	container := corev1.Container{}
	container.Name = taskStep.Name
	container.Image = taskStep.Image
	container.ImagePullPolicy = corev1.PullIfNotPresent
	container.VolumeMounts = []corev1.VolumeMount{
		corev1.VolumeMount{
			Name:      "entrypoint-volume",
			ReadOnly:  false,
			MountPath: "/entrypoint/bin",
		},
		corev1.VolumeMount{
			Name:      "podinfo",
			ReadOnly:  false,
			MountPath: "/etc/podinfo",
		},
	}
	if len(taskStep.Command) == 0 {
		reference, err := name.ParseReference(taskStep.Image)
		if err != nil {
			return container, err
		}
		imageInterface, ok := ImageCache.Get(reference)
		if !ok {
			imageInterface, err = ParseImage(taskStep.Image)
			if err != nil {
				return container, err
			}
			ImageCache.Add(reference, imageInterface)
		}
		image := imageInterface.(*Image)

		taskStep.Command = image.Command["linux/amd64"].Command
		if len(taskStep.Command) == 0 && len(taskStep.Args) == 0 {
			taskStep.Args = image.Command["linux/amd64"].Args
		}
	}

	container.Command = []string{"/entrypoint/bin/entrypoint"}
	container.Args = []string{
		"--wait", "/etc/podinfo/order",
		"--waitcontent", strconv.Itoa(index),
		"--out", "stdout", // entrypoint 中 写上stdout 就会定向到标准输出
		"--command",
	}
	if len(taskStep.Command) > 0 {
		container.Args = append(container.Args, strings.Join(taskStep.Command, " "))
	}
	if len(taskStep.Args) > 0 {
		container.Args = append(container.Args, taskStep.Args...)
	}

	return container, nil
}

func (pb *PodBuilder) setPodMetadata(pod *corev1.Pod) {
	pod.Namespace = pb.task.Namespace
	pod.Name = PodPrefix + pb.task.Name // pod名称
	pod.Annotations = map[string]string{
		TaskOrderAnnotation: "0",
	}
}
