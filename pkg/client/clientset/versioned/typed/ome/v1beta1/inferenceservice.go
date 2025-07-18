// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	context "context"

	omev1beta1 "github.com/sgl-project/ome/pkg/apis/ome/v1beta1"
	scheme "github.com/sgl-project/ome/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// InferenceServicesGetter has a method to return a InferenceServiceInterface.
// A group's client should implement this interface.
type InferenceServicesGetter interface {
	InferenceServices(namespace string) InferenceServiceInterface
}

// InferenceServiceInterface has methods to work with InferenceService resources.
type InferenceServiceInterface interface {
	Create(ctx context.Context, inferenceService *omev1beta1.InferenceService, opts v1.CreateOptions) (*omev1beta1.InferenceService, error)
	Update(ctx context.Context, inferenceService *omev1beta1.InferenceService, opts v1.UpdateOptions) (*omev1beta1.InferenceService, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, inferenceService *omev1beta1.InferenceService, opts v1.UpdateOptions) (*omev1beta1.InferenceService, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*omev1beta1.InferenceService, error)
	List(ctx context.Context, opts v1.ListOptions) (*omev1beta1.InferenceServiceList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *omev1beta1.InferenceService, err error)
	InferenceServiceExpansion
}

// inferenceServices implements InferenceServiceInterface
type inferenceServices struct {
	*gentype.ClientWithList[*omev1beta1.InferenceService, *omev1beta1.InferenceServiceList]
}

// newInferenceServices returns a InferenceServices
func newInferenceServices(c *OmeV1beta1Client, namespace string) *inferenceServices {
	return &inferenceServices{
		gentype.NewClientWithList[*omev1beta1.InferenceService, *omev1beta1.InferenceServiceList](
			"inferenceservices",
			c.RESTClient(),
			scheme.ParameterCodec,
			namespace,
			func() *omev1beta1.InferenceService { return &omev1beta1.InferenceService{} },
			func() *omev1beta1.InferenceServiceList { return &omev1beta1.InferenceServiceList{} },
		),
	}
}
