package env

const (
	// DockerUsername for Docker Hub.
	DockerUsername = "DOCKER_USERNAME"
	// DockerPassword for Docker Hub.
	DockerPassword = "DOCKER_PASSWORD"
	// DockerImage stored on Docker Hub.
	DockerImage = "DOCKER_IMAGE"

	// MySQLFile for dumping and loading.
	MySQLFile = "MYSQL_FILE"
	// MySQLHostname for connecting to MySQL database.
	MySQLHostname = "MYSQL_HOSTNAME"
	// MySQLUsername for connecting to MySQL database.
	MySQLUsername = "MYSQL_USERNAME"
	// MySQLPassword for connecting to MySQL database.
	MySQLPassword = "MYSQL_PASSWORD"
	// MySQLDatabase for connecting to MySQL database.
	MySQLDatabase = "MYSQL_DATABASE"
	// MySQLProtocol for connecting to MySQL database.
	MySQLProtocol = "MYSQL_PROTOCOL"
	// MySQLPort for connecting to MySQL database.
	MySQLPort = "MYSQL_PORT"
	// MySQLMaxConn for connecting to MySQL database.
	MySQLMaxConn = "MYSQL_MAX_CONN"
	// MySQLConfig for connecting to MySQL database.
	MySQLConfig = "MYSQL_CONFIG"

	// AWSRegion where CodeBuild will be executed.
	AWSRegion = "AWS_REGION"
	// AWSAccessKeyID for logging into AWS.
	AWSAccessKeyID = "AWS_ACCESS_KEY_ID"
	// AWSSecretAccessKey for logging into AWS.
	AWSSecretAccessKey = "AWS_SECRET_ACCESS_KEY"

	// AWSIAMRole for CodeBuild to load build context.
	AWSIAMRole = "AWS_IAM_ROLE"

	// AWSS3Bucket for storing build context.
	AWSS3Bucket = "AWS_S3_BUCKET"

	// AWSCodeBuildDockerfile packaged in build context.
	AWSCodeBuildDockerfile = "AWS_CODEBUILD_DOCKERFILE"
	// AWSCodeBuildSpec packaged in build context.
	AWSCodeBuildSpec = "AWS_CODEBUILD_SPEC"
	// AWSCodeBuildProject name.
	AWSCodeBuildProject = "AWS_CODEBUILD_PROJECT"
	// AWSCodeBuildCompute used to package image.
	AWSCodeBuildCompute = "AWS_CODEBUILD_COMPUTE"
	// AWSCodeBuildImage to use for building the image.
	AWSCodeBuildImage = "AWS_CODEBUILD_IMAGE"

	// K8sNamespace to ConfigMap discovery.
	K8sNamespace = "K8S_NAMESPACE"
	// K8sCronJobFrequency for running created/updates CronJobs.
	K8sCronJobFrequency = "K8S_CRONJOB_FREQUENCY"
	// K8sCronJobImage for running created/updates CronJobs.
	K8sCronJobImage = "K8S_CRONJOB_IMAGE"
	// K8sCronJobCPU for running created/updates CronJobs.
	K8sCronJobCPU = "K8S_CRONJOB_CPU"
	// K8sCronJobMemory for running created/updates CronJobs.
	K8sCronJobMemory = "K8S_CRONJOB_MEMORY"
	// K8sConfigMapKeyHostname for ConfigMap discovery.
	K8sConfigMapKeyHostname = "K8S_CONFIGMAP_KEY_HOSTNAME"
	// K8sConfigMapKeyUsername for ConfigMap discovery.
	K8sConfigMapKeyUsername = "K8S_CONFIGMAP_KEY_USERNAME"
	// K8sConfigMapKeyPassword for ConfigMap discovery.
	K8sConfigMapKeyPassword = "K8S_CONFIGMAP_KEY_PASSWORD"
	// K8sConfigMapKeyDatabase for ConfigMap discovery.
	K8sConfigMapKeyDatabase = "K8S_CONFIGMAP_KEY_DATABASE"
	// K8sConfigMapKeyImage for ConfigMap discovery.
	K8sConfigMapKeyImage = "K8S_CONFIGMAP_KEY_IMAGE"
)
