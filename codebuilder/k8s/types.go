package k8s

// Keys for ConfigMap discovery.
type Keys struct {
	Hostname string
	Username string
	Password string
	Database string
	Image    string
}

// Resources for dumping and pushing CodeBuild context.
type Resources struct {
	CPU    string
	Memory string
}

// Docker credentials.
type Docker struct {
	Username string
	Password string
}

// AWSCredentials for connecting to S3 and CodeBuild.
type AWSCredentials struct {
	KeyID     string
	AccessKey string
}
