package provider

import (
	"context"

	"github.com/Berndinox/tf-provider-virtua-cloud/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type cloudServerResource struct {
	client *virtuacloudProviderData
}

type cloudServerResourceModel struct {
	UUID               types.String `tfsdk:"uuid"`
	ProjectUUID        types.String `tfsdk:"project_uuid"`
	OfferUUID          types.String `tfsdk:"offer_uuid"`
	SystemUUID         types.String `tfsdk:"system_uuid"`
	Ipv6Enabled        types.Bool   `tfsdk:"ipv6_enabled"`
	Hostname           types.String `tfsdk:"hostname"`
	PowerState         types.String `tfsdk:"power_state"`
	ResizeDisk         types.Bool   `tfsdk:"resize_disk"`
	RestartTriggeredAt types.String `tfsdk:"restart_triggered_at"`
	Name               types.String `tfsdk:"name"`
	Label              types.String `tfsdk:"label"`
	Status             types.String `tfsdk:"status"`
	Vcpus              types.String `tfsdk:"vcpus"`
	VcpusUsed          types.String `tfsdk:"vcpus_used"`
	MemorySize         types.String `tfsdk:"memory_size"`
	MemorySizeUsed     types.String `tfsdk:"memory_size_used"`
	RootSpace          types.String `tfsdk:"root_space"`
	RootDiskType       types.String `tfsdk:"root_disk_type"`
	Uptime             types.String `tfsdk:"uptime"`
	VmType             types.String `tfsdk:"vm_type"`
	Keyboard           types.String `tfsdk:"keyboard"`
	IsSetup            types.Bool   `tfsdk:"is_setup"`
	IsError            types.Bool   `tfsdk:"is_error"`
	IsProcessing       types.Bool   `tfsdk:"is_processing"`
	IsSuspended        types.Bool   `tfsdk:"is_suspended"`
	IsSmtpAllowed      types.Bool   `tfsdk:"is_smtp_allowed"`
	IsIpv6Enabled      types.Bool   `tfsdk:"is_ipv6_enabled_read"`
	StartTime          types.String `tfsdk:"start_time"`
	EndTime            types.String `tfsdk:"end_time"`
	MonthlyUsage       types.String `tfsdk:"monthly_usage"`
	SshKeys            types.String `tfsdk:"ssh_keys"`
	ShortDescription   types.String `tfsdk:"short_description"`
	SuspensionReason   types.String `tfsdk:"suspension_reason"`
	SetupStep          types.String `tfsdk:"setup_step"`
	NetbootOn          types.String `tfsdk:"netboot_on"`
	NetbootIsSetup     types.String `tfsdk:"netboot_is_setup"`
	IsSetupAt          types.String `tfsdk:"is_setup_at"`
	Offer              types.Object `tfsdk:"offer"`
	CloudZone          types.Object `tfsdk:"cloud_zone"`
	System             types.Object `tfsdk:"system"`
}

type offerInfoModel struct {
	UUID       types.String `tfsdk:"uuid"`
	Category   types.String `tfsdk:"category"`
	Name       types.String `tfsdk:"name"`
	PriceMonth types.String `tfsdk:"price_month"`
	PriceHour  types.String `tfsdk:"price_hour"`
}

type systemInfoModel struct {
	UUID types.String `tfsdk:"uuid"`
	Name types.String `tfsdk:"name"`
}

func NewCloudServerResource() resource.Resource {
	return &cloudServerResource{}
}

func (r *cloudServerResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "virtuacloud_cloud_server"
}

