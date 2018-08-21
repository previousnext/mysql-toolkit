package envar

const (
	// AcquiaUsername for the Acquia platform.
	AcquiaUsername = "ACQUIA_USERNAME"
	// AcquiaPassword for the Acquia platform.
	AcquiaPassword = "ACQUIA_PASSWORD"
	// AcquiaSite hosted on the Acquia platform.
	AcquiaSite = "ACQUIA_SITE"
	// AcquiaEnvironment hosted on the Acquia platform.
	AcquiaEnvironment = "ACQUIA_ENVIRONMENT"
	// AcquiaDatabase hosted on the Acquia platform.
	AcquiaDatabase = "ACQUIA_DATABASE"

	// DockerUsername for Docker Hub.
	DockerUsername = "DOCKER_USERNAME"
	// DockerPassword for Docker Hub.
	DockerPassword = "DOCKER_PASSWORD"
	// DockerImage stored on Docker Hub.
	DockerImage = "DOCKER_IMAGE"
	// Dockerfile used to build an image.
	Dockerfile = "DOCKER_FILE"

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

	// OperatorNamespaceWatch these namespaces for CustomResources.
	OperatorNamespaceWatch = "OPERATOR_NAMESPACE_WATCH"
	// OperatorNamespace used execute Jobs.
	OperatorNamespace = "OPERATOR_NAMESPACE"
	// OperatorSecret used to loading configuration.
	OperatorSecret = "OPERATOR_SECRET"
	// OperatorJobImage used to loading configuration.
	OperatorJobImage = "OPERATOR_JOB_IMAGE"
	// OperatorJobCPU used to loading configuration.
	OperatorJobCPU = "OPERATOR_JOB_CPU"
	// OperatorJobMemory used to loading configuration.
	OperatorJobMemory = "OPERATOR_JOB_MEMORY"
	// OperatorResync tells the operator how long before it resyncs CRDs.
	OperatorResync = "OPERATOR_RESYNC"
)
