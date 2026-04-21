package client

type Account struct {
	Balance                   string `json:"balance"`
	OutstandingBalance        string `json:"outstanding_balance"`
	Currency                  string `json:"currency"`
	MonthlyCloudUsage         string `json:"monthly_cloud_usage"`
	MonthlyCloudUsageEstimate string `json:"monthly_cloud_usage_estimate"`
	Timezone                  string `json:"timezone"`
	TodayCloudUsage           string `json:"today_cloud_usage"`
	CloudServersLimit         string `json:"cloud_servers_limit"`
}

type LimitsResponse struct {
	Success bool        `json:"success"`
	Usage   LimitsUsage `json:"usage"`
	Limits  LimitsInfo  `json:"limits"`
}

type LimitsUsage struct {
	CloudServers int        `json:"cloud_servers"`
	Vcpus        FlexString `json:"vcpus"`
	MemorySize   FlexString `json:"memory_size"`
	RootSpace    FlexString `json:"root_space"`
	IpAddressV4  FlexString `json:"ip_address_v4"`
	IpAddressV6  FlexString `json:"ip_address_v6"`
}

type LimitsInfo struct {
	CloudServers FlexString `json:"cloud_servers"`
	SmtpEnabled  FlexInt    `json:"smtp_enabled"`
}

type Project struct {
	UUID                      string  `json:"uuid"`
	Name                      string  `json:"name"`
	Description               string  `json:"description"`
	CloudServersCount         FlexInt `json:"cloud_servers_count"`
	MonthlyCloudUsage         string  `json:"monthly_cloud_usage"`
	MonthlyCloudUsageEstimate string  `json:"monthly_cloud_usage_estimate"`
	DomainsCount              FlexInt `json:"domains_count"`
	Environment               *string `json:"environment"`
	CreatedAt                 string  `json:"created_at"`
}

type ProjectsResponse struct {
	Success  bool      `json:"success"`
	Projects []Project `json:"projects"`
}

type Offer struct {
	UUID                              string     `json:"uuid"`
	Category                          string     `json:"category"`
	Name                              string     `json:"name"`
	PriceMonth                        FlexString `json:"price_month"`
	PriceHour                         FlexString `json:"price_hour"`
	AdditionalRootSpacePriceMonth     FlexString `json:"additional_root_space_price_month"`
	AdditionalRootSpacePriceHour      FlexString `json:"additional_root_space_price_hour"`
	Vcpus                             FlexString `json:"vcpus"`
	CpuFamily                         string     `json:"cpu_family"`
	MemorySize                        FlexString `json:"memory_size"`
	RootSpace                         FlexString `json:"root_space"`
	RootDiskType                      string     `json:"root_disk_type"`
	Bandwidth                         FlexString `json:"bandwidth"`
	IsWindowsIncluded                 FlexString `json:"is_windows_included"`
	AdditionalIpv4AddressesMax        FlexString `json:"additional_ipv4_addresses_max"`
	AdditionalIpv4AddressesPriceHour  FlexString `json:"additional_ipv4_addresses_price_hour"`
	AdditionalIpv4AddressesPriceMonth FlexString `json:"additional_ipv4_addresses_price_month"`
	CloudZone                         CloudZone  `json:"cloud_zone"`
}

type OffersResponse struct {
	Success bool    `json:"success"`
	Offers  []Offer `json:"offers"`
}

type System struct {
	UUID         string     `json:"uuid"`
	Name         string     `json:"name"`
	Distribution string     `json:"distribution"`
	Version      FlexString `json:"version"`
	Category     string     `json:"category"`
	IsWindows    FlexString `json:"is_windows"`
}

type SystemsResponse struct {
	Success bool     `json:"success"`
	Systems []System `json:"systems"`
}

type CloudZone struct {
	CountryCode    string `json:"country_code"`
	CountryName    string `json:"country_name"`
	CityName       string `json:"city_name"`
	DatacenterName string `json:"datacenter_name"`
	Timezone       string `json:"timezone"`
}

type OfferInfo struct {
	UUID       string     `json:"uuid"`
	Category   string     `json:"category"`
	Name       string     `json:"name"`
	PriceMonth FlexString `json:"price_month"`
	PriceHour  FlexString `json:"price_hour"`
}

type SystemInfo struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type CloudServer struct {
	UUID             string     `json:"uuid"`
	Name             string     `json:"name"`
	Label            *string    `json:"label"`
	Offer            OfferInfo  `json:"offer"`
	StartTime        string     `json:"start_time"`
	EndTime          string     `json:"end_time"`
	MonthlyUsage     FlexString `json:"monthly_usage"`
	ShortDescription *string    `json:"short_description"`
	IsSuspended      FlexString `json:"is_suspended"`
	SuspensionReason *string    `json:"suspension_reason"`
	CloudZone        CloudZone  `json:"cloud_zone"`
	Vcpus            FlexString `json:"vcpus"`
	VcpusUsed        FlexString `json:"vcpus_used"`
	MemorySize       FlexString `json:"memory_size"`
	MemorySizeUsed   FlexString `json:"memory_size_used"`
	RootSpace        FlexString `json:"root_space"`
	RootDiskType     string     `json:"root_disk_type"`
	Status           string     `json:"status"`
	IsError          FlexString `json:"is_error"`
	IsSetup          FlexString `json:"is_setup"`
	System           SystemInfo `json:"system"`
	Uptime           FlexString `json:"uptime"`
	VmType           string     `json:"vm_type"`
	SetupStep        *string    `json:"setup_step"`
	NetbootOn        *string    `json:"netboot_on"`
	NetbootIsSetup   FlexString `json:"netboot_is_setup"`
	IsSmtpAllowed    FlexString `json:"is_smtp_allowed"`
	Hostname         *string    `json:"hostname"`
	IsSetupAt        string     `json:"is_setup_at"`
	IsIpv6Enabled    FlexString `json:"is_ipv6_enabled"`
	IsProcessing     FlexString `json:"is_processing"`
	SshKeys          *string    `json:"ssh_keys"`
	Keyboard         string     `json:"keyboard"`
}

type CloudServerResponse struct {
	Success     bool        `json:"success"`
	UUID        string      `json:"uuid"`
	CloudServer CloudServer `json:"cloud_server"`
}

type CloudServersResponse struct {
	Success      bool          `json:"success"`
	CloudServers []CloudServer `json:"cloud_servers"`
}

type CreateCloudServerRequest struct {
	ProjectUUID string `json:"PROJECT_UUID"`
	Offer       string `json:"offer"`
	System      string `json:"system"`
	Ipv6Enable  int    `json:"ipv6_enable"`
	Hostname    string `json:"hostname,omitempty"`
}

type CreateCloudServerResponse struct {
	Success bool   `json:"success"`
	UUID    string `json:"uuid"`
}

type ResizeCloudServerRequest struct {
	OfferUUID  string `json:"NEW_OFFER_UUID"`
	ResizeDisk bool   `json:"resize_disk"`
}

type CloudServerActionResponse struct {
	Success  bool            `json:"success"`
	UUID     string          `json:"uuid"`
	Messages []ActionMessage `json:"messages"`
}

type ActionMessage struct {
	Content string `json:"content"`
	Type    string `json:"type"`
}

type PasswordResponse struct {
	Success  bool   `json:"success"`
	UUID     string `json:"uuid"`
	Password string `json:"password"`
}

type ApiErrorResponse struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors,omitempty"`
	Message string   `json:"message,omitempty"`
}
