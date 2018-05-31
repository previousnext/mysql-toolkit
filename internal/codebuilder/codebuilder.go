package codebuilder

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/previousnext/mysql-toolkit/internal/awsutils"
)

// BuildParams for the codebuilder.
type BuildParams struct {
	Project    string
	Region     string
	Compute    string
	Image      string
	Role       string
	Bucket     string
	Dockerfile string
	BuildSpec  string
	Database   string
	Docker     Docker
}

// Build the Docker image with AWS Codebuild.
// Inspired by: https://docs.aws.amazon.com/codebuild/latest/userguide/sample-docker.html
func Build(w io.Writer, params BuildParams) error {
	fmt.Fprintln(w, "Validating configuration")

	// Validate that the params.
	err := params.Validate()
	if err != nil {
		return errors.Wrap(err, "validation failed")
	}

	// Get a session which we can share with the packaging and building components.
	session, err := awsutils.NewSession(params.Region)
	if err != nil {
		return errors.Wrap(err, "session failed")
	}

	fmt.Fprintln(w, "Compiling source (Dockerfile/BuildSpec/Database)")

	pkg, err := compile(params.Project, params.Dockerfile, params.BuildSpec, params.Database)
	if err != nil {
		return errors.Wrap(err, "failed to package source")
	}

	var (
		client   = codebuild.New(session)
		uploader = s3manager.NewUploader(session)
	)

	fmt.Fprintln(w, "Uploading source")

	err = uploadFile(uploader, params.Bucket, pkg)
	if err != nil {
		return errors.Wrap(err, "failed to upload source")
	}

	fmt.Fprintln(w, "Deleting local copy of source")

	err = os.Remove(pkg)
	if err != nil {
		return errors.Wrap(err, "failed to remove source")
	}

	fmt.Fprintln(w, "Creating Codebuild project")

	err = createProject(client, params, filepath.Base(pkg))
	if err != nil {
		return errors.Wrap(err, "failed to create project")
	}

	fmt.Fprintln(w, "Starting Codebuild")

	resp, err := client.StartBuild(&codebuild.StartBuildInput{
		ProjectName: aws.String(params.Project),
	})
	if err != nil {
		return errors.Wrap(err, "failed to start build")
	}

	fmt.Fprintln(w, "Waiting for build to finish")

	err = waitForBuild(client, resp.Build.Id)
	if err != nil {
		return errors.Wrap(err, "build failed")
	}

	fmt.Fprintln(w, "Complete!")

	return nil
}

// Helper function to upload a file to an S3 bucket.
func uploadFile(uploader *s3manager.Uploader, bucket, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return errors.Wrap(err, "failed to open file")
	}
	defer f.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filepath.Base(file)),
		Body:   f,
	})
	if err != nil {
		return errors.Wrap(err, "failed to upload file")
	}

	return nil
}

// Helper function to create/update a CodeBuild project.
func createProject(client *codebuild.CodeBuild, params BuildParams, pkg string) error {
	_, err := client.CreateProject(params.CreateProjectInput(pkg))
	if err != nil {
		if awserr, ok := err.(awserr.Error); ok {
			if awserr.Code() == codebuild.ErrCodeResourceAlreadyExistsException {
				_, err := client.UpdateProject(params.UpdateProjectInput(pkg))
				if err != nil {
					return errors.Wrap(err, "failed to update")
				}
			} else {
				return errors.Wrap(err, "failed to create")
			}
		}
	}

	return nil
}

// Helper function to wait for the CodeBuild to finish.
func waitForBuild(client *codebuild.CodeBuild, id *string) error {
	limiter := time.Tick(time.Second * 10)

	for {
		<-limiter

		resp, err := client.BatchGetBuilds(&codebuild.BatchGetBuildsInput{
			Ids: []*string{id},
		})
		if err != nil {
			return errors.Wrap(err, "failed to get build")
		}

		if len(resp.Builds) == 0 {
			return errors.New("cannot find build")
		}

		build := resp.Builds[0]

		if *build.BuildComplete {
			if *build.BuildStatus == codebuild.StatusTypeSucceeded {
				return nil
			}

			return fmt.Errorf("build finished with status: %s", *build.BuildStatus)
		}
	}

	return errors.New("unable to determine build status")
}
