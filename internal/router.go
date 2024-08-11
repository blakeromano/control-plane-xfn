package internal

import (
	"errors"

	"github.com/crossplane/control-plane-function/internal/services/s3"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/response"
)

func Router(req *fnv1beta1.RunFunctionRequest, rsp *fnv1beta1.RunFunctionResponse, oxr *resource.Composite, log logging.Logger) (*fnv1beta1.RunFunctionRequest, *fnv1beta1.RunFunctionResponse) {
	switch oxr.Resource.GetKind() {
	case "XS3Bucket":
		log.Debug("rendering s3 bucket resource")	
		s3.RenderS3BucketResources(req, rsp, oxr, log)
	default:
		response.Warning(rsp, errors.New("object kind not recognized"))
		return req, rsp
	}

	configureManagementPolicies(req, rsp, oxr, log)
	readyCheck(req, rsp, oxr, log)

	return req, rsp
}
