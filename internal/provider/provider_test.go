package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/Berndinox/tf-provider-virtua-cloud/internal/provider"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"virtuacloud": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("VIRTUACLOUD_API_KEY"); v == "" {
		t.Fatal("VIRTUACLOUD_API_KEY must be set for acceptance tests")
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
