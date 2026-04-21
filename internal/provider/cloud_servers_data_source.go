package provider

import (
	"context"

	"github.com/Berndinox/tf-provider-virtua-cloud/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type cloudServersDataSource struct {
	client *virtuacloudProviderData
}

var _ datasource.DataSource = (*cloudServersDataSource)(nil)

type cloudServersDataSourceModel struct {
	Servers []cloudServerListItemModel `tfsdk:"servers"`
}

type cloudServerListItemModel struct {
	UUID             types.String `tfsdk:"uuid"`
	Name             types.String `tfsdk:"name"`
	Label            types.String `tfsdk:"label"`
	Status           types.String `tfsdk:"status"`
	Vcpus            types.String `tfsdk:"vcpus"`
	MemorySize       types.String `tfsdk:"memory_size"`
	RootSpace        types.String `tfsdk:"root_space"`
	RootDiskType     types.String `tfsdk:"root_disk_type"`
	IsSetup          types.Bool   `tfsdk:"is_setup"`
	IsError          types.Bool   `tfsdk:"is_error"`
	IsProcessing     types.Bool   `tfsdk:"is_processing"`
	IsSuspended      types.Bool   `tfsdk:"is_suspended"`
	IsSmtpAllowed    types.Bool   `tfsdk:"is_smtp_allowed"`
	IsIpv6Enabled    types.Bool   `tfsdk:"is_ipv6_enabled"`
	StartTime        types.String `tfsdk:"start_time"`
	EndTime          types.String `tfsdk:"end_time"`
	MonthlyUsage     types.String `tfsdk:"monthly_usage"`
	Hostname         types.String `tfsdk:"hostname"`
	VmType           types.String `tfsdk:"vm_type"`
	Keyboard         types.String `tfsdk:"keyboard"`
	Uptime           types.String `tfsdk:"uptime"`
	Offer            types.Object `tfsdk:"offer"`
	CloudZone        types.Object `tfsdk:"cloud_zone"`
	System           types.Object `tfsdk:"system"`
	ProjectUUID      types.String `tfsdk:"project_uuid"`
	OfferUUID        types.String `tfsdk:"offer_uuid"`
	SystemUUID       types.String `tfsdk:"system_uuid"`
	Ipv6Enabled      types.Bool   `tfsdk:"ipv6_enabled"`
	ShortDescription types.String `tfsdk:"short_description"`
	SuspensionReason types.String `tfsdk:"suspension_reason"`
	SetupStep        types.String `tfsdk:"setup_step"`
	NetbootOn        types.String `tfsdk:"netboot_on"`
	NetbootIsSetup   types.String `tfsdk:"netboot_is_setup"`
	IsSetupAt        types.String `tfsdk:"is_setup_at"`
	SshKeys          types.String `tfsdk:"ssh_keys"`
	VcpusUsed        types.String `tfsdk:"vcpus_used"`
	MemorySizeUsed   types.String `tfsdk:"memory_size_used"`
}

func NewCloudServersDataSource() datasource.DataSource {
	return &cloudServersDataSource{}
}

func (d *cloudServersDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "virtuacloud_cloud_servers"
}

