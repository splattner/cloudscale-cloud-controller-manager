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
	"fmt"
	"io"
	"os"

	cloudscale "github.com/cloudscale-ch/cloudscale-go-sdk"
	"golang.org/x/oauth2"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
)

const (
	cloudscaleTokenENVVar = "CLOUDSCALE_TOKEN"
	nodeNameENVVar        = "NODE_NAME"
	providerName          = "cloudscale"
)

type cloud struct {
	client    *cloudscale.Client
	instances cloudprovider.Instances
	zones     cloudprovider.Zones
}

func newCloud(config io.Reader) (cloudprovider.Interface, error) {
	token := os.Getenv(cloudscaleTokenENVVar)
	if token == "" {
		return nil, fmt.Errorf("environment variable %q is required", cloudscaleTokenENVVar)
	}
	nodeName := os.Getenv(nodeNameENVVar)
	if nodeName == "" {
		return nil, fmt.Errorf("environment variable %q is required", nodeNameENVVar)
	}

	token := &oauth2.Token{AccessToken: token}
	tokenSource := oauth2.StaticTokenSource(token)
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)

	client = cloudscale.NewClient(oauthClient)

	return &cloud{
		client:    client,
		zones:     newZones(client, nodeName),
		instances: newInstances(client),
	}, nil
}

func (c *cloud) Initialize(clientBuilder controller.ControllerClientBuilder) {}

func (c *cloud) Instances() (cloudprovider.Instances, bool) {
	return c.instances, true
}

func (c *cloud) Zones() (cloudprovider.Zones, bool) {
	return c.zones, true
}

func (c *cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return nil, false
}

func (c *cloud) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

func (c *cloud) Routes() (cloudprovider.Routes, bool) {
	return nil, false
}

func (c *cloud) ProviderName() string {
	return providerName
}

func (c *cloud) ScrubDNS(nameservers, searches []string) (nsOut, srchOut []string) {
	return nil, nil
}

func (c *cloud) HasClusterID() bool {
	return false
}

func init() {
	cloudprovider.RegisterCloudProvider(providerName, func(config io.Reader) (cloudprovider.Interface, error) {
		return newCloud(config)
	})
}
