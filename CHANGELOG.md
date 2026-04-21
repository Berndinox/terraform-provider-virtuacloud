# Changelog

## 0.1.0 (Unreleased)

### Features

- **Provider**: Initial Virtua.Cloud provider with API key authentication
- **Resource**: `virtuacloud_cloud_server` - Full lifecycle management (create, read, update, delete, import)
  - Power state management (`running`/`stopped`)
  - Server resize via offer change
  - Restart trigger via `restart_triggered_at`
  - IPv6 support with configurable block size
  - Custom hostname support
  - Automatic polling for server setup completion (1-2 min)
- **Data Source**: `virtuacloud_account` - Account info (balance, usage, limits)
- **Data Source**: `virtuacloud_limits` - Resource usage and limits
- **Data Source**: `virtuacloud_projects` - Project listing
- **Data Source**: `virtuacloud_offers` - Available server offers
- **Data Source**: `virtuacloud_systems` - Available operating systems
- **Data Source**: `virtuacloud_cloud_server_password` - Server password retrieval (sensitive)

### Security

- API key marked as sensitive in provider configuration
- Password data source uses sensitive attribute
- HTTPS-only communication with Virtua.Cloud API
- Exponential backoff retry on transient errors
- No credentials stored in resource state