func (d *cloudServersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all cloud servers for the authenticated Virtua.Cloud account.",
		Attributes: map[string]schema.Attribute{
			"servers": schema.ListNestedAttribute{
				Description: "List of cloud servers.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":              schema.StringAttribute{Description: "Server UUID.", Computed: true},
						"name":              schema.StringAttribute{Description: "Server name.", Computed: true},
						"label":             schema.StringAttribute{Description: "Server label.", Computed: true},
						"status":            schema.StringAttribute{Description: "Current server status.", Computed: true},
						"vcpus":             schema.StringAttribute{Description: "Number of vCPUs.", Computed: true},
						"vcpus_used":        schema.StringAttribute{Description: "Number of vCPUs currently in use.", Computed: true},
						"memory_size":       schema.StringAttribute{Description: "Memory size in MB.", Computed: true},
						"memory_size_used":  schema.StringAttribute{Description: "Memory currently in use in MB.", Computed: true},
						"root_space":        schema.StringAttribute{Description: "Root disk space in GB.", Computed: true},
						"root_disk_type":    schema.StringAttribute{Description: "Root disk type (ssd or nvme).", Computed: true},
						"is_setup":          schema.BoolAttribute{Description: "Whether the server setup is complete.", Computed: true},
						"is_error":          schema.BoolAttribute{Description: "Whether the server is in an error state.", Computed: true},
						"is_processing":     schema.BoolAttribute{Description: "Whether the server is currently processing an operation.", Computed: true},
						"is_suspended":      schema.BoolAttribute{Description: "Whether the server is suspended.", Computed: true},
						"is_smtp_allowed":   schema.BoolAttribute{Description: "Whether SMTP is allowed.", Computed: true},
						"is_ipv6_enabled":   schema.BoolAttribute{Description: "Whether IPv6 is enabled.", Computed: true},
						"start_time":        schema.StringAttribute{Description: "Server start time.", Computed: true},
						"end_time":          schema.StringAttribute{Description: "Server end time.", Computed: true},
						"monthly_usage":     schema.StringAttribute{Description: "Current monthly usage cost.", Computed: true},
						"hostname":          schema.StringAttribute{Description: "Custom hostname.", Computed: true},
						"vm_type":           schema.StringAttribute{Description: "Virtualization type (e.g. qemu).", Computed: true},
						"keyboard":          schema.StringAttribute{Description: "Keyboard layout.", Computed: true},
						"uptime":            schema.StringAttribute{Description: "Server uptime in seconds.", Computed: true},
						"short_description": schema.StringAttribute{Description: "Short description of the server.", Computed: true},
						"suspension_reason": schema.StringAttribute{Description: "Reason for suspension.", Computed: true},
						"setup_step":        schema.StringAttribute{Description: "Current setup step.", Computed: true},
						"netboot_on":        schema.StringAttribute{Description: "Netboot on status.", Computed: true},
						"netboot_is_setup":  schema.StringAttribute{Description: "Whether netboot is set up.", Computed: true},
						"is_setup_at":       schema.StringAttribute{Description: "Timestamp when setup completed.", Computed: true},
						"ssh_keys":          schema.StringAttribute{Description: "SSH keys assigned to the server.", Computed: true},
						"project_uuid":      schema.StringAttribute{Description: "UUID of the project the server belongs to.", Computed: true},
						"offer_uuid":        schema.StringAttribute{Description: "UUID of the offer.", Computed: true},
						"system_uuid":       schema.StringAttribute{Description: "UUID of the operating system.", Computed: true},
						"ipv6_enabled":      schema.BoolAttribute{Description: "Whether IPv6 is enabled on the server.", Computed: true},
						"offer": schema.SingleNestedAttribute{
							Description: "Offer details.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"uuid":        schema.StringAttribute{Computed: true},
								"category":    schema.StringAttribute{Computed: true},
								"name":        schema.StringAttribute{Computed: true},
								"price_month": schema.StringAttribute{Computed: true},
								"price_hour":  schema.StringAttribute{Computed: true},
							},
						},
						"cloud_zone": schema.SingleNestedAttribute{
							Description: "Cloud zone details.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"country_code":    schema.StringAttribute{Computed: true},
								"country_name":    schema.StringAttribute{Computed: true},
								"city_name":       schema.StringAttribute{Computed: true},
								"datacenter_name": schema.StringAttribute{Computed: true},
								"timezone":        schema.StringAttribute{Computed: true},
							},
						},
						"system": schema.SingleNestedAttribute{
							Description: "Operating system details.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"uuid": schema.StringAttribute{Computed: true},
								"name": schema.StringAttribute{Computed: true},
							},
						},
					},
				},
			},
		},
	}
}

func (d *cloudServersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = nil
	if req.ProviderData != nil {
		d.client = req.ProviderData.(*virtuacloudProviderData)
	}
}

