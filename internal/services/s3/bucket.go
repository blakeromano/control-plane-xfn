package s3

import (
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/response"
	"github.com/pkg/errors"
)

func RenderS3BucketResources(req *fnv1beta1.RunFunctionRequest, rsp *fnv1beta1.RunFunctionResponse, oxr *resource.Composite, log logging.Logger) (*fnv1beta1.RunFunctionRequest, *fnv1beta1.RunFunctionResponse) {
	desired, err := request.GetDesiredComposedResources(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get desired composed resources in %T", rsp))
		return req, rsp
	}

	observed, err := request.GetObservedComposedResources(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get observed composed resources in %T", rsp))
		return req, rsp
	}

	dataClassification, _ := oxr.Resource.GetString("spec.resourceConfig.dataClassification")

	desired["bucket"] = resource.NewDesiredComposed()

	RenderS3BucketResource(desired["bucket"], oxr, dataClassification)

	desired["ses-encryption"] = resource.NewDesiredComposed()

	RenderS3SesConfig(desired["ses-encryption"], oxr)

	desired["bucket-lifecycle-rule"] = resource.NewDesiredComposed()

	GenerateS3BucketLifecyclePolicy(desired["bucket-lifecycle-rule"], oxr, dataClassification)

	if _, exists := observed["bucket"]; exists {
		log.Debug("Bucket exists already", "bucket", oxr.Resource.GetClaimReference().Name)
		bucketArn, _ := observed["bucket"].Resource.GetString("status.atProvider.arn")
		bucketName, _ := observed["bucket"].Resource.GetString("status.atProvider.id")
		oxr.Resource.SetString("status.bucketName", bucketName)
		oxr.Resource.SetString("status.bucketArn", bucketArn)
	}

	if err := response.SetDesiredComposedResources(rsp, desired); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot set desired composed resources in %T", rsp))
		return req, rsp
	}

	oxr.Resource.SetValue("metadata.managedFields", nil)

	if err := response.SetDesiredCompositeResource(rsp, oxr); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot set desired composite resource in %T", rsp))
		return req, rsp
	}

	return req, rsp
}

func RenderS3BucketResource(desired *resource.DesiredComposed, oxr *resource.Composite, dataClassification string) {
	desired.Resource.Object = map[string]interface{}{
		"apiVersion": "s3.aws.upbound.io/v1beta1",
		"kind":       "Bucket",
		"metadata": map[string]interface{}{
			"name":          oxr.Resource.GetClaimReference().Name,
		},
		"spec": map[string]interface{}{
			"forProvider": map[string]interface{}{
				"region": "us-east-2",
				"tags":   map[string]interface{}{
					"name": oxr.Resource.GetClaimReference().Name,
					"namespace": oxr.Resource.GetClaimReference().Namespace,
					"dataClassification": dataClassification,
				},
			},
			"providerConfigRef": map[string]interface{}{
				"name": "provider-aws",
			},
		},
	}
}

func RenderS3SesConfig(desired *resource.DesiredComposed, oxr *resource.Composite) {
	desired.Resource.Object = map[string]interface{}{
		"apiVersion": "s3.aws.upbound.io/v1beta1",
		"kind":       "BucketServerSideEncryptionConfiguration",
		"metadata": map[string]interface{}{
			"name":          oxr.Resource.GetClaimReference().Name,
		},
		"spec": map[string]interface{}{
			"forProvider": map[string]interface{}{
				"region": "us-east-2",
				"bucketSelector": map[string]interface{}{
					"matchControllerRef": true,
				},
				"rule": []interface{}{
					map[string]interface{}{
						"applyServerSideEncryptionByDefault": []interface{}{
							map[string]interface{}{
								"sseAlgorithm": "AES256",
							},
						},
					},
				},
			},
			"providerConfigRef": map[string]interface{}{
				"name": "provider-aws",
			},
		},
	}
}

func GenerateS3BucketLifecyclePolicy(desired *resource.DesiredComposed, oxr *resource.Composite, dataClassification string) {
	desired.Resource.Object = map[string]interface{}{
		"apiVersion": "s3.aws.upbound.io/v1beta1",
		"kind":       "BucketLifecycleConfiguration",
		"metadata": map[string]interface{}{
			"name":          oxr.Resource.GetClaimReference().Name,
		},
		"spec": map[string]interface{}{
			"forProvider": map[string]interface{}{
				"region": "us-east-2",
				"bucketSelector": map[string]interface{}{
					"matchControllerRef": true,
				},
			},
			"providerConfigRef": map[string]interface{}{
				"name": "provider-aws",
			},
		},
	}

	switch dataClassification {
		case "public":
			desired.Resource.SetValue("spec.forProvider.rule", publicRule)
		case "internal":
			desired.Resource.SetValue("spec.forProvider.rule", internalRule)
		case "confidential":
			desired.Resource.SetValue("spec.forProvider.rule", confidentialRule)
	}
}