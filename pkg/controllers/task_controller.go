package controllers

import (
	"context"
	"github.com/cleverhu/task-cicd/pkg/builder"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/event"

	taskv1alpha1 "github.com/cleverhu/task-cicd/pkg/apis/task/v1alpha1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type TaskController struct {
	client.Client
	//*clientset.Clientset
	E record.EventRecorder //记录事件
}

func NewTaskController(e record.EventRecorder) *TaskController {
	return &TaskController{E: e}
}

func (r *TaskController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	task := &taskv1alpha1.Task{}
	err := r.Get(ctx, req.NamespacedName, task)
	if err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}
	if !task.DeletionTimestamp.IsZero() {
		klog.V(2).Infof("task %s is been deleting.\n", task.Name)
		return reconcile.Result{}, nil
	}

	return reconcile.Result{}, r.UpdatePodIfNeeded(ctx, r.Client, task)
}

func (r *TaskController) InjectClient(c client.Client) error {
	r.Client = c
	return nil
}

func (r *TaskController) UpdatePodIfNeeded(ctx context.Context, c client.Client, task *taskv1alpha1.Task) error {
	return builder.NewPodBuilder(c, task).Build(ctx)
}

func (r *TaskController) UpdatePodFunc(event event.UpdateEvent, limitingInterface workqueue.RateLimitingInterface) {
	for _, ref := range event.ObjectNew.GetOwnerReferences() {
		if ref.Controller != nil && *ref.Controller && ref.Kind == taskv1alpha1.TaskKind && ref.APIVersion == taskv1alpha1.TaskApiVersion {
			limitingInterface.Add(reconcile.Request{NamespacedName: types.NamespacedName{
				Namespace: event.ObjectNew.GetNamespace(),
				Name:      ref.Name,
			}})
		}
	}
}
