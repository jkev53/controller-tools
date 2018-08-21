/*
Copyright 2018 The Kubernetes Authors.

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

package external

import (
	"fmt"
	"regexp"
	// "strings"

	"github.com/markbates/inflect"
	"sigs.k8s.io/controller-tools/pkg/scaffold/resource"
)

// External scaffolds a External for a Resource
type External struct {
	Resource *resource.Resource

	// Domain the domain of the external type
	Domain string

	// ImportPath the import path of the external type
	ImportPath string
}

// Validate checks the Resource values to make sure they are valid.
func (e *External) Validate() error {
	if len(e.Resource.Group) == 0 {
		return fmt.Errorf("group cannot be empty")
	}
	if len(e.Resource.Version) == 0 {
		return fmt.Errorf("version cannot be empty")
	}
	if len(e.Resource.Kind) == 0 {
		return fmt.Errorf("kind cannot be empty")
	}

	// rs := inflect.NewDefaultRuleset()
	// if len(e.Resource.Resource) == 0 {
	// 	e.Resource.Resource = rs.Pluralize(strings.ToLower(e.Resource.Kind))
	// }

	groupMatch := regexp.MustCompile("^[a-z]+$")
	if !groupMatch.MatchString(e.Resource.Group) {
		return fmt.Errorf("group must match ^[a-z]+$ (was %s)", e.Resource.Group)
	}

	versionMatch := regexp.MustCompile("^v\\d+(alpha\\d+|beta\\d+)?$")
	if !versionMatch.MatchString(e.Resource.Version) {
		return fmt.Errorf(
			"version must match ^v\\d+(alpha\\d+|beta\\d+)?$ (was %s)", e.Resource.Version)
	}

	if e.Resource.Kind != inflect.Camelize(e.Resource.Kind) {
		return fmt.Errorf("Kind must be camelcase (expected %s was %s)", inflect.Camelize(e.Resource.Kind), e.Resource.Kind)
	}

	return nil
}
