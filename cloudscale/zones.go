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
	cloudscale "github.com/cloudscale-ch/cloudscale-go-sdk"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

type zones struct {
	client   *cloudscale.Client
	nodeName string // name of the node the programm is running on
}

func newZones(client *cloudscale.Client, nodeName string) cloudprovider.Zones {
	return zones{client, nodeName}
}

func (z zones) GetZone() (zone cloudprovider.Zone, err error) {
	var server *cloudscale.Server
	server, err = getServerByName(z.client, z.nodeName)
	if err != nil {
		return
	}
	zone = zoneFromServer(server)
	return
}

func (z zones) GetZoneByProviderID(providerID string) (zone cloudprovider.Zone, err error) {
	var server *cloudscale.Server
	server, err = getServerByID(z.client, providerID)
	if err != nil {
		return
	}
	zone = zoneFromServer(server)
	return
}

func (z zones) GetZoneByNodeName(nodeName types.NodeName) (zone cloudprovider.Zone, err error) {
	var server *cloudscale.Server
	server, err = getServerByName(z.client, string(nodeName))
	if err != nil {
		return
	}
	zone = zoneFromServer(server)
	return
}

func zoneFromServer(server *cloudscale.Server) (zone cloudprovider.Zone) {
	return cloudprovider.Zone{
		Region:        "default",
		FailureDomain: "default",
	}
}
