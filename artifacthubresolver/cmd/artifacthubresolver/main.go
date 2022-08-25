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

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	hub "github.com/tektoncd/resolution/artifacthubresolver/pkg/artifacthub"
	"github.com/tektoncd/resolution/pkg/resolver/framework"
	"knative.dev/pkg/injection/sharedmain"
)

func main() {
	apiURL := os.Getenv("ARTIFACT_HUB_API")
	hubURL := hub.DefaultArtifactHubURL
	if apiURL == "" {
		hubURL = hub.DefaultArtifactHubURL
	} else {
		if !strings.HasSuffix(apiURL, "/") {
			apiURL += "/"
		}
		hubURL = apiURL + hub.YamlEndpoint
	}
	fmt.Println("RUNNING WITH HUB URL PATTERN:", hubURL)
	resolver := hub.Resolver{ArtifactHubURL: hubURL}
	sharedmain.Main("artifacthub-controller",
		framework.NewController(context.Background(), &resolver),
	)
}
