AWS CodeBuild
=============

## Overview

Offload your image build pipeline to AWS.

https://aws.amazon.com/codebuild

## Requirements

* S3 Bucket
* [IAM Role](https://docs.aws.amazon.com/codebuild/latest/userguide/setting-up.html#setting-up-service-role)

## Usage

* Upload a `.zip` to S3 with:
 * Dockerfile
 * buildspec.yml 
 * db.sql
* Create a CodeBuild project
* Trigger a CodeBuild project build

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