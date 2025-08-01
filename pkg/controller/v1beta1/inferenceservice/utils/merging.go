package utils

import (
	"bytes"
	"encoding/json"
	"html/template"

	goerrors "github.com/pkg/errors"
	"github.com/sgl-project/ome/pkg/apis/ome/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

// MergeRuntimeContainers Merge the predictor Container struct with the runtime Container struct, allowing users
func MergeRuntimeContainers(runtimeContainer *v1.Container, predictorContainer *v1.Container) (*v1.Container, error) {
	// Save runtime container name, as the name can be overridden as empty string during the Unmarshal below
	// since the Name field does not have the 'omitempty' struct tag.
	runtimeContainerName := runtimeContainer.Name

	// Use JSON Marshal/Unmarshal to merge Container structs using strategic merge patch
	runtimeContainerJson, err := json.Marshal(runtimeContainer)
	if err != nil {
		return nil, err
	}

	overrides, err := json.Marshal(predictorContainer)
	if err != nil {
		return nil, err
	}

	mergedContainer := v1.Container{}
	jsonResult, err := strategicpatch.StrategicMergePatch(runtimeContainerJson, overrides, mergedContainer)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonResult, &mergedContainer); err != nil {
		return nil, err
	}

	if mergedContainer.Name == "" {
		mergedContainer.Name = runtimeContainerName
	}

	// Strategic merge patch will replace args but more useful behaviour here is to concatenate
	mergedContainer.Args = append(append([]string{}, runtimeContainer.Args...), predictorContainer.Args...)

	return &mergedContainer, nil
}

// mergeSpec merges two Kubernetes-style specs using strategic merge patch.
// `runtimeInit` is the base (typically from a runtime default), and `override` comes from the user-defined resource.
// Fields from `override` will overwrite corresponding fields in `runtimeInit`.
func mergeSpec[T any](runtimeInit T, override T) (*T, error) {
	baseJSON, err := json.Marshal(runtimeInit)
	if err != nil {
		return nil, err
	}

	overrideJSON, err := json.Marshal(override)
	if err != nil {
		return nil, err
	}

	var merged T
	jsonResult, err := strategicpatch.StrategicMergePatch(baseJSON, overrideJSON, merged)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonResult, &merged); err != nil {
		return nil, err
	}

	return &merged, nil
}

// MergeRouterSpec merges a runtime-provided RouterSpec with a user-provided RouterSpec from InferenceService.
// The user-provided (isvcRouter) fields take precedence over the runtime defaults.
func MergeRouterSpec(isvcRouter, runtimeRouter *v1beta1.RouterSpec) (*v1beta1.RouterSpec, error) {
	switch {
	case isvcRouter == nil:
		// if router is not specified in isvc, return nil
		return nil, nil
	case runtimeRouter == nil:
		// if router is not specified in runtime, return a copy of isvcRouter
		return isvcRouter.DeepCopy(), nil
	}

	return mergeSpec(v1beta1.RouterSpec{
		ComponentExtensionSpec: runtimeRouter.ComponentExtensionSpec,
		PodSpec:                runtimeRouter.PodSpec,
		Runner:                 runtimeRouter.Runner,
		Config:                 runtimeRouter.Config,
	}, *isvcRouter)
}

// MergeEngineSpec merges a runtime-provided EngineSpec with a user-provided EngineSpec from InferenceService.
// The user-provided (isvcEngine) fields take precedence over the runtime defaults.
func MergeEngineSpec(runtimeEngine, isvcEngine *v1beta1.EngineSpec) (*v1beta1.EngineSpec, error) {
	switch {
	case isvcEngine == nil:
		// if engine is not specified in isvc, return nil
		return nil, nil
	case runtimeEngine == nil:
		// if engine is not specified in runtime, return a copy of isvcEngine
		return isvcEngine.DeepCopy(), nil
	}

	return mergeSpec(v1beta1.EngineSpec{
		ComponentExtensionSpec: runtimeEngine.ComponentExtensionSpec,
		PodSpec:                runtimeEngine.PodSpec,
		Runner:                 runtimeEngine.Runner,
		Leader:                 runtimeEngine.Leader,
		Worker:                 runtimeEngine.Worker,
	}, *isvcEngine)
}

// MergeDecoderSpec merges a runtime-provided DecoderSpec with a user-provided DecoderSpec from InferenceService.
// The user-provided (isvcDecoder) fields take precedence over the runtime defaults.
func MergeDecoderSpec(runtimeDecoder, isvcDecoder *v1beta1.DecoderSpec) (*v1beta1.DecoderSpec, error) {
	switch {
	case isvcDecoder == nil:
		// if decoder is not specified in isvc, return nil
		return nil, nil
	case runtimeDecoder == nil:
		// if decoder is not specified in runtime, return a copy of isvcDecoder
		return isvcDecoder.DeepCopy(), nil
	}

	return mergeSpec(v1beta1.DecoderSpec{
		ComponentExtensionSpec: runtimeDecoder.ComponentExtensionSpec,
		PodSpec:                runtimeDecoder.PodSpec,
		Runner:                 runtimeDecoder.Runner,
		Leader:                 runtimeDecoder.Leader,
		Worker:                 runtimeDecoder.Worker,
	}, *isvcDecoder)
}

// ConvertPodSpec converts v1beta1.PodSpec to v1.PodSpec
// This handles the conversion between the custom v1beta1.PodSpec type and the core v1.PodSpec type
func ConvertPodSpec(spec *v1beta1.PodSpec) (*v1.PodSpec, error) {
	if spec == nil {
		return nil, goerrors.New("cannot convert nil PodSpec")
	}

	// Use JSON marshaling to convert between the types
	data, err := json.Marshal(spec)
	if err != nil {
		return nil, goerrors.Wrap(err, "failed to marshal v1beta1.PodSpec")
	}

	var podSpec v1.PodSpec
	if err := json.Unmarshal(data, &podSpec); err != nil {
		return nil, goerrors.Wrap(err, "failed to unmarshal to v1.PodSpec")
	}

	return &podSpec, nil
}

// ReplacePlaceholders Replace placeholders in runtime container by values from inferenceservice metadata
func ReplacePlaceholders(container *v1.Container, meta metav1.ObjectMeta) error {
	data, _ := json.Marshal(container)
	tmpl, err := template.New("container-tmpl").Parse(string(data))
	if err != nil {
		return err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, meta)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf.Bytes(), container)
}
