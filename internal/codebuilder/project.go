package codebuilder

import (
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/codebuild"

	"github.com/previousnext/mysql-toolkit/internal/envar"
)

// CreateProjectInput used for CreateProject CodeBuild method.
func (b *BuildParams) CreateProjectInput(source string) *codebuild.CreateProjectInput {
	return &codebuild.CreateProjectInput{
		Name: aws.String(b.Project),
		Source: &codebuild.ProjectSource{
			Type:     aws.String(codebuild.ArtifactsTypeS3),
			Location: aws.String(filepath.Join(b.Bucket, source)),
		},
		Artifacts: &codebuild.ProjectArtifacts{
			Type: aws.String(codebuild.ArtifactsTypeNoArtifacts),
		},
		Environment: &codebuild.ProjectEnvironment{
			Type:        aws.String(codebuild.EnvironmentTypeLinuxContainer),
			ComputeType: aws.String(b.Compute),
			Image:       aws.String(b.Image),
			EnvironmentVariables: []*codebuild.EnvironmentVariable{
				{
					Name:  aws.String(envar.DockerUsername),
					Value: aws.String(b.Docker.Username),
				},
				{
					Name:  aws.String(envar.DockerPassword),
					Value: aws.String(b.Docker.Password),
				},
				{
					Name:  aws.String(envar.DockerImage),
					Value: aws.String(b.Docker.Image),
				},
			},
		},
		ServiceRole: aws.String(b.Role),
	}
}

// UpdateProjectInput used for UpdateProject CodeBuild method.
func (b *BuildParams) UpdateProjectInput(source string) *codebuild.UpdateProjectInput {
	return &codebuild.UpdateProjectInput{
		Name: aws.String(b.Project),
		Source: &codebuild.ProjectSource{
			Type:     aws.String(codebuild.ArtifactsTypeS3),
			Location: aws.String(filepath.Join(b.Bucket, source)),
		},
		Artifacts: &codebuild.ProjectArtifacts{
			Type: aws.String(codebuild.ArtifactsTypeNoArtifacts),
		},
		Environment: &codebuild.ProjectEnvironment{
			Type:        aws.String(codebuild.EnvironmentTypeLinuxContainer),
			ComputeType: aws.String(b.Compute),
			Image:       aws.String(b.Image),
			EnvironmentVariables: []*codebuild.EnvironmentVariable{
				{
					Name:  aws.String(envar.DockerUsername),
					Value: aws.String(b.Docker.Username),
				},
				{
					Name:  aws.String(envar.DockerPassword),
					Value: aws.String(b.Docker.Password),
				},
				{
					Name:  aws.String(envar.DockerImage),
					Value: aws.String(b.Docker.Image),
				},
			},
		},
		ServiceRole: aws.String(b.Role),
	}
}
