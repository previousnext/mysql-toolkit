package secrets

// Secrets used for Acquia to AWS CodeBuild to Docker Hub.
type Secrets struct {
	AcquiaUsername string `k8s_secret:"acquia.username"`
	AcquiaPassword string `k8s_secret:"acquia.password"`
	DockerUsername string `k8s_secret:"docker.username"`
	DockerPassword string `k8s_secret:"docker.password"`
	AWSRole        string `k8s_secret:"aws.role"`
	AWSKey         string `k8s_secret:"aws.key.id"`
	AWSAccess      string `k8s_secret:"aws.key.access"`
	AWSBucket      string `k8s_secret:"aws.bucket"`
}