func (r *cloudServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Virtua.Cloud cloud server. Server creation may take 1-2 minutes.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Description: "Server UUID. This is the unique identifier for the server.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_uuid": schema.StringAttribute{
				Description: "UUID of the project to assign the server to. Changing this forces a new resource.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"offer_uuid": schema.StringAttribute{
				Description: "UUID of the offer to use. Changing this triggers a resize operation.",
				Required:    true,
			},
			"system_uuid": schema.StringAttribute{
				Description: "UUID of the operating system to install. Changing this forces a new resource.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"ipv6_enabled": schema.BoolAttribute{
				Description: "Whether to enable IPv6. Changing this forces a new resource.",
				Required:    true,
				PlanModifiers: []planmodifier.Bool{
					boolRequiresReplace(),
				},
			},
			"hostname": schema.StringAttribute{
				Description: "Custom hostname for the server.",
				Optional:    true,
			},
			"power_state": schema.StringAttribute{
				Description: "Desired power state of the server: running or stopped. Default is running.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("running"),
				Validators: []validator.String{
					stringvalidator.OneOf("running", "stopped"),
				},
			},
			"resize_disk": schema.BoolAttribute{
				Description: "Whether to resize the disk when changing the offer. Default is false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"restart_triggered_at": schema.StringAttribute{
				Description: "Trigger a server restart by updating this to the current timestamp. The server will restart when this value changes.",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Auto-generated server name.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"label": schema.StringAttribute{
				Description: "Server label.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "Current server status.",
				Computed:    true,
			},
			"vcpus": schema.StringAttribute{
				Description: "Number of vCPUs.",
				Computed:    true,
			},
			"vcpus_used": schema.StringAttribute{
				Description: "Number of vCPUs currently in use.",
				Computed:    true,
			},
			"memory_size": schema.StringAttribute{
				Description: "Memory size in MB.",
				Computed:    true,
			},
			"memory_size_used": schema.StringAttribute{
				Description: "Memory currently in use in MB.",
				Computed:    true,
			},
			"root_space": schema.StringAttribute{
				Description: "Root disk space in GB.",
				Computed:    true,
			},
			"root_disk_type": schema.StringAttribute{
				Description: "Root disk type (ssd or nvme).",
				Computed:    true,
			},
			"uptime": schema.StringAttribute{
				Description: "Server uptime in seconds.",
				Computed:    true,
			},
			"vm_type": schema.StringAttribute{
				Description: "Virtualization type (e.g. qemu).",
				Computed:    true,
			},
			"keyboard": schema.StringAttribute{
				Description: "Keyboard layout.",
				Computed:    true,
			},
			"is_setup": schema.BoolAttribute{
				Description: "Whether the server setup is complete.",
				Computed:    true,
			},
			"is_error": schema.BoolAttribute{
				Description: "Whether the server is in an error state.",
				Computed:    true,
			},
			"is_processing": schema.BoolAttribute{
				Description: "Whether the server is currently processing an operation.",
				Computed:    true,
			},
			"is_suspended": schema.BoolAttribute{
				Description: "Whether the server is suspended.",
				Computed:    true,
			},
			"is_smtp_allowed": schema.BoolAttribute{
				Description: "Whether SMTP is allowed.",
				Computed:    true,
			},
			"is_ipv6_enabled_read": schema.BoolAttribute{
				Description: "Whether IPv6 is enabled on the server (read-only, reflects actual server state).",
				Computed:    true,
			},
			"start_time": schema.StringAttribute{
				Description: "Server start time.",
				Computed:    true,
			},
			"end_time": schema.StringAttribute{
				Description: "Server end time.",
				Computed:    true,
			},
			"monthly_usage": schema.StringAttribute{
				Description: "Current monthly usage cost.",
				Computed:    true,
			},
			"ssh_keys": schema.StringAttribute{
				Description: "SSH keys assigned to the server.",
				Computed:    true,
			},
			"short_description": schema.StringAttribute{
				Description: "Short description of the server.",
				Computed:    true,
			},
			"suspension_reason": schema.StringAttribute{
				Description: "Reason for suspension, if suspended.",
				Computed:    true,
			},
			"setup_step": schema.StringAttribute{
				Description: "Current setup step of the server.",
				Computed:    true,
			},
			"netboot_on": schema.StringAttribute{
				Description: "Netboot on status.",
				Computed:    true,
			},
			"netboot_is_setup": schema.StringAttribute{
				Description: "Whether netboot is set up.",
				Computed:    true,
			},
			"is_setup_at": schema.StringAttribute{
				Description: "Timestamp when setup completed.",
				Computed:    true,
			},
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
	}
}

