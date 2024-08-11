package internal

import (

	"github.com/crossplane/crossplane-runtime/pkg/logging"
	fnv1beta1 "github.com/crossplane/function-sdk-go/proto/v1beta1"
	"github.com/crossplane/function-sdk-go/request"
	"github.com/crossplane/function-sdk-go/resource"
	"github.com/crossplane/function-sdk-go/response"
	"github.com/pkg/errors"
)

func getPausePolicy() []string {
	return []string{}
}

func getObservePolicy() []string {
	return []string{"Observe"}
}

func getFullControlPolicy() []string {
	return []string{"*"}
}

func getDefaultControlPolicy() []string {
	return []string{"Create", "Observe", "Update", "LateInitialize"}
}

func configureManagementPolicies(req *fnv1beta1.RunFunctionRequest, rsp *fnv1beta1.RunFunctionResponse, oxr *resource.Composite, log logging.Logger) (*fnv1beta1.RunFunctionRequest, *fnv1beta1.RunFunctionResponse) {
	desired, err := request.GetDesiredComposedResources(req)
	if err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot get desired composed resources from %T", req))
		return req, rsp
	}
	managementPolicy, _ := oxr.Resource.GetString("spec.managementPolicy")
	for _, dr := range desired {
		switch managementPolicy {
		case "pause":
			dr.Resource.SetValue("spec.managementPolicies", getPausePolicy())
		case "observe":
			dr.Resource.SetValue("spec.managementPolicies", getObservePolicy())
		case "full":
			dr.Resource.SetValue("spec.managementPolicies", getFullControlPolicy())
		default:
			dr.Resource.SetValue("spec.managementPolicies", getDefaultControlPolicy())
		}
	}
	if err := response.SetDesiredComposedResources(rsp, desired); err != nil {
		response.Fatal(rsp, errors.Wrapf(err, "cannot set desired composed resources in %T", rsp))
		return req, rsp
	}

	return req, rsp
}
