package main

import (
	"github.com/cleverhu/task-cicd/pkg/builder"
	"github.com/cleverhu/task-cicd/pkg/k8sconfig"
)

func main() {
	builder.InitImageCache(100)
	k8sconfig.InitManager()
}
