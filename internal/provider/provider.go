package provider

import (
	"context"
	"os"
	"time"

	"github.com/Berndinox/tf-provider-virtua-cloud/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	defaultBaseURL       = "https://api.virtua.cloud/v1"
	defaultTimeout       = 60
	serverCreateTimeout  = 10 * time.Minute
	serverDeleteTimeout  = 3 * time.Minute
	serverPowerTimeout   = 3 * time.Minute
	serverResizeTimeout  = 5 * time.Minute
	serverRestartTimeout = 3 * time.Minute
)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &virtuacloudProvider{version: version}
	}
}

type virtuacloudProvider struct {
	version string
}

var _ provider.Provider = (*virtuacloudProvider)(nil)

type virtuacloudProviderModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	BaseURL types.String `tfsdk:"base_url"`
	Timeout types.Int64  `tfsdk:"timeout"`
}

type virtuacloudProviderData struct {
	Client *client.Client
}

func (p *virtuacloudProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "virtuacloud"
	resp.Version = p.version
}

func (p *virtuacloudProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provider for Virtua.Cloud platform to manage cloud servers and resources.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "Virtua.Cloud API key. Can also be set via VIRTUACLOUD_API_KEY environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
			"base_url": schema.StringAttribute{
				Description: "Virtua.Cloud API base URL. Defaults to https://api.virtua.cloud/v1.",
				Optional:    true,
			},
			"timeout": schema.Int64Attribute{
				Description: "HTTP client timeout in seconds. Defaults to 60.",
				Optional:    true,
			},
		},
	}
}

func (p *virtuacloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config virtuacloudProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv("VIRTUACLOUD_API_KEY")
	if !config.APIKey.IsNull() && config.APIKey.IsUnknown() == false {
		apiKey = config.APIKey.ValueString()
	}
	if apiKey == "" {
		resp.Diagnostics.AddError("Missing API Key", "API key must be set via api_key attribute or VIRTUACLOUD_API_KEY environment variable")
		return
	}

	baseURL := defaultBaseURL
	if !config.BaseURL.IsNull() && config.BaseURL.IsUnknown() == false {
		baseURL = config.BaseURL.ValueString()
	}

	timeoutSeconds := int64(defaultTimeout)
	if !config.Timeout.IsNull() && config.Timeout.IsUnknown() == false {
		timeoutSeconds = config.Timeout.ValueInt64()
	}

	cl, err := client.NewClient(client.Config{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to create client", err.Error())
		return
	}

	providerData := &virtuacloudProviderData{Client: cl}
	resp.DataSourceData = providerData
	resp.ResourceData = providerData
}

func (p *virtuacloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccountDataSource,
		NewLimitsDataSource,
		NewProjectsDataSource,
		NewOffersDataSource,
		NewSystemsDataSource,
		NewCloudServerPasswordDataSource,
		NewCloudServersDataSource,
	}
}

func (p *virtuacloudProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCloudServerResource,
	}
}

func getClientFromDataSourceData(data any) *client.Client {
	providerData, ok := data.(*virtuacloudProviderData)
	if !ok {
		return nil
	}
	return providerData.Client
}
