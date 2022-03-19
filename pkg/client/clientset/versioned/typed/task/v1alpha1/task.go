// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/cleverhu/task-cicd/pkg/apis/task/v1alpha1"
	scheme "github.com/cleverhu/task-cicd/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// TasksGetter has a method to return a TaskInterface.
// A group's client should implement this interface.
type TasksGetter interface {
	Tasks(namespace string) TaskInterface
}

// TaskInterface has methods to work with Task resources.
type TaskInterface interface {
	Create(ctx context.Context, task *v1alpha1.Task, opts v1.CreateOptions) (*v1alpha1.Task, error)
	Update(ctx context.Context, task *v1alpha1.Task, opts v1.UpdateOptions) (*v1alpha1.Task, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Task, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.TaskList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Task, err error)
	TaskExpansion
}

// tasks implements TaskInterface
type tasks struct {
	client rest.Interface
	ns     string
}

// newTasks returns a Tasks
func newTasks(c *CicdV1alpha1Client, namespace string) *tasks {
	return &tasks{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the task, and returns the corresponding task object, and an error if there is any.
func (c *tasks) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Task, err error) {
	result = &v1alpha1.Task{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("tasks").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Tasks that match those selectors.
func (c *tasks) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.TaskList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.TaskList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("tasks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested tasks.
func (c *tasks) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("tasks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a task and creates it.  Returns the server's representation of the task, and an error, if there is any.
func (c *tasks) Create(ctx context.Context, task *v1alpha1.Task, opts v1.CreateOptions) (result *v1alpha1.Task, err error) {
	result = &v1alpha1.Task{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("tasks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(task).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a task and updates it. Returns the server's representation of the task, and an error, if there is any.
func (c *tasks) Update(ctx context.Context, task *v1alpha1.Task, opts v1.UpdateOptions) (result *v1alpha1.Task, err error) {
	result = &v1alpha1.Task{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("tasks").
		Name(task.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(task).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the task and deletes it. Returns an error if one occurs.
func (c *tasks) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("tasks").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *tasks) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("tasks").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched task.
func (c *tasks) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Task, err error) {
	result = &v1alpha1.Task{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("tasks").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
