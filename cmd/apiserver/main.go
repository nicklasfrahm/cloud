/*
Copyright 2025.

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
	"fmt"
	"os"

	"k8s.io/klog"
	"sigs.k8s.io/apiserver-runtime/pkg/builder"

	// +kubebuilder:scaffold:resource-imports
	cloudv1 "github.com/nicklasfrahm/cloud/pkg/apis/cloud/v1"
	"github.com/nicklasfrahm/cloud/pkg/storage/blob"
	"github.com/spf13/pflag"
)

func main() {
	configPath := pflag.String("config-path", "", "storage prefix for the object storage")

	objectStorageConfig, err := os.ReadFile(*configPath)
	if err != nil {
		klog.Fatal(fmt.Errorf("failed to read object storage config: %v", err))
	}

	err = builder.APIServer.
		// +kubebuilder:scaffold:resource-register
		WithResourceAndHandler(&cloudv1.Machine{}, blob.NewJSONBLOBStorageProvider(&cloudv1.Machine{}, objectStorageConfig)).
		WithoutEtcd().
		Execute()
	if err != nil {
		klog.Fatal(err)
	}
}
