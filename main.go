//go:build darwin
// +build darwin

/*
Copyright 2016 The Kubernetes Authors All rights reserved.
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

	"github.com/densityops/machine-driver-hyperkit/pkg/hyperkit"
	"github.com/densityops/machine/libmachine/drivers/plugin"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "version" {
			fmt.Printf("Driver version: %s\n", hyperkit.DriverVersion)
			os.Exit(0)
		}
	}
	plugin.RegisterDriver(hyperkit.NewDriver())
}
