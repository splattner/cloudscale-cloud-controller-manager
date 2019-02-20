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
	"errors"
	"strconv"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudscale "github.com/cloudscale-ch/cloudscale-go-sdk"
)

type instances struct {
	client *cloudscale.Client
}

func newInstances(client *cloudscale.Client) *instances {
	return &instances{client}
}

func (i *instances) NodeAddressesByProviderID(providerID string) ([]v1.NodeAddress, error) {


	server, err := getServerByID(i.client, providerID)
	if err != nil {
		return nil, err
	}
	return nodeAddresses(server), nil
}

func (i *instances) NodeAddresses(nodeName types.NodeName) ([]v1.NodeAddress, error) {
	server, err := getServerByName(i.client, string(nodeName))
	if err != nil {
		return nil, err
	}
	return nodeAddresses(server), nil
}

func (i *instances) ExternalID(nodeName types.NodeName) (string, error) {
	return i.InstanceID(nodeName)
}

func (i *instances) InstanceID(nodeName types.NodeName) (string, error) {
	server, err := getServerByName(i.client, string(nodeName))
	if err != nil {
		return "", err
	}
	return strconv.Itoa(server.uuid), nil
}

func (i *instances) InstanceType(nodeName types.NodeName) (string, error) {
	server, err := getServerByName(i.client, string(nodeName))
	if err != nil {
		return "", err
	}
	return server.Flavor, nil
}

func (i *instances) InstanceTypeByProviderID(providerID string) (string, error) {
	server, err := getServerByID(i.client, providerID)
	if err != nil {
		return "", err
	}
	return server.Flavor, nil
}

func (i *instances) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	return errors.New("not implemented")
}

func (i *instances) CurrentNodeName(hostname string) (types.NodeName, error) {
	return types.NodeName(hostname), nil
}

func (i instances) InstanceExistsByProviderID(providerID string) (exists bool, err error) {

	var server *cloudscale.Server
	server, _, err = i.client.Server.Get(context.Background(), providerID)
	if err != nil {
		return
	}
	exists = server != nil
	return
}

func nodeAddresses(server *cloudscale.Server) []v1.NodeAddress {
	var addresses []v1.NodeAddress

	var publicIP = ""

	for iface := range server.Interfaces {
		if iface.Type == "public" {
			for address : range iface.Addresses {
				if address.version == 4 {
					publicIP = address.address
				}
			}
		}
	}

	addresses = append(
		addresses,
		v1.NodeAddress{Type: v1.NodeHostName, Address: server.Name},
		v1.NodeAddress{Type: v1.NodeExternalIP, Address: publicIP},
	)
	return addresses
}