func (r *cloudServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = nil
	if req.ProviderData != nil {
		r.client = req.ProviderData.(*virtuacloudProviderData)
	}
}

func (r *cloudServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Provider client not configured")
		return
	}

	var plan cloudServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ipv6Enable := 0
	if plan.Ipv6Enabled.ValueBool() {
		ipv6Enable = 1
	}

	createReq := client.CreateCloudServerRequest{
		ProjectUUID: plan.ProjectUUID.ValueString(),
		Offer:       plan.OfferUUID.ValueString(),
		System:      plan.SystemUUID.ValueString(),
		Ipv6Enable:  ipv6Enable,
		Hostname:    plan.Hostname.ValueString(),
	}

	result, err := r.client.Client.CreateCloudServer(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create cloud server", err.Error())
		return
	}

	plan.UUID = types.StringValue(result.UUID)

	err = r.client.Client.WaitForCloudServerSetup(ctx, result.UUID, serverCreateTimeout)
	if err != nil {
		resp.Diagnostics.AddWarning("Server created but setup not yet complete", err.Error())
	}

	if plan.PowerState.ValueString() == "running" {
		err = r.client.Client.WaitForCloudServerStatus(ctx, result.UUID, "running", serverPowerTimeout)
		if err != nil {
			resp.Diagnostics.AddWarning("Server created but not yet running", err.Error())
		}
	} else if plan.PowerState.ValueString() == "stopped" {
		_, stopErr := r.client.Client.StopCloudServer(ctx, result.UUID)
		if stopErr != nil {
			resp.Diagnostics.AddWarning("Server created but failed to stop", stopErr.Error())
		} else {
			waitErr := r.client.Client.WaitForCloudServerStatus(ctx, result.UUID, "stopped", serverPowerTimeout)
			if waitErr != nil {
				resp.Diagnostics.AddWarning("Server stop initiated but not yet confirmed", waitErr.Error())
			}
		}
	}

	server, err := r.client.Client.GetCloudServer(ctx, result.UUID)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read server after creation", err.Error())
		return
	}

	r.populateModel(ctx, server, &plan)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *cloudServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Provider client not configured")
		return
	}

	var state cloudServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	server, err := r.client.Client.GetCloudServer(ctx, state.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read cloud server", err.Error())
		return
	}

	if server.UUID == "" {
		resp.State.RemoveResource(ctx)
		return
	}

	r.populateModel(ctx, server, &state)
	state.PowerState = mapStatusToPowerState(server.Status)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *cloudServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Provider client not configured")
		return
	}

	var plan, state cloudServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	uuid := plan.UUID.ValueString()

	if plan.OfferUUID.ValueString() != state.OfferUUID.ValueString() {
		_, err := r.client.Client.ResizeCloudServer(ctx, uuid, client.ResizeCloudServerRequest{
			OfferUUID:  plan.OfferUUID.ValueString(),
			ResizeDisk: plan.ResizeDisk.ValueBool(),
		})
		if err != nil {
			resp.Diagnostics.AddError("Failed to resize cloud server", err.Error())
			return
		}
		err = r.client.Client.WaitForCloudServerStatus(ctx, uuid, "running", serverResizeTimeout)
		if err != nil {
			resp.Diagnostics.AddWarning("Resize initiated but server not yet running", err.Error())
		}
	}

	if plan.PowerState.ValueString() != state.PowerState.ValueString() {
		switch plan.PowerState.ValueString() {
		case "running":
			_, err := r.client.Client.StartCloudServer(ctx, uuid)
			if err != nil {
				resp.Diagnostics.AddError("Failed to start cloud server", err.Error())
				return
			}
			err = r.client.Client.WaitForCloudServerStatus(ctx, uuid, "running", serverPowerTimeout)
			if err != nil {
				resp.Diagnostics.AddWarning("Start initiated but server not yet running", err.Error())
			}
		case "stopped":
			_, err := r.client.Client.StopCloudServer(ctx, uuid)
			if err != nil {
				resp.Diagnostics.AddError("Failed to stop cloud server", err.Error())
				return
			}
			err = r.client.Client.WaitForCloudServerStatus(ctx, uuid, "stopped", serverPowerTimeout)
			if err != nil {
				resp.Diagnostics.AddWarning("Stop initiated but server not yet stopped", err.Error())
			}
		}
	}

	if plan.RestartTriggeredAt.ValueString() != state.RestartTriggeredAt.ValueString() && !plan.RestartTriggeredAt.IsNull() {
		_, err := r.client.Client.RestartCloudServer(ctx, uuid)
		if err != nil {
			resp.Diagnostics.AddError("Failed to restart cloud server", err.Error())
			return
		}
		err = r.client.Client.WaitForCloudServerStatus(ctx, uuid, "running", serverRestartTimeout)
		if err != nil {
			resp.Diagnostics.AddWarning("Restart initiated but server not yet running", err.Error())
		}
	}

	server, err := r.client.Client.GetCloudServer(ctx, uuid)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read server after update", err.Error())
		return
	}

	r.populateModel(ctx, server, &plan)
	plan.PowerState = mapStatusToPowerState(server.Status)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *cloudServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Provider client not configured")
		return
	}

	var state cloudServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.Client.DeleteCloudServer(ctx, state.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to delete cloud server", err.Error())
		return
	}
}

