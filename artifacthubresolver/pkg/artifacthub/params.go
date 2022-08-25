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

// DefaultArtifactHubURL is the default url for the Artifact Hub API.
const DefaultArtifactHubURL = "https://artifacthub.io/api/v1/packages/%s/%s/%s/%s"

// YamlEndpoint is the suffix for a private custom artifact hub instance.
const YamlEndpoint = "api/v1/packages/%s/%s/%s/%s"

// ParamName is the parameter defining resource name in the catalog.
const ParamName = "name"

// ParamKind is the parameter defining the resource kind in the catalog.
const ParamKind = "kind"

// ParamVersion is the parameter defining the resource version in the catalog.
const ParamVersion = "version"
