package k8sconfig

import (
	corev1 "k8s.io/api/core/v1"
	"log"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	taskv1alpha1 "github.com/cleverhu/task-cicd/pkg/apis/task/v1alpha1"
	"github.com/cleverhu/task-cicd/pkg/controllers"

	"sigs.k8s.io/controller-runtime/pkg/builder"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

//初始化 控制器管理器
func InitManager() {
	//taskClient := versioned.NewForConfigOrDie(K8sRestConfig())
	logf.SetLogger(zap.New())
	mgr, err := manager.New(K8sRestConfig(),
		manager.Options{
			Logger: logf.Log.WithName("mycicd"),
		})

	if err != nil {
		log.Fatal("创建管理器失败:", err.Error())
	}
	//Schema定义了资源序列化和反序列化的方法以及资源类型和版本的对应关系
	err = taskv1alpha1.SchemeBuilder.AddToScheme(mgr.GetScheme())
	if err != nil {
		mgr.GetLogger().Error(err, "unable add schema")
		os.Exit(1)
	}
	//初始化控制器对象
	cicdController := controllers.NewTaskController(
		mgr.GetEventRecorderFor("mycicd"),
	)

	if err = builder.ControllerManagedBy(mgr).
		For(&taskv1alpha1.Task{}).
		Watches(&source.Kind{Type: &corev1.Pod{}}, handler.Funcs{UpdateFunc: cicdController.UpdatePodFunc}).
		Complete(cicdController); err != nil {
		mgr.GetLogger().Error(err, "unable to create manager")
		os.Exit(1)
	}

	if err = mgr.Start(signals.SetupSignalHandler()); err != nil {
		mgr.GetLogger().Error(err, "unable to start manager")
	}
}
