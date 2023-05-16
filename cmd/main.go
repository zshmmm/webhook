package main

import (
	"flag"
	"os"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"podwebhook/webhook"
)

var (
	webhookLogger = log.Log.WithName("pod-admission-webhook")
)

func init() {
	log.SetLogger(zap.New())
}

func main() {
	var (
		certDir string
		port    int
		enableLeaderElection bool
	)

	// 初始化参数
	flag.IntVar(&port, "port", 8443, "Webhook server port.")
	flag.StringVar(&certDir, "certDir", "/etc/webhook/certs", "certDir is the directory that contains the server key and certificate.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
	"Enable leader election for controller manager. "+
		"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	// 初始化 manager 实例
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		// webhook 服务端口
		Port:   port,
		// webhook 服务端证书目录，使用 controller-runtime 证书文件必须指定为：tls.key 和 tls.crt
		// webhook deployment 中使用的 secret 生成时必须满足当前需求
		CertDir: certDir,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "f96a6927.qwwebhook.io",
	})

	if err != nil {
		webhookLogger.Error(err, "create manager failed")
		os.Exit(1)
	}

	// 通过 manager 创建 webhook  server，并将自定义的处理逻辑绑定到 webhook 的 Handler 中
	if err := ctrl.NewWebhookManagedBy(mgr).
		For(&corev1.Pod{}).
		WithDefaulter(&webhook.PodAnnotations{}).
		RecoverPanic().
		Complete(); err != nil {
		webhookLogger.Error(err, "create webhook failed")
		os.Exit(1)
	}

	webhookLogger.Info("strating manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		webhookLogger.Error(err, "start manager failed")
		os.Exit(1)
	}
}