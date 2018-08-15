package tests

type Authentication struct {
	Username string `k8s_secret:"auth.username"`
	Password string `k8s_secret:"auth.password"`
}
