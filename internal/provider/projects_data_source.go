package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type projectsDataSource struct {
	client *virtuacloudProviderData
}

var _ datasource.DataSource = (*projectsDataSource)(nil)

type projectsDataSourceModel struct {
	Projects []projectModel `tfsdk:"projects"`
}

type projectModel struct {
	UUID                      types.String `tfsdk:"uuid"`
	Name                      types.String `tfsdk:"name"`
	Description               types.String `tfsdk:"description"`
	CloudServersCount         types.String `tfsdk:"cloud_servers_count"`
	MonthlyCloudUsage         types.String `tfsdk:"monthly_cloud_usage"`
	MonthlyCloudUsageEstimate types.String `tfsdk:"monthly_cloud_usage_estimate"`
	DomainsCount              types.String `tfsdk:"domains_count"`
	Environment               types.String `tfsdk:"environment"`
	CreatedAt                 types.String `tfsdk:"created_at"`
}

func NewProjectsDataSource() datasource.DataSource {
	return &projectsDataSource{}
}

func (d *projectsDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "virtuacloud_projects"
}

func (d *projectsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all projects for the authenticated Virtua.Cloud account.",
		Attributes: map[string]schema.Attribute{
			"projects": schema.ListNestedAttribute{
				Description: "List of projects.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":                         schema.StringAttribute{Description: "Project UUID.", Computed: true},
						"name":                         schema.StringAttribute{Description: "Project name.", Computed: true},
						"description":                  schema.StringAttribute{Description: "Project description.", Computed: true},
						"cloud_servers_count":          schema.StringAttribute{Description: "Number of cloud servers in project.", Computed: true},
						"monthly_cloud_usage":          schema.StringAttribute{Description: "Monthly cloud usage.", Computed: true},
						"monthly_cloud_usage_estimate": schema.StringAttribute{Description: "Estimated monthly cloud usage.", Computed: true},
						"domains_count":                schema.StringAttribute{Description: "Number of domains in project.", Computed: true},
						"environment":                  schema.StringAttribute{Description: "Project environment (e.g. production, staging).", Computed: true},
						"created_at":                   schema.StringAttribute{Description: "Project creation timestamp.", Computed: true},
					},
				},
			},
		},
	}
}

func (d *projectsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = nil
	if req.ProviderData != nil {
		d.client = req.ProviderData.(*virtuacloudProviderData)
	}
}

func (d *projectsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Provider client not configured")
		return
	}

	result, err := d.client.Client.GetProjects(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read projects", err.Error())
		return
	}

	var state projectsDataSourceModel
	for _, p := range result.Projects {
		env := types.StringNull()
		if p.Environment != nil {
			env = types.StringValue(*p.Environment)
		}
		state.Projects = append(state.Projects, projectModel{
			UUID:                      types.StringValue(p.UUID),
			Name:                      types.StringValue(p.Name),
			Description:               types.StringValue(p.Description),
			CloudServersCount:         types.StringValue(fmt.Sprintf("%d", p.CloudServersCount)),
			MonthlyCloudUsage:         types.StringValue(p.MonthlyCloudUsage),
			MonthlyCloudUsageEstimate: types.StringValue(p.MonthlyCloudUsageEstimate),
			DomainsCount:              types.StringValue(fmt.Sprintf("%d", p.DomainsCount)),
			Environment:               env,
			CreatedAt:                 types.StringValue(p.CreatedAt),
		})
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
