/*
Copyright 2022 The Tekton Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package artifacthub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/tektoncd/resolution/pkg/common"
	"github.com/tektoncd/resolution/pkg/resolver/framework"
)

// LabelValueArtifactHubResolverType is the value to use for the resolution.tekton.dev/type label on resource requests.
const LabelValueArtifactHubResolverType string = "artifacthub"

// Resolver implements a framework.Resolver that can fetch files from Artifact Hub.
type Resolver struct {
	// ArtifactHubURL is the URL for artifact hub resolver.
	ArtifactHubURL string
}

// Initialize sets up any dependencies needed by the resolver.
// None at the moment.
func (r *Resolver) Initialize(context.Context) error {
	return nil
}

// GetName returns a string name to refer to this resolver by.
func (r *Resolver) GetName(context.Context) string {
	return "ArtifactHub"
}

// GetConfigName returns the name of the artifact hub resolver's configmap.
func (r *Resolver) GetConfigName(context.Context) string {
	return "artifacthubresolver-config"
}

// GetSelector returns a map of labels to match requests to the artifact hub resolver.
func (r *Resolver) GetSelector(context.Context) map[string]string {
	return map[string]string{
		common.LabelKeyResolverType: LabelValueArtifactHubResolverType,
	}
}

// ValidateParams ensures parameters from a request are as expected.
func (r *Resolver) ValidateParams(ctx context.Context, params map[string]string) error {
	if _, ok := params[ParamName]; !ok {
		return errors.New("must include name param")
	}
	if _, ok := params[ParamVersion]; !ok {
		return errors.New("must include version param")
	}
	if kind, ok := params[ParamKind]; ok {
		if kind != "task" && kind != "pipeline" {
			return errors.New("kind param must be task or pipeline")
		}
	}
	return nil
}

type dataResponse struct {
	YAML string `json:"manifestRaw"`
}

type hubResponse struct {
	Data dataResponse `json:"data"`
}

// Resolve uses the given params to resolve the requested file or resource.
func (r *Resolver) Resolve(ctx context.Context, params map[string]string) (framework.ResolvedResource, error) {
	kind, err := validateKind(ctx, params)
	if err != nil {
		return nil, err
	}
	params[ParamKind] = kind

	url := fmt.Sprintf(r.ArtifactHubURL, catalog(kind), repo(kind), params[ParamName], params[ParamVersion])
	return resolve(url)
}

func resolve(url string) (framework.ResolvedResource, error) {
	fmt.Println("FETCHING FROM ARTIFACT HUB URL:", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error requesting resource from hub: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}
	hr := hubResponse{}
	err = json.Unmarshal(body, &hr)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling json response: %w", err)
	}
	return &ResolvedArtifactHubResource{
		Content: []byte(hr.Data.YAML),
	}, nil
}

func validateKind(ctx context.Context, params map[string]string) (string, error) {
	conf := framework.GetResolverConfigFromContext(ctx)
	kind, ok := params[ParamKind]
	if !ok {
		if kindString, ok := conf[ConfigKind]; ok {
			kind = kindString
		} else {
			return "", fmt.Errorf("default resource Kind was not set during installation of the hub resolver")
		}
	}
	if kind != "task" && kind != "pipeline" {
		return "", fmt.Errorf("kind param must be task or pipeline")
	}
	return kind, nil
}

func catalog(kind string) string {
	return fmt.Sprintf("tekton-%s", kind)
}

func repo(kind string) string {
	return fmt.Sprintf("tekton-catalog-%ss", kind)
}

// ResolvedArtifactHubResource wraps the data we want to return to Pipelines
type ResolvedArtifactHubResource struct {
	Content []byte
}

var _ framework.ResolvedResource = &ResolvedArtifactHubResource{}

// Data returns the bytes of our hard-coded Pipeline
func (rr *ResolvedArtifactHubResource) Data() []byte {
	return rr.Content
}

// Annotations returns any metadata needed alongside the data. None atm.
func (*ResolvedArtifactHubResource) Annotations() map[string]string {
	return nil
}
