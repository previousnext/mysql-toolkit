package env

const (
	DockerUsername = "DOCKER_USERNAME"
	DockerPassword = "DOCKER_PASSWORD"
	DockerImage    = "DOCKER_IMAGE"

	MySQLFile     = "MYSQL_FILE"
	MySQLHostname = "MYSQL_HOSTNAME"
	MySQLUsername = "MYSQL_USERNAME"
	MySQLPassword = "MYSQL_PASSWORD"
	MySQLDatabase = "MYSQL_DATABASE"
	MySQLProtocol = "MYSQL_PROTOCOL"
	MySQLPort     = "MYSQL_PORT"
	MySQLMaxConn  = "MYSQL_MAX_CONN"
	MySQLConfig   = "MYSQL_CONFIG"

	AWSRegion          = "AWS_REGION"
	AWSAccessKeyID     = "AWS_ACCESS_KEY_ID"
	AWSSecretAccessKey = "AWS_SECRET_ACCESS_KEY"

	AWSIAMRole = "AWS_IAM_ROLE"

	AWSS3Bucket = "AWS_S3_BUCKET"

	AWSCodeBuildDockerfile = "AWS_CODEBUILD_DOCKERFILE"
	AWSCodeBuildSpec       = "AWS_CODEBUILD_SPEC"
	AWSCodeBuildProject    = "AWS_CODEBUILD_PROJECT"
	AWSCodeBuildCompute    = "AWS_CODEBUILD_COMPUTE"
	AWSCodeBuildImage      = "AWS_CODEBUILD_IMAGE"

	K8sNamespace            = "K8S_NAMESPACE"
	K8sCronJobFrequency     = "K8S_CRONJOB_FREQUENCY"
	K8sCronJobImage         = "K8S_CRONJOB_IMAGE"
	K8sCronJobCPU           = "K8S_CRONJOB_CPU"
	K8sCronJobMemory        = "K8S_CRONJOB_MEMORY"
	K8sConfigMapKeyHostname = "K8S_CONFIGMAP_KEY_HOSTNAME"
	K8sConfigMapKeyUsername = "K8S_CONFIGMAP_KEY_USERNAME"
	K8sConfigMapKeyPassword = "K8S_CONFIGMAP_KEY_PASSWORD"
	K8sConfigMapKeyDatabase = "K8S_CONFIGMAP_KEY_DATABASE"
	K8sConfigMapKeyImage    = "K8S_CONFIGMAP_KEY_IMAGE"
)
