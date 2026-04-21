# terraform-provider-virtuacloud

OpenTofu/Terraform provider for [Virtua.Cloud](https://www.virtua.cloud) platform. Manage cloud servers and resources using Infrastructure as Code.

**Note:**
* The project is in its eraly stages (**Beta**)
* This is an unofficial Terraform provider and is not affiliated with, endorsed by, or supported by Virtua.Cloud / HashiCorp
* Provided “as is” without warranty of any kind; use at your own risk.
* Contributions are welcome and appreciated.


## Requirements

- [OpenTofu](https://opentofu.org/) >= 1.6 or [Terraform](https://www.terraform.io/) >= 1.5
- Go >= 1.21 (for building from source)

## Quick Start

```hcl
terraform {
  required_providers {
    virtuacloud = {
      source  = "registry.opentofu.org/Berndinox/virtuacloud"
      version = "~> 0.1"
    }
  }
}

provider "virtuacloud" {
  api_key = var.virtuacloud_api_key
}

data "virtuacloud_offers" "available" {}

data "virtuacloud_systems" "available" {}

resource "virtuacloud_cloud_server" "web" {
  project_uuid  = "your-project-uuid"
  offer_uuid    = data.virtuacloud_offers.available.offers[0].uuid
  system_uuid   = data.virtuacloud_systems.available.systems[0].uuid
  ipv6_enabled  = true
  power_state   = "running"
}
```

## Authentication

Set your Virtua.Cloud API key either in the provider configuration or via the `VIRTUACLOUD_API_KEY` environment variable:

```hcl
provider "virtuacloud" {
  api_key = var.virtuacloud_api_key  # Recommended: use a variable
}
```

```bash
export VIRTUACLOUD_API_KEY="your-api-key"
```

## Resources

### `virtuacloud_cloud_server`

Manages a cloud server. Creation may take 1-2 minutes.

| Attribute | Type | Required | ForceNew | Description |
|---|---|---|---|---|
| `project_uuid` | string | yes | yes | Project UUID |
| `offer_uuid` | string | yes | no | Offer UUID (change triggers resize) |
| `system_uuid` | string | yes | yes | OS system UUID |
| `ipv6_enabled` | bool | yes | yes | Enable IPv6 |
| `ipv6_block_size` | int | no | yes | IPv6 block size (128=/128, 64=/64) |
| `hostname` | string | no | no | Custom hostname |
| `power_state` | string | no | no | `running` or `stopped` (default: running) |
| `resize_disk` | bool | no | no | Resize disk on offer change (default: false) |
| `restart_triggered_at` | string | no | no | Trigger restart by updating this value |
| `uuid` | string | computed | - | Server UUID |
| `name` | string | computed | - | Auto-generated name |
| `status` | string | computed | - | Server status |
| `vcpus` | string | computed | - | Number of vCPUs |
| `memory_size` | string | computed | - | Memory in MB |
| `root_space` | string | computed | - | Disk in GB |
| `root_disk_type` | string | computed | - | Disk type (ssd/nvme) |
| `offer` | object | computed | - | Current offer details |
| `cloud_zone` | object | computed | - | Zone/location details |
| `system` | object | computed | - | OS details |

**Import:** `terraform import virtuacloud_cloud_server.example <uuid>`

## Data Sources

### `virtuacloud_account`
Account info: balance, currency, usage.

### `virtuacloud_limits`
Resource usage and limits.

### `virtuacloud_projects`
List all projects.

### `virtuacloud_offers`
List available server offers.

### `virtuacloud_systems`
List available operating systems.

### `virtuacloud_cloud_server_password`
Retrieve root or rescue password (sensitive).

| Argument | Type | Required | Description |
|---|---|---|---|
| `server_uuid` | string | yes | Server UUID |
| `password_type` | string | yes | `root` or `rescue` |

## Development

```bash
make build     # Build the provider binary
make install   # Install locally
make test      # Run unit tests
make docs      # Generate documentation
```
This project was developed making use of AI.

## License

MIT
