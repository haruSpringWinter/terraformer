// Copyright 2018 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package datadog

import (
	"context"
	"fmt"

	datadogV2 "github.com/DataDog/datadog-api-client-go/api/v2/datadog"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
)

var (
	// RoleAllowEmptyValues ...
	RoleAllowEmptyValues = []string{}
)

// RoleGenerator ...
type RoleGenerator struct {
	DatadogService
}

func (g *RoleGenerator) createResources(roles []datadogV2.Role) []terraformutils.Resource {
	resources := []terraformutils.Resource{}
	for _, role := range roles {
		resourceName := role.GetId()
		resource := g.createResource(resourceName)
		resource.IgnoreKeys = append(resource.IgnoreKeys, "permission.([0-9]+).name")
		resources = append(resources, resource)
	}

	return resources
}

func (g *RoleGenerator) createResource(roleID string) terraformutils.Resource {
	return terraformutils.NewSimpleResource(
		roleID,
		fmt.Sprintf("role_%s", roleID),
		"datadog_role",
		"datadog",
		RoleAllowEmptyValues,
	)
}

// InitResources Generate TerraformResources from Datadog API,
// from each role create 1 TerraformResource.
// Need Role ID as ID for terraform resource
func (g *RoleGenerator) InitResources() error {
	datadogClientV1 := g.Args["datadogClientV2"].(*datadogV2.APIClient)
	authV1 := g.Args["authV2"].(context.Context)

	roles, _, err := datadogClientV1.RolesApi.ListRoles(authV1).Execute()
	if err != nil {
		return err
	}
	g.Resources = g.createResources(roles.GetData())
	return nil
}
