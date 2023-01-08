package controller

import (
	"context"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"try_operator/pkg/sysconfig"
)

const (
	ProxyControllerAnnotation = "jtthink"
	ingressAnnotationKey = "kubernetes.io/ingress.class"
)

type ProxyController struct {
	client.Client
}

func NewProxyController() *ProxyController {
	return &ProxyController{}
}

func (r *ProxyController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {

	obj := &networkingv1.Ingress{}
	err := r.Get(ctx, req.NamespacedName, obj)
	if err != nil {
		return reconcile.Result{}, err
	}
	// 取到需要的ingress对象
	if r.isNeedIngress(obj.Annotations) {
		klog.Info("get ingress!")
		err = sysconfig.AppConfig(obj)
		if err != nil {
			return reconcile.Result{}, nil
		}
	}

	return reconcile.Result{}, nil
}


func(r *ProxyController) InjectClient(c client.Client) error {
	r.Client = c
	return nil
}

func (r *ProxyController) isNeedIngress(annotations map[string]string) bool {
	if v, ok := annotations[ingressAnnotationKey]; ok && v== ProxyControllerAnnotation {
		return true
	}
	return false
}

func (r *ProxyController) IngressDeleteHandler(event event.DeleteEvent, limitingInterface workqueue.RateLimitingInterface) {
	// 取到需要的ingress
	if r.isNeedIngress(event.Object.GetAnnotations()) {
		klog.Info("delete ingress", event.Object.GetName())
		if err := sysconfig.DeleteConfig(event.Object.GetName(), event.Object.GetNamespace()); err != nil {
			klog.Error(err)
		}
	}
}

