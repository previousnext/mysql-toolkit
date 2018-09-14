MySQL Toolkit
=============

[![CircleCI](https://circleci.com/gh/previousnext/mysql-toolkit.svg?style=svg)](https://circleci.com/gh/previousnext/mysql-toolkit)

**Maintainer**: Nick Schuch

A toolkit for santizing, packaging and deploying MySQL images.

https://www.previousnext.com.au/blog/sharing-state-containerized-databases-developers

![Diagram](/docs/diagram.png "Diagram")

## Commands

Dump a MySQL database to a sanitized SQL file.

```
mtk db dump --hostname=HOSTNAME --username=USERNAME --password=PASSWORD --database=DATABASE --file=FILE [<flags>]
  Dump the database using a MySQL connection

  --hostname=HOSTNAME    Hostname for connecting to Mysql
  --username=USERNAME    Username for connecting to Mysql
  --password=PASSWORD    Password for connecting to Mysql
  --database=DATABASE    Database for connecting to Mysql
  --protocol="tcp"       Protocol for connecting to Mysql
  --port="3306"          Port for connecting to Mysql
  --max-conn=50          Maximum amount of open connections
  --config="config.yml"  Policy for dumping the database
  --file=FILE            Location to save the dumped database
```

Run the kubernetes operator.

```
mtk db operator --role=ROLE --bucket=BUCKET --key-id=KEY-ID --access-key=ACCESS-KEY --docker-username=DOCKER-USERNAME --docker-password=DOCKER-PASSWORD [<flags>]
  Continuous image builder using a MySQL connection as the source

  --namespace=""           Namespace to lookup ConfigMaps
  --frequency="@daily"     How ofter CronJobs should create new CodeBuild
                           project builds
  --image="previousnext/mysql-toolkit:latest"  
                           Image to use for running the CronJob
  --cpu="250m"             How much CPU resource should be assigned to the
                           CronJob
  --memory="256Mi"         How much memory resource should be assigned to the
                           CronJob
  --key-hostname="mysql.hostname"  
                           ConfigMap key which containers the MySQL hostname
  --key-username="mysql.username"  
                           ConfigMap key which containers the MySQL username
  --key-password="mysql.password"  
                           ConfigMap key which containers the MySQL password
  --key-database="mysql.database"  
                           ConfigMap key which containers the MySQL database
  --key-image="mysql.docker.image"  
                           ConfigMap key which containers the MySQL image
  --role=ROLE              ServiceRole or IAM resource which grants access to
                           the S3 bucket
  --bucket=BUCKET          Bucket to upload the file temporarily before
                           CodeBuild runs
  --key-id=KEY-ID          AWS Credentials
  --access-key=ACCESS-KEY  AWS Credentials
  --docker-username=DOCKER-USERNAME  
                           Username for the Docker Registry
  --docker-password=DOCKER-PASSWORD  
                           Password for the Docker Registry
```

Build a pre-loaded database container image with AWS CodeBuild.

```
mtk build aws --project=PROJECT --dockerfile=DOCKERFILE --spec=SPEC --bucket=BUCKET --role=ROLE --docker-username=DOCKER-USERNAME --docker-password=DOCKER-PASSWORD --docker-image=DOCKER-IMAGE --file=FILE [<flags>]
  Build an image using AWS CodeBuild

  --region="ap-southeast-2"    Region to run the build
  --project=PROJECT            Name for the CodeBuild project
  --compute="BUILD_GENERAL1_SMALL"  
                               Size of the compute for the build
  --image="aws/codebuild/docker:17.09.0"  
                               CodeBuild image to use for executing the build
  --dockerfile=DOCKERFILE      Path to the Dockerfile use to build the image
  --spec=SPEC                  Path to the BuildSpec use to build the image
  --bucket=BUCKET              Bucket to upload the file temporarily before
                               CodeBuild runs
  --role=ROLE                  ServiceRole or IAM resource which grants access
                               to the S3 bucket
  --docker-username=DOCKER-USERNAME  
                               Username for the Docker Registry
  --docker-password=DOCKER-PASSWORD  
                               Password for the Docker Registry
  --docker-image=DOCKER-IMAGE  Name to push to the registry
  --file=FILE                  Path to the Mysql database use to build the
                               image
```

## Installation

Installation instructions are availabe via the [releases page](https://github.com/previousnext/mysql-toolkit/releases)

## Development

### Getting started

For steps on getting started with Go:

https://golang.org/doc/install

To get a checkout of the project run the following commands:

```bash
# Make sure the parent directories exist.
mkdir -p $GOPATH/src/github.com/previousnext

# Checkout the codebase.
git clone git@github.com:previousnext/mysql-toolkit $GOPATH/src/github.com/previousnext/mysql-toolkit

# Change into the project to run workflow commands.
cd $GOPATH/src/github.com/previousnext/mysql-toolkit
```

### Documentation

See `/docs`

### Resources

* [Dave Cheney - Reproducible Builds](https://www.youtube.com/watch?v=c3dW80eO88I)
* [Bryan Cantril - Debugging under fire](https://www.youtube.com/watch?v=30jNsCVLpAE&t=2675s)
* [Sam Boyer - The New Era of Go Package Management](https://www.youtube.com/watch?v=5LtMb090AZI)
* [Kelsey Hightower - From development to production](https://www.youtube.com/watch?v=XL9CQobFB8I&t=787s)

### Tools

```bash
# Dependency management
go get -u github.com/golang/dep/cmd/dep

# Testing
go get -u github.com/golang/lint/golint

# Release management.
go get -u github.com/tcnksm/ghr

# Build
go get -u github.com/mitchellh/gox
```

### Workflow

**Testing**

```bash
make lint
```

**Building**

```bash
make build
```

**Releasing**

```bash
make release
```
