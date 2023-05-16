package webhook

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type PodAnnotations struct{}

// 返回的 error 如果不为空，本次请求会被决绝 admission.Denied(err.Error())
// 在 sigs.k8s.io/controller-runtime@v0.14.6/pkg/webhook/admission/defaulter_custom.go 中
// 如下函数被 admission.Handle 调用
func (p *PodAnnotations) Default(ctx context.Context, obj runtime.Object) error {
	log := logf.FromContext(ctx)

	// 参数 obj 为 api 请求的k8s资源对象，如想要获取原始请求 admission.Request 需要通过如下方式从 ctx 中获取
	// 在 admission.Handle 中会将原始请求放到 ctx 中，并提供了如下方法取出
	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		err := apierrors.NewInternalError(err)
		log.Error(err, "admission.Request not found in context")
		return err
	}

	// 如果是 tryRun 不做任何操作返回(https://kubernetes.io/zh-cn/docs/reference/access-authn-authz/extensible-admission-controllers/)
	if *req.DryRun {
		log.Info("dry run")
		return nil
	}

	pod, ok := obj.(*corev1.Pod)
	if !ok {
		// 外层函数会解析 apierrors.APIStatus 接口类型的错误，可以通过 apierrors 模块中的错误生成函数来初始化得到相应类型错误
		// 如果直接返回错误，则外层直接返回 admission.Denied(err.Error()) 错误
		err := apierrors.NewInternalError(fmt.Errorf("expected a pod but get %+v", obj))
		log.Error(err, "")
		return err
	}

	if pod.Annotations == nil {
		pod.Annotations = make(map[string]string)
	}

	if pod.Annotations["example-mutating-admission-webhook"] == "qwopt" {
		return nil
	}

	pod.Annotations["example-mutating-admission-webhook"] = "qwopt"
	log.Info("pod inject annotation", "pod", pod.GetName(), "example-mutating-admission-webhook", "qwopt")

	return nil
}