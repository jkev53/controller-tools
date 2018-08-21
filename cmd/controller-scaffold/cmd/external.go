// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"
	// "os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
	"sigs.k8s.io/controller-tools/pkg/scaffold"
	"sigs.k8s.io/controller-tools/pkg/scaffold/external"
	"sigs.k8s.io/controller-tools/pkg/scaffold/input"
	"sigs.k8s.io/controller-tools/pkg/scaffold/resource"
)

var e *external.External

var extResourceFlag, extControllerFlag *flag.Flag

// var doResource, doController, doMake bool

// ExternalCmd represents the resource command
var ExternalCmd = &cobra.Command{
	Use:   "external",
	Short: "Scaffold a Kubernetes API",
	Long: `Scaffold a Kubernetes API by creating a Resource definition and / or a Controller.

api will prompt the user for if it should scaffold the Resource and / or Controller.  To only
scaffold a Controller for an existing Resource, select "n" for Resource.  To only define
the schema for a Resource without writing a Controller, select "n" for Controller.

After the scaffold is written, api will run make on the project.
`,
	Example: `	# Create a frigates API with Group: ship, Version: v1beta1 and Kind: Frigate
	controller-scaffold external --group ship --version v1beta1 --kind Frigate

	# Edit the Controller
	nano pkg/controller/frigate/frigate_controller.go

	# Edit the Controller Test
	nano pkg/controller/frigate/frigate_controller_test.go

	# Regenerate code and run against the Kubernetes cluster configured by ~/.kube/config
	make run
`,
	Run: func(cmd *cobra.Command, args []string) {
		DieIfNoProjectx()

		if !extResourceFlag.Changed {
			fmt.Println("Create Resource under pkg/apis [y/n]?")
			doResource = yesno()
		}
		if !extControllerFlag.Changed {
			// fmt.Println("Create Controller under pkg/controller [y/n]?")
			doController = yesno()
		}

		fmt.Println("Writing scaffold for you to edit...")

		if doResource {
			fmt.Println(filepath.Join("pkg", "apis", e.Resource.Group, e.Resource.Version,
				fmt.Sprintf("%s_types.go", strings.ToLower(e.Resource.Kind))))
			fmt.Println(filepath.Join("pkg", "apis", e.Resource.Group, e.Resource.Version,
				fmt.Sprintf("%s_types_test.go", strings.ToLower(e.Resource.Kind))))

			err := (&scaffold.Scaffold{}).Execute(input.Options{},
				&external.AddToScheme{External: e},
			)
			if err != nil {
				log.Fatal(err)
			}
		}

		if doController {
			fmt.Println(filepath.Join("pkg", "controller", strings.ToLower(e.Resource.Kind),
				fmt.Sprintf("%s_controller.go", strings.ToLower(e.Resource.Kind))))
			fmt.Println(filepath.Join("pkg", "apis", strings.ToLower(e.Resource.Kind),
				fmt.Sprintf("%s_controller_test.go", strings.ToLower(e.Resource.Kind))))

			// resrc := &resource.Resource{
			// 	Namespaced = e.Namespaced,
			// 	Group = e.Group,
			// 	Version = e.Version,
			// 	Kind = e.Kind,
			// 	Resource = e.Resource.Resource,
			// 	ShortNames = e.ShortNames,
			// 	CreateExampleReconcileBody = e.CreateExampleReconcileBody,
			// }
			err := (&scaffold.Scaffold{}).Execute(input.Options{},
				&external.Controller{External: e},
				&external.AddController{External: e},
				// &controller.Test{External: e},
				// &controller.SuiteTest{External: e},
			)
			if err != nil {
				log.Fatal(err)
			}
		}

		// if doMake {
		// 	fmt.Println("Running make...")
		// 	cm := exec.Command("make") // #nosec
		// 	cm.Stderr = os.Stderr
		// 	cm.Stdout = os.Stdout
		// 	if err := cm.Run(); err != nil {
		// 		log.Fatal(err)
		// 	}
		// }
	},
}

func init() {
	rootCmd.AddCommand(ExternalCmd)
	ExternalCmd.Flags().BoolVar(&doMake, "make", true,
		"if true, run make after generating files")
	ExternalCmd.Flags().BoolVar(&doResource, "resource", true,
		"if set, generate the resource without prompting the user")
	extResourceFlag = ExternalCmd.Flag("resource")
	ExternalCmd.Flags().BoolVar(&doController, "controller", true,
		"if set, generate the controller without prompting the user")
	extControllerFlag = ExternalCmd.Flag("controller")
	e = ExternalForFlags(ExternalCmd.Flags())
}

// DieIfNoProjectx checks to make sure the command is run from a directory containing a project file.
func DieIfNoProjectx() {
	if _, err := os.Stat("PROJECT"); os.IsNotExist(err) {
		log.Fatalf("Command must be run from a diretory containing %s", "PROJECT")
	}
}

// ExternalForFlags registers flags for Resource fields and returns the Resource
func ExternalForFlags(f *flag.FlagSet) *external.External {
	e := &external.External{}
	r := &resource.Resource{}
	f.StringVar(&r.Kind, "kind", "", "resource Kind")
	f.StringVar(&r.Group, "group", "", "resource Group")
	f.StringVar(&r.Version, "version", "", "resource Version")
	f.BoolVar(&r.Namespaced, "namespaced", true, "true if the resource is namespaced")
	f.BoolVar(&r.CreateExampleReconcileBody, "example", true,
		"true if an example reconcile body should be written")

	f.StringVar(&e.Domain, "domain", "", "external types domain")
	f.StringVar(&e.ImportPath, "importpath", "", "external types importPath")

	e.Resource = r
	return e
}
