package k8s

type Keys struct {
	Hostname string
	Username string
	Password string
	Database string
	Image    string
}

type Resources struct {
	CPU    string
	Memory string
}

type Docker struct {
	Username string
	Password string
}

type AWSCredentials struct {
	KeyID     string
	AccessKey string
}
