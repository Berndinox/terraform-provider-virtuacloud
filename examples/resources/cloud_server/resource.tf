# Look up offers, systems, and projects first
data "virtuacloud_offers" "available" {}

data "virtuacloud_systems" "available" {}

data "virtuacloud_projects" "available" {}

# Create a cloud server
resource "virtuacloud_cloud_server" "web" {
  project_uuid = data.virtuacloud_projects.available.projects[0].uuid
  offer_uuid   = data.virtuacloud_offers.available.offers[0].uuid
  system_uuid  = data.virtuacloud_systems.available.systems[0].uuid
  ipv6_enabled = true
  power_state  = "running"
  hostname     = "myserver.example.com"
}

# Create a server that starts in stopped state
resource "virtuacloud_cloud_server" "staging" {
  project_uuid = data.virtuacloud_projects.available.projects[0].uuid
  offer_uuid   = data.virtuacloud_offers.available.offers[0].uuid
  system_uuid  = data.virtuacloud_systems.available.systems[0].uuid
  ipv6_enabled = false
  power_state  = "stopped"
}