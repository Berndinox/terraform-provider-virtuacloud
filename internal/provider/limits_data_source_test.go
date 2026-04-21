package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLimitsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLimitsDataSourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.virtuacloud_limits.test", "usage.cloud_servers"),
					resource.TestCheckResourceAttrSet("data.virtuacloud_limits.test", "limits.cloud_servers"),
				),
			},
		},
	})
}

func testAccLimitsDataSourceConfig() string {
	return `
data "virtuacloud_limits" "test" {}
`
}