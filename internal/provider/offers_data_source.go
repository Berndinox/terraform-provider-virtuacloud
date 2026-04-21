package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type offersDataSource struct {
	client *virtuacloudProviderData
}

var _ datasource.DataSource = (*offersDataSource)(nil)

type offersDataSourceModel struct {
	Offers []offerModel `tfsdk:"offers"`
}

type offerModel struct {
	UUID                              types.String   `tfsdk:"uuid"`
	Category                          types.String   `tfsdk:"category"`
	Name                              types.String   `tfsdk:"name"`
	PriceMonth                        types.String   `tfsdk:"price_month"`
	PriceHour                         types.String   `tfsdk:"price_hour"`
	AdditionalRootSpacePriceMonth     types.String   `tfsdk:"additional_root_space_price_month"`
	AdditionalRootSpacePriceHour      types.String   `tfsdk:"additional_root_space_price_hour"`
	Vcpus                             types.String   `tfsdk:"vcpus"`
	CpuFamily                         types.String   `tfsdk:"cpu_family"`
	MemorySize                        types.String   `tfsdk:"memory_size"`
	RootSpace                         types.String   `tfsdk:"root_space"`
	RootDiskType                      types.String   `tfsdk:"root_disk_type"`
	Bandwidth                         types.String   `tfsdk:"bandwidth"`
	IsWindowsIncluded                 types.String   `tfsdk:"is_windows_included"`
	AdditionalIpv4AddressesMax        types.String   `tfsdk:"additional_ipv4_addresses_max"`
	AdditionalIpv4AddressesPriceHour  types.String   `tfsdk:"additional_ipv4_addresses_price_hour"`
	AdditionalIpv4AddressesPriceMonth types.String   `tfsdk:"additional_ipv4_addresses_price_month"`
	CloudZone                         cloudZoneModel `tfsdk:"cloud_zone"`
}

type cloudZoneModel struct {
	CountryCode    types.String `tfsdk:"country_code"`
	CountryName    types.String `tfsdk:"country_name"`
	CityName       types.String `tfsdk:"city_name"`
	DatacenterName types.String `tfsdk:"datacenter_name"`
	Timezone       types.String `tfsdk:"timezone"`
}

func NewOffersDataSource() datasource.DataSource {
	return &offersDataSource{}
}

func (d *offersDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "virtuacloud_offers"
}

func (d *offersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all available cloud server offers from Virtua.Cloud.",
		Attributes: map[string]schema.Attribute{
			"offers": schema.ListNestedAttribute{
				Description: "List of available offers.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid":                                  schema.StringAttribute{Description: "Offer UUID.", Computed: true},
						"category":                              schema.StringAttribute{Description: "Offer category (e.g. vcs, vcs-win).", Computed: true},
						"name":                                  schema.StringAttribute{Description: "Offer name.", Computed: true},
						"price_month":                           schema.StringAttribute{Description: "Monthly price.", Computed: true},
						"price_hour":                            schema.StringAttribute{Description: "Hourly price.", Computed: true},
						"additional_root_space_price_month":     schema.StringAttribute{Description: "Monthly price per additional GB of root space.", Computed: true},
						"additional_root_space_price_hour":      schema.StringAttribute{Description: "Hourly price per additional GB of root space.", Computed: true},
						"vcpus":                                 schema.StringAttribute{Description: "Number of vCPUs.", Computed: true},
						"cpu_family":                            schema.StringAttribute{Description: "CPU family (e.g. intel-xeon, amd-epyc, amd-ryzen).", Computed: true},
						"memory_size":                           schema.StringAttribute{Description: "Memory size in MB.", Computed: true},
						"root_space":                            schema.StringAttribute{Description: "Root disk space in GB.", Computed: true},
						"root_disk_type":                        schema.StringAttribute{Description: "Root disk type (ssd or nvme).", Computed: true},
						"bandwidth":                             schema.StringAttribute{Description: "Bandwidth in Mbps.", Computed: true},
						"is_windows_included":                   schema.StringAttribute{Description: "Whether Windows license is included (1=yes, 0=no).", Computed: true},
						"additional_ipv4_addresses_max":         schema.StringAttribute{Description: "Maximum additional IPv4 addresses.", Computed: true},
						"additional_ipv4_addresses_price_hour":  schema.StringAttribute{Description: "Hourly price per additional IPv4 address.", Computed: true},
						"additional_ipv4_addresses_price_month": schema.StringAttribute{Description: "Monthly price per additional IPv4 address.", Computed: true},
						"cloud_zone": schema.SingleNestedAttribute{
							Description: "Cloud zone details for this offer.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"country_code":    schema.StringAttribute{Description: "Country code.", Computed: true},
								"country_name":    schema.StringAttribute{Description: "Country name.", Computed: true},
								"city_name":       schema.StringAttribute{Description: "City name.", Computed: true},
								"datacenter_name": schema.StringAttribute{Description: "Datacenter name.", Computed: true},
								"timezone":        schema.StringAttribute{Description: "Timezone.", Computed: true},
							},
						},
					},
				},
			},
		},
	}
}

func (d *offersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = nil
	if req.ProviderData != nil {
		d.client = req.ProviderData.(*virtuacloudProviderData)
	}
}

func (d *offersDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Provider not configured", "Provider client not configured")
		return
	}

	result, err := d.client.Client.GetOffers(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read offers", err.Error())
		return
	}

	var state offersDataSourceModel
	for _, o := range result.Offers {
		state.Offers = append(state.Offers, offerModel{
			UUID:                              types.StringValue(o.UUID),
			Category:                          types.StringValue(o.Category),
			Name:                              types.StringValue(o.Name),
			PriceMonth:                        types.StringValue(string(o.PriceMonth)),
			PriceHour:                         types.StringValue(string(o.PriceHour)),
			AdditionalRootSpacePriceMonth:     types.StringValue(string(o.AdditionalRootSpacePriceMonth)),
			AdditionalRootSpacePriceHour:      types.StringValue(string(o.AdditionalRootSpacePriceHour)),
			Vcpus:                             types.StringValue(string(o.Vcpus)),
			CpuFamily:                         types.StringValue(o.CpuFamily),
			MemorySize:                        types.StringValue(string(o.MemorySize)),
			RootSpace:                         types.StringValue(string(o.RootSpace)),
			RootDiskType:                      types.StringValue(o.RootDiskType),
			Bandwidth:                         types.StringValue(string(o.Bandwidth)),
			IsWindowsIncluded:                 types.StringValue(string(o.IsWindowsIncluded)),
			AdditionalIpv4AddressesMax:        types.StringValue(string(o.AdditionalIpv4AddressesMax)),
			AdditionalIpv4AddressesPriceHour:  types.StringValue(string(o.AdditionalIpv4AddressesPriceHour)),
			AdditionalIpv4AddressesPriceMonth: types.StringValue(string(o.AdditionalIpv4AddressesPriceMonth)),
			CloudZone: cloudZoneModel{
				CountryCode:    types.StringValue(o.CloudZone.CountryCode),
				CountryName:    types.StringValue(o.CloudZone.CountryName),
				CityName:       types.StringValue(o.CloudZone.CityName),
				DatacenterName: types.StringValue(o.CloudZone.DatacenterName),
				Timezone:       types.StringValue(o.CloudZone.Timezone),
			},
		})
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
