package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type accountDataSource struct {
	client *virtuacloudProviderData
}

var _ datasource.DataSource = (*accountDataSource)(nil)

type accountDataSourceModel struct {
	Balance                   types.String `tfsdk:"balance"`
	OutstandingBalance        types.String `tfsdk:"outstanding_balance"`
	Currency                  types.String `tfsdk:"currency"`
	MonthlyCloudUsage         types.String `tfsdk:"monthly_cloud_usage"`
	MonthlyCloudUsageEstimate types.String `tfsdk:"monthly_cloud_usage_estimate"`
	Timezone                  types.String `tfsdk:"timezone"`
	TodayCloudUsage           types.String `tfsdk:"today_cloud_usage"`
	CloudServersLimit         types.String `tfsdk:"cloud_servers_limit"`
}

func NewAccountDataSource() datasource.DataSource {
	return &accountDataSource{}
}

func (d *accountDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "virtuacloud_account"
}

func (d *accountDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Retrieve account information for the authenticated Virtua.Cloud user.",
		Attributes: map[string]schema.Attribute{
			"balance": schema.StringAttribute{
				Description: "Current account balance.",
				Computed:    true,
			},
			"outstanding_balance": schema.StringAttribute{
				Description: "Outstanding account balance.",
				Computed:    true,
			},
			"currency": schema.StringAttribute{
				Description: "Account currency (e.g. EUR).",
				Computed:    true,
			},
			"monthly_cloud_usage": schema.StringAttribute{
				Description: "Current monthly cloud usage.",
				Computed:    true,
			},
			"monthly_cloud_usage_estimate": schema.StringAttribute{
				Description: "Estimated monthly cloud usage at end of month.",
				Computed:    true,
			},
			"timezone": schema.StringAttribute{
				Description: "Account timezone.",
				Computed:    true,
			},
			"today_cloud_usage": schema.StringAttribute{
				Description: "Today's cloud usage.",
				Computed:    true,
			},
			"cloud_servers_limit": schema.StringAttribute{
				Description: "Maximum number of cloud servers allowed.",
				Computed:    true,
			},
		},
	}
}

func (d *accountDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = nil
	if req.ProviderData != nil {
		d.client = req.ProviderData.(*virtuacloudProviderData)
	}
}

func (d *accountDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Provider client not configured")
		return
	}

	account, err := d.client.Client.GetAccount(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read account", err.Error())
		return
	}

	var state accountDataSourceModel
	state.Balance = types.StringValue(account.Balance)
	state.OutstandingBalance = types.StringValue(account.OutstandingBalance)
	state.Currency = types.StringValue(account.Currency)
	state.MonthlyCloudUsage = types.StringValue(account.MonthlyCloudUsage)
	state.MonthlyCloudUsageEstimate = types.StringValue(account.MonthlyCloudUsageEstimate)
	state.Timezone = types.StringValue(account.Timezone)
	state.TodayCloudUsage = types.StringValue(account.TodayCloudUsage)
	state.CloudServersLimit = types.StringValue(account.CloudServersLimit)

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
