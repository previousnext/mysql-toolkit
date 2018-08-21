package tests

type Backend struct {
	Host   string `k8s_configmap:"backend.host"`
	Port   int64  `k8s_configmap:"backend.port"`
	Secure bool   `k8s_configmap:"backend.secure"`
	Token  bool   `k8s_secret:"backend.token"`
}
