package init

import (
	"os"
	"testing"

	backendAtlas "github.com/hashicorp/terraform/backend/atlas"
	backendLocal "github.com/hashicorp/terraform/backend/local"
	backendRemote "github.com/hashicorp/terraform/backend/remote"
	backendAzure "github.com/hashicorp/terraform/backend/remote-state/azure"
	backendConsul "github.com/hashicorp/terraform/backend/remote-state/consul"
	backendEtcdv3 "github.com/hashicorp/terraform/backend/remote-state/etcdv3"
	backendGCS "github.com/hashicorp/terraform/backend/remote-state/gcs"
	backendInmem "github.com/hashicorp/terraform/backend/remote-state/inmem"
	backendManta "github.com/hashicorp/terraform/backend/remote-state/manta"
	backendS3 "github.com/hashicorp/terraform/backend/remote-state/s3"
	backendSwift "github.com/hashicorp/terraform/backend/remote-state/swift"
)

func TestInit_backend(t *testing.T) {
	// Initialize the backends map
	Init(nil)

	backends := []string{
		"local",
		"remote",
		"atlas",
		"azurerm",
		"consul",
		"etcdv3",
		"gcs",
		"inmem",
		"manta",
		"s3",
		"swift",
		"azure",
	}

	// Make sure we get the requested backend
	for _, name := range backends {
		b := Backend(name)

		ok := false
		switch name {
		case "local":
			_, ok = b().(*backendLocal.Local)
		case "remote":
			_, ok = b().(*backendRemote.Remote)
		case "atlas":
			_, ok = b().(*backendAtlas.Backend)
		case "azurerm":
			_, ok = b().(*backendAzure.Backend)
		case "consul":
			_, ok = b().(*backendConsul.Backend)
		case "etcdv3":
			_, ok = b().(*backendEtcdv3.Backend)
		case "gcs":
			_, ok = b().(*backendGCS.Backend)
		case "inmem":
			_, ok = b().(*backendInmem.Backend)
		case "manta":
			_, ok = b().(*backendManta.Backend)
		case "s3":
			_, ok = b().(*backendS3.Backend)
		case "swift":
			_, ok = b().(*backendSwift.Backend)
		case "azure":
			_, ok = b().(deprecatedBackendShim)
		}

		if !ok {
			t.Fatalf("unable to assert the %q backend to its expected type", name)
		}
	}
}

func TestInit_forceLocalBackend(t *testing.T) {
	// Initialize the backends map
	Init(nil)

	enhancedBackends := []string{
		"local",
		"remote",
	}

	// Set the TF_FORCE_LOCAL_BACKEND flag so all enhanced backends will
	// return a local.Local backend with themselves as embedded backend.
	if err := os.Setenv("TF_FORCE_LOCAL_BACKEND", "1"); err != nil {
		t.Fatalf("error setting environment variable TF_FORCE_LOCAL_BACKEND: %v", err)
	}

	// Make sure we always get the local backend.
	for _, name := range enhancedBackends {
		b := Backend(name)

		local, ok := b().(*backendLocal.Local)
		if !ok {
			t.Fatalf("expected the %q enhanced backend to be of type \"*local.Local\"", name)
		}

		switch name {
		case "local":
			if local.Backend != nil {
				t.Fatalf("expected local.Backend to be nil, got: %T", local.Backend)
			}
		case "remote":
			if _, ok := local.Backend.(*backendRemote.Remote); !ok {
				t.Fatalf("expected local.Backend to be *remote.Remote, got: %T", local.Backend)
			}
		}
	}
}
