package codebuilder

import "github.com/aws/aws-sdk-go/service/codebuild"

const (
	DefaultBuildImage  = "aws/codebuild/docker:17.09.0"
	DefaultComputeSize = codebuild.ComputeTypeBuildGeneral1Small
)

func (b *BuildParams) Defaults() {
	if b.Compute == "" {
		b.Compute = DefaultComputeSize
	}

	if b.Image == "" {
		b.Image = DefaultBuildImage
	}
}
