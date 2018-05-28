MySQL Toolkit
=============

[![CircleCI](https://circleci.com/gh/previousnext/mysql-toolkit.svg?style=svg)](https://circleci.com/gh/previousnext/mysql-toolkit)

**Maintainer**: Nick Schuch

Toolkit for working with MySQL databases.

* Exporting
* Building a container image for developers (Docker Compose etc)
* Integrating with container orchestrators

## Usage

The below example will cover:

* Dumping a sanitized database
* Building a Docker image with AWS CodeBuild

### Dumping a MySQL database

The `db dump` command will:

* Connect to the MySQL database
* Dump the database
* Sanitize based on the `mysql-config` file

```bash
mtk db dump --hostname=127.0.0.1 \
            --username=root \
            --password=password \
            --database=mydb \
            --config=example/config.yml \
            --file=/tmp/db.sql
```

You can now store this file as you see fit, allowing developers to have a sanitized database for local development.

### Building a MySQL container image

The below command will:

* Upload a ZIP containing a Dockerfile, buildspec.yml and db.sql to S3
* Create a CodeBuild project
* Trigger a CodeBuild project build

For this command you will need to create an IAM role which CodeBuild can use to access to S3 bucket.

https://docs.aws.amazon.com/codebuild/latest/userguide/setting-up.html#setting-up-service-role

```bash
mtk build aws --project=mysql-toolkit-example \
              --dockerfile=example/Dockerfile \
              --spec=example/buildspec.yml \
              --bucket=mysql-sanitized \
              --role=arn:aws:iam::XXXXXXXXXXXXX:role/mysql-toolkit \
              --docker-username=dockeruser \
              --docker-password=password \
              --docker-image=example/database:latest \
              --file=/tmp/db.sql
```

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
