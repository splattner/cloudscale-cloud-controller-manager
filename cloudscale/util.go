/*
Copyright 2018 Hetzner Cloud GmbH.

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

package cloudscale

import (
	"context"
	"strings"

	cloudscale "github.com/cloudscale-ch/cloudscale-go-sdk"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

func getServerByName(c *cloudscale.Client, name string) (server *cloudscale.Server, err error) {

	allServer, err := c.Servers.List(context.Background())

	for i := range allServer {
		if strings.Compare(allServer[i].Name, name) == 0 {
			server = &allServer[i]
			break
		}
	}

	if err != nil {
		return
	}
	if server == nil {
		err = cloudprovider.InstanceNotFound
		return
	}
	return
}

func getServerByID(c *cloudscale.Client, id string) (server *cloudscale.Server, err error) {
	server, err = c.Servers.Get(context.Background(), id)
	if err != nil {
		return
	}
	if server == nil {
		err = cloudprovider.InstanceNotFound
		return
	}
	return
}
