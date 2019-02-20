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
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	cloudscale "github.com/cloudscale-ch/cloudscale-go-sdk"
	"golang.org/x/oauth2"
)

type testEnv struct {
	Server *httptest.Server
	Mux    *http.ServeMux
	Client *hcloud.Client
}

func (env *testEnv) Teardown() {
	env.Server.Close()
	env.Server = nil
	env.Mux = nil
	env.Client = nil
}

func newTestEnv() testEnv {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	token := &oauth2.Token{AccessToken: "token"}
	tokenSource := oauth2.StaticTokenSource(token)
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)

	client = cloudscale.NewClient(oauthClient)

	return testEnv{
		Server: server,
		Mux:    mux,
		Client: client,
	}
}

func TestNewCloud(t *testing.T) {
	os.Setenv("CLOUDSCALE_TOKEN", "test")
	os.Setenv("NODE_NAME", "test")

	var config bytes.Buffer
	_, err := newCloud(&config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestCloud(t *testing.T) {
	var config bytes.Buffer
	cloud, err := newCloud(&config)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	t.Run("Instances", func(t *testing.T) {
		_, supported := cloud.Instances()
		if !supported {
			t.Error("Instances interface should be supported")
		}
	})

	t.Run("Zones", func(t *testing.T) {
		_, supported := cloud.Zones()
		if !supported {
			t.Error("Zones interface should be supported")
		}
	})

	t.Run("LoadBalancer", func(t *testing.T) {
		_, supported := cloud.LoadBalancer()
		if supported {
			t.Error("LoadBalancer interface should not be supported")
		}
	})

	t.Run("Clusters", func(t *testing.T) {
		_, supported := cloud.Clusters()
		if supported {
			t.Error("Clusters interface should not be supported")
		}
	})

	t.Run("Routes", func(t *testing.T) {
		_, supported := cloud.Routes()
		if supported {
			t.Error("Routes interface should not be supported")
		}
	})

	t.Run("HasClusterID", func(t *testing.T) {
		if cloud.HasClusterID() {
			t.Error("HasClusterID should be false")
		}
	})

	t.Run("ProviderName", func(t *testing.T) {
		if cloud.ProviderName() != "hcloud" {
			t.Error("ProviderName should be hcloud")
		}
	})
}