func (d *cloudServersDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Provider client not configured")
		return
	}

	result, err := d.client.Client.ListCloudServers(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list cloud servers", err.Error())
		return
	}

	var state cloudServersDataSourceModel
	for _, s := range result.CloudServers {
		item := cloudServerListItemModel{
			UUID:           types.StringValue(s.UUID),
			Name:           types.StringValue(s.Name),
			Status:         types.StringValue(s.Status),
			Vcpus:          types.StringValue(string(s.Vcpus)),
			VcpusUsed:      types.StringValue(string(s.VcpusUsed)),
			MemorySize:     types.StringValue(string(s.MemorySize)),
			MemorySizeUsed: types.StringValue(string(s.MemorySizeUsed)),
			RootSpace:      types.StringValue(string(s.RootSpace)),
			RootDiskType:   types.StringValue(s.RootDiskType),
			IsSetup:        types.BoolValue(string(s.IsSetup) == "1"),
			IsError:        types.BoolValue(string(s.IsError) == "1"),
			IsProcessing:   types.BoolValue(string(s.IsProcessing) == "1"),
			IsSuspended:    types.BoolValue(string(s.IsSuspended) == "1"),
			IsSmtpAllowed:  types.BoolValue(string(s.IsSmtpAllowed) == "1"),
			IsIpv6Enabled:  types.BoolValue(string(s.IsIpv6Enabled) == "1"),
			StartTime:      types.StringValue(s.StartTime),
			EndTime:        types.StringValue(s.EndTime),
			MonthlyUsage:   types.StringValue(string(s.MonthlyUsage)),
			VmType:         types.StringValue(s.VmType),
			Keyboard:       types.StringValue(s.Keyboard),
			Uptime:         types.StringValue(string(s.Uptime)),
			NetbootIsSetup: types.StringValue(string(s.NetbootIsSetup)),
			IsSetupAt:      types.StringValue(s.IsSetupAt),
		}

		if s.Label != nil {
			item.Label = types.StringValue(*s.Label)
		} else {
			item.Label = types.StringNull()
		}

		if s.Hostname != nil {
			item.Hostname = types.StringValue(*s.Hostname)
		} else {
			item.Hostname = types.StringNull()
		}

		if s.SshKeys != nil {
			item.SshKeys = types.StringValue(*s.SshKeys)
		} else {
			item.SshKeys = types.StringNull()
		}

		if s.ShortDescription != nil {
			item.ShortDescription = types.StringValue(*s.ShortDescription)
		} else {
			item.ShortDescription = types.StringNull()
		}

		if s.SuspensionReason != nil {
			item.SuspensionReason = types.StringValue(*s.SuspensionReason)
		} else {
			item.SuspensionReason = types.StringNull()
		}

		if s.SetupStep != nil {
			item.SetupStep = types.StringValue(*s.SetupStep)
		} else {
			item.SetupStep = types.StringNull()
		}

		if s.NetbootOn != nil {
			item.NetbootOn = types.StringValue(*s.NetbootOn)
		} else {
			item.NetbootOn = types.StringNull()
		}

		offerObj := offerInfoModel{
			UUID:       types.StringValue(s.Offer.UUID),
			Category:   types.StringValue(s.Offer.Category),
			Name:       types.StringValue(s.Offer.Name),
			PriceMonth: types.StringValue(string(s.Offer.PriceMonth)),
			PriceHour:  types.StringValue(string(s.Offer.PriceHour)),
		}
		offerObjValue, _ := types.ObjectValueFrom(ctx, item.Offer.AttributeTypes(ctx), offerObj)
		item.Offer = offerObjValue

		zoneObj := cloudZoneModel{
			CountryCode:    types.StringValue(s.CloudZone.CountryCode),
			CountryName:    types.StringValue(s.CloudZone.CountryName),
			CityName:       types.StringValue(s.CloudZone.CityName),
			DatacenterName: types.StringValue(s.CloudZone.DatacenterName),
			Timezone:       types.StringValue(s.CloudZone.Timezone),
		}
		zoneObjValue, _ := types.ObjectValueFrom(ctx, item.CloudZone.AttributeTypes(ctx), zoneObj)
		item.CloudZone = zoneObjValue

		sysObj := systemInfoModel{
			UUID: types.StringValue(s.System.UUID),
			Name: types.StringValue(s.System.Name),
		}
		sysObjValue, _ := types.ObjectValueFrom(ctx, item.System.AttributeTypes(ctx), sysObj)
		item.System = sysObjValue

		state.Servers = append(state.Servers, item)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// cloudServerListItemToModel is used by cloud_servers data source
func cloudServerToListItem(ctx context.Context, s client.CloudServer) cloudServerListItemModel {
	return cloudServerListItemModel{}
}
