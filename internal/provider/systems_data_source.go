package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type systemsDataSource struct {
	client *virtuacloudProviderData
}

var _ datasource.DataSource = (*systemsDataSource)(nil)

type systemsDataSourceModel struct {
	Systems []systemModel `tfsdk:"systems"`
}

type systemModel struct {
	UUID         types.String `tfsdk:"uuid"`
	Name         types.String `tfsdk:"name"`
	Distribution types.String `tfsdk:"distribution"`
	Version      types.String `tfsdk:"version"`
	Category     types.String `tfsdk:"category"`
	IsWindows    types.String `tfsdk:"is_windows"`
}

func NewSystemsDataSource() datasource.DataSource {
	return &systemsDataSource{}
}

func (d *systemsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "virtuacloud_systems"
}

func (d *systemsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all available operating systems for Virtua.Cloud servers.",
		Attributes: map[string]schema.Attribute{
			"systems": schema.ListNestedAttribute{
				Description: "List of available systems.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":         schema.StringAttribute{Description: "System UUID.", Computed: true},
						"name":         schema.StringAttribute{Description: "System display name.", Computed: true},
						"distribution": schema.StringAttribute{Description: "Distribution name (e.g. debian, ubuntu, windows-server-2019).", Computed: true},
						"version":      schema.StringAttribute{Description: "Distribution version.", Computed: true},
						"category":     schema.StringAttribute{Description: "System category (e.g. minimal, windows, apps).", Computed: true},
						"is_windows":   schema.StringAttribute{Description: "Whether this is a Windows system (1=yes, 0=no).", Computed: true},
					},
				},
			},
		},
	}
}

func (d *systemsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = nil
	if req.ProviderData != nil {
		d.client = req.ProviderData.(*virtuacloudProviderData)
	}
}

func (d *systemsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Provider client not configured")
		return
	}

	result, err := d.client.Client.GetSystems(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read systems", err.Error())
		return
	}

	var state systemsDataSourceModel
	for _, s := range result.Systems {
		state.Systems = append(state.Systems, systemModel{
			UUID:         types.StringValue(s.UUID),
			Name:         types.StringValue(s.Name),
			Distribution: types.StringValue(s.Distribution),
			Version:      types.StringValue(string(s.Version)),
			Category:     types.StringValue(s.Category),
			IsWindows:    types.StringValue(string(s.IsWindows)),
		})
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
