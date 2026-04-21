package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type cloudServerPasswordDataSource struct {
	client *virtuacloudProviderData
}

var _ datasource.DataSource = (*cloudServerPasswordDataSource)(nil)

type cloudServerPasswordDataSourceModel struct {
	ServerUUID    types.String `tfsdk:"server_uuid"`
	PasswordType  types.String `tfsdk:"password_type"`
	Password      types.String `tfsdk:"password"`
}

func NewCloudServerPasswordDataSource() datasource.DataSource {
	return &cloudServerPasswordDataSource{}
}

func (d *cloudServerPasswordDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "virtuacloud_cloud_server_password"
}

func (d *cloudServerPasswordDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve the root or rescue password for a Virtua.Cloud server. The password is sensitive and will not be displayed in logs.",
		Attributes: map[string]schema.Attribute{
			"server_uuid": schema.StringAttribute{
				Description: "UUID of the cloud server.",
				Required:    true,
			},
			"password_type": schema.StringAttribute{
				Description: "Type of password to retrieve: root or rescue.",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "The server password. This value is sensitive.",
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func (d *cloudServerPasswordDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = nil
	if req.ProviderData != nil {
		d.client = req.ProviderData.(*virtuacloudProviderData)
	}
}

func (d *cloudServerPasswordDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Provider client not configured")
		return
	}

	var config cloudServerPasswordDataSourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	passwordType := config.PasswordType.ValueString()
	if passwordType != "root" && passwordType != "rescue" {
		resp.Diagnostics.AddError("Invalid password type", "password_type must be 'root' or 'rescue'")
		return
	}

	result, err := d.client.Client.GetCloudServerPassword(ctx, config.ServerUUID.ValueString(), passwordType)
	if err != nil {
		resp.Diagnostics.AddError("Failed to retrieve server password", err.Error())
		return
	}

	state := cloudServerPasswordDataSourceModel{
		ServerUUID:   config.ServerUUID,
		PasswordType: config.PasswordType,
		Password:     types.StringValue(result.Password),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}