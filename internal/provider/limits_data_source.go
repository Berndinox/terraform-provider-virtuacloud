package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type limitsDataSource struct {
	client *virtuacloudProviderData
}

var _ datasource.DataSource = (*limitsDataSource)(nil)

type limitsDataSourceModel struct {
	Usage  limitsUsageModel `tfsdk:"usage"`
	Limits limitsInfoModel  `tfsdk:"limits"`
}

type limitsUsageModel struct {
	CloudServers types.Int64  `tfsdk:"cloud_servers"`
	Vcpus        types.String `tfsdk:"vcpus"`
	MemorySize   types.String `tfsdk:"memory_size"`
	RootSpace    types.String `tfsdk:"root_space"`
	IpAddressV4  types.String `tfsdk:"ip_address_v4"`
	IpAddressV6  types.String `tfsdk:"ip_address_v6"`
}

type limitsInfoModel struct {
	CloudServers types.String `tfsdk:"cloud_servers"`
	SmtpEnabled  types.String `tfsdk:"smtp_enabled"`
}

func NewLimitsDataSource() datasource.DataSource {
	return &limitsDataSource{}
}

func (d *limitsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "virtuacloud_limits"
}

func (d *limitsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve resource usage and limits for the authenticated Virtua.Cloud account.",
		Attributes: map[string]schema.Attribute{
			"usage": schema.SingleNestedAttribute{
				Description: "Current resource usage.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"cloud_servers": schema.Int64Attribute{
						Description: "Number of active cloud servers.",
						Computed:    true,
					},
					"vcpus": schema.StringAttribute{
						Description: "Number of vCPUs in use.",
						Computed:    true,
					},
					"memory_size": schema.StringAttribute{
						Description: "Memory size in MB in use.",
						Computed:    true,
					},
					"root_space": schema.StringAttribute{
						Description: "Root disk space in GB in use.",
						Computed:    true,
					},
					"ip_address_v4": schema.StringAttribute{
						Description: "Number of IPv4 addresses in use.",
						Computed:    true,
					},
					"ip_address_v6": schema.StringAttribute{
						Description: "Number of IPv6 addresses in use.",
						Computed:    true,
					},
				},
			},
			"limits": schema.SingleNestedAttribute{
				Description: "Resource limits for the account.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"cloud_servers": schema.StringAttribute{
						Description: "Maximum number of cloud servers allowed.",
						Computed:    true,
					},
					"smtp_enabled": schema.StringAttribute{
						Description: "Whether SMTP is enabled (1=enabled, 0=disabled).",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (d *limitsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = nil
	if req.ProviderData != nil {
		d.client = req.ProviderData.(*virtuacloudProviderData)
	}
}

func (d *limitsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Provider client not configured")
		return
	}

	limits, err := d.client.Client.GetLimits(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read limits", err.Error())
		return
	}

	var state limitsDataSourceModel
	state.Usage = limitsUsageModel{
		CloudServers: types.Int64Value(int64(limits.Usage.CloudServers)),
		Vcpus:        types.StringValue(string(limits.Usage.Vcpus)),
		MemorySize:   types.StringValue(string(limits.Usage.MemorySize)),
		RootSpace:    types.StringValue(string(limits.Usage.RootSpace)),
		IpAddressV4:  types.StringValue(string(limits.Usage.IpAddressV4)),
		IpAddressV6:  types.StringValue(string(limits.Usage.IpAddressV6)),
	}
	state.Limits = limitsInfoModel{
		CloudServers: types.StringValue(string(limits.Limits.CloudServers)),
		SmtpEnabled:  types.StringValue(fmt.Sprintf("%d", limits.Limits.SmtpEnabled)),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