func (r *cloudServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

func (r *cloudServerResource) populateModel(ctx context.Context, server *client.CloudServer, model *cloudServerResourceModel) {
	model.UUID = types.StringValue(server.UUID)
	model.Name = types.StringValue(server.Name)
	model.Status = types.StringValue(server.Status)
	model.Vcpus = types.StringValue(string(server.Vcpus))
	model.VcpusUsed = types.StringValue(string(server.VcpusUsed))
	model.MemorySize = types.StringValue(string(server.MemorySize))
	model.MemorySizeUsed = types.StringValue(string(server.MemorySizeUsed))
	model.RootSpace = types.StringValue(string(server.RootSpace))
	model.RootDiskType = types.StringValue(server.RootDiskType)
	model.Uptime = types.StringValue(string(server.Uptime))
	model.VmType = types.StringValue(server.VmType)
	model.Keyboard = types.StringValue(server.Keyboard)
	model.StartTime = types.StringValue(server.StartTime)
	model.EndTime = types.StringValue(server.EndTime)
	model.MonthlyUsage = types.StringValue(string(server.MonthlyUsage))
	model.IsSetup = types.BoolValue(string(server.IsSetup) == "1")
	model.IsError = types.BoolValue(string(server.IsError) == "1")
	model.IsProcessing = types.BoolValue(string(server.IsProcessing) == "1")
	model.IsSuspended = types.BoolValue(string(server.IsSuspended) == "1")
	model.IsSmtpAllowed = types.BoolValue(string(server.IsSmtpAllowed) == "1")
	model.IsIpv6Enabled = types.BoolValue(string(server.IsIpv6Enabled) == "1")

	if server.Label != nil {
		model.Label = types.StringValue(*server.Label)
	} else {
		model.Label = types.StringNull()
	}

	if server.SshKeys != nil {
		model.SshKeys = types.StringValue(*server.SshKeys)
	} else {
		model.SshKeys = types.StringNull()
	}

	if server.ShortDescription != nil {
		model.ShortDescription = types.StringValue(*server.ShortDescription)
	} else {
		model.ShortDescription = types.StringNull()
	}

	if server.SuspensionReason != nil {
		model.SuspensionReason = types.StringValue(*server.SuspensionReason)
	} else {
		model.SuspensionReason = types.StringNull()
	}

	if server.SetupStep != nil {
		model.SetupStep = types.StringValue(*server.SetupStep)
	} else {
		model.SetupStep = types.StringNull()
	}

	if server.NetbootOn != nil {
		model.NetbootOn = types.StringValue(*server.NetbootOn)
	} else {
		model.NetbootOn = types.StringNull()
	}

	model.NetbootIsSetup = types.StringValue(string(server.NetbootIsSetup))
	model.IsSetupAt = types.StringValue(server.IsSetupAt)

	offerObj := offerInfoModel{
		UUID:       types.StringValue(server.Offer.UUID),
		Category:   types.StringValue(server.Offer.Category),
		Name:       types.StringValue(server.Offer.Name),
		PriceMonth: types.StringValue(string(server.Offer.PriceMonth)),
		PriceHour:  types.StringValue(string(server.Offer.PriceHour)),
	}
	offerObjValue, _ := types.ObjectValueFrom(ctx, model.Offer.AttributeTypes(ctx), offerObj)
	model.Offer = offerObjValue

	zoneObj := cloudZoneModel{
		CountryCode:    types.StringValue(server.CloudZone.CountryCode),
		CountryName:    types.StringValue(server.CloudZone.CountryName),
		CityName:       types.StringValue(server.CloudZone.CityName),
		DatacenterName: types.StringValue(server.CloudZone.DatacenterName),
		Timezone:       types.StringValue(server.CloudZone.Timezone),
	}
	zoneObjValue, _ := types.ObjectValueFrom(ctx, model.CloudZone.AttributeTypes(ctx), zoneObj)
	model.CloudZone = zoneObjValue

	sysObj := systemInfoModel{
		UUID: types.StringValue(server.System.UUID),
		Name: types.StringValue(server.System.Name),
	}
	sysObjValue, _ := types.ObjectValueFrom(ctx, model.System.AttributeTypes(ctx), sysObj)
	model.System = sysObjValue
}

func mapStatusToPowerState(status string) types.String {
	switch status {
	case "running":
		return types.StringValue("running")
	case "stopped":
		return types.StringValue("stopped")
	default:
		return types.StringValue(status)
	}
}

func boolRequiresReplace() planmodifier.Bool {
	return boolRequiresReplaceModifier{}
}

type boolRequiresReplaceModifier struct{}

func (m boolRequiresReplaceModifier) Description(_ context.Context) string {
	return "Changing this attribute forces a resource replacement."
}

func (m boolRequiresReplaceModifier) MarkdownDescription(_ context.Context) string {
	return "Changing this attribute forces a resource replacement."
}

func (m boolRequiresReplaceModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.StateValue.IsNull() || req.PlanValue.IsNull() {
		return
	}
	if req.StateValue.ValueBool() != req.PlanValue.ValueBool() {
		resp.RequiresReplace = true
	}
}

func int64RequiresReplace() planmodifier.Int64 {
	return int64RequiresReplaceModifier{}
}

type int64RequiresReplaceModifier struct{}

func (m int64RequiresReplaceModifier) Description(_ context.Context) string {
	return "Changing this attribute forces a resource replacement."
}

func (m int64RequiresReplaceModifier) MarkdownDescription(_ context.Context) string {
	return "Changing this attribute forces a resource replacement."
}

func (m int64RequiresReplaceModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if req.StateValue.IsNull() || req.PlanValue.IsNull() {
		return
	}
	if req.StateValue.ValueInt64() != req.PlanValue.ValueInt64() {
		resp.RequiresReplace = true
	}
}

var _ resource.Resource = (*cloudServerResource)(nil)
var _ resource.ResourceWithImportState = (*cloudServerResource)(nil)
var _ planmodifier.Bool = boolRequiresReplaceModifier{}
var _ planmodifier.Int64 = int64RequiresReplaceModifier{}
