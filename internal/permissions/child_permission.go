package permissions

// Dashboard Permissions
const (
	DashboardOrderShipmentOverview = 1 << 0
	DashboardOrderShipmentCarrier  = 1 << 1
	DashboardOrderShipmentSaleReps = 1 << 2
	DashboardOrderShipmentFilter   = 1 << 3
	DashboardCarePageOverview      = 1 << 4
	DashboardCarePagePerformance   = 1 << 5
	DashboardCarePageReasons       = 1 << 6
	DashboardCarePageFilter        = 1 << 7
	DashboardTelesalesOverview     = 1 << 8
	DashboardTelesalesPerformance  = 1 << 9
	DashboardTelesalesReasons      = 1 << 10
	DashboardTelesalesFilter       = 1 << 11
	DashboardMarketingCarePage     = 1 << 12
	DashboardMarketingTelesales    = 1 << 13
)

// Order Permissions
const (
	OrderFetchMany               = 1 << 0
	OrderBasicSearch             = 1 << 1
	OrderAdvanceSearch           = 1 << 2
	OrderDuplicateFilter         = 1 << 3
	OrderExportOrders            = 1 << 4
	OrderCreate                  = 1 << 5
	OrderBulkCreate              = 1 << 6
	OrderFetchOne                = 1 << 7
	OrderEditProductsAndFees     = 1 << 8
	OrderEditGeneralInformation  = 1 << 9
	OrderEditCustomerInformation = 1 << 10
	OrderEditDeliveryInformation = 1 << 11
	OrderEditTags                = 1 << 12
	OrderUpdateStatus            = 1 << 13
	OrderCancel                  = 1 << 14
	OrderReadHistories           = 1 << 15
	OrderBulkUpdateStatus        = 1 << 16
	OrderBulkUpdateSaleReps      = 1 << 17
	OrderBulkUpdateTags          = 1 << 18
	OrderBulkUpdateSource        = 1 << 19
	OrderBulkSyncOrderToFFM      = 1 << 20
)

// Customer Permissions
const (
	CustomerFetchMany = 1 << 0
	CustomerFetchOne  = 1 << 1
	CustomerUpdate    = 1 << 2
	CustomerCreate    = 1 << 3
)

// Product Permissions
const (
	ProductFetchMany   = 1 << 0
	ProductExport      = 1 << 1
	ProductImport      = 1 << 2
	ProductCreate      = 1 << 3
	ProductCreateCombo = 1 << 4
	ProductFetchOne    = 1 << 5
	ProductUpdate      = 1 << 6
)

// Supplier Permissions
const (
	SupplierFetchMany = 1 << 0
	SupplierFetchOne  = 1 << 1
	SupplierCreate    = 1 << 2
	SupplierUpdate    = 1 << 3
)

// Inventory Permissions
const (
	InventoryFetchMany = 1 << 0
)

// Inbound Permissions
const (
	InboundFetchMany    = 1 << 0
	InboundFetchOne     = 1 << 1
	InboundCreate       = 1 << 2
	InboundUpdate       = 1 << 3
	InboundUpdateStatus = 1 << 4
)

// Outbound Permissions
const (
	OutboundFetchMany    = 1 << 0
	OutboundFetchOne     = 1 << 1
	OutboundCreate       = 1 << 2
	OutboundUpdate       = 1 << 3
	OutboundUpdateStatus = 1 << 4
)

// Stocktaking Permissions
const (
	StocktakingFetchMany    = 1 << 0
	StocktakingFetchOne     = 1 << 1
	StocktakingCreate       = 1 << 2
	StocktakingUpdate       = 1 << 3
	StocktakingUpdateStatus = 1 << 4
)

// ReturnHandling Permissions
const (
	ReturnHandlingFetchMany    = 1 << 0
	ReturnHandlingFetchOne     = 1 << 1
	ReturnHandlingCreate       = 1 << 2
	ReturnHandlingUpdate       = 1 << 3
	ReturnHandlingUpdateStatus = 1 << 4
)

// Telesales Permissions
const (
	TelesalesFetchAssignedLeads        = 1 << 0
	TelesalesAssignedLeadsFilter       = 1 << 1
	TelesalesTakeCareLeads             = 1 << 2
	TelesalesAppointments              = 1 << 3
	TelesalesFetchOne                  = 1 << 4
	TelesalesEditProductsAndFees       = 1 << 5
	TelesalesEditGeneralInformation    = 1 << 6
	TelesalesEditCustomerInformation   = 1 << 7
	TelesalesEditDeliveryInformation   = 1 << 8
	TelesalesEditTags                  = 1 << 9
	TelesalesCreateCareReason          = 1 << 10
	TelesalesEditSource                = 1 << 11
	TelesalesActionLogs                = 1 << 12
	TelesalesFetchLeads                = 1 << 13
	TelesalesLeadsFilter               = 1 << 14
	TelesalesManualDistribute          = 1 << 15
	TelesalesManualRevoke              = 1 << 16
	TelesalesExportExcel               = 1 << 17
	TelesalesImportExcel               = 1 << 18
	TelesalesDistributeConfig          = 1 << 19
	TelesalesProcessingProcedureConfig = 1 << 20
)

// CarePage Permissions
const (
	CarePageFetchPageGroups         = 1 << 0
	CarePageCreatePageGroup         = 1 << 1
	CarePageUpdatePageGroup         = 1 << 2
	CarePagePageGroupsAdvanceFilter = 1 << 3
	CarePageManualDistribute        = 1 << 4
	CarePageManualRevoke            = 1 << 5
	CarePageBulkUpdate              = 1 << 6
	CarePageProcess                 = 1 << 7
	CarePageCreateOrder             = 1 << 8
	CarePageCreateAppointment       = 1 << 9
	CarePageActionLogs              = 1 << 10
	CarePageFetchConfigGroups       = 1 << 11
	CarePageCreateConfigGroup       = 1 << 12
	CarePageUpdateConfigGroup       = 1 << 13
	CarePageRemoveConfigGroup       = 1 << 14
	CarePageLimitationSettings      = 1 << 15
	CarePageAIConfigurations        = 1 << 16
	CarePageAIProductConfigurations = 1 << 17
)

// Remaining Modules' Permissions
// BotManagement Permissions
const (
	BotManagementFetchMany     = 1 << 0
	BotManagementCreate        = 1 << 1
	BotManagementUpdate        = 1 << 2
	BotManagementCrawlBotList  = 1 << 3
	BotManagementBotAdsManager = 1 << 4
)

// FanPages Permissions
const (
	FanPagesFetchMany = 1 << 0
	FanPagesLink      = 1 << 1
)

// Campaigns Permissions
const (
	CampaignsFetchMany = 1 << 0
	CampaignsFetchOne  = 1 << 1
	CampaignsCreate    = 1 << 2
)

// Projects Permissions (Module 15)
const (
	ProjectsFetchMany    = 1 << 0
	ProjectsFetchOne     = 1 << 1
	ProjectsCreate       = 1 << 2
	ProjectsUpdate       = 1 << 3
	ProjectsUpdateStatus = 1 << 4
	ProjectsFFMIntegrate = 1 << 5
	ProjectsActionLogs   = 1 << 6
)

// Countries Permissions (Module 16)
const (
	CountriesFetchMany    = 1 << 0
	CountriesCreate       = 1 << 1
	CountriesUpdateStatus = 1 << 2
)

// Accounts Permissions (Module 17)
const (
	AccountsFetchMany    = 1 << 0
	AccountsFetchOne     = 1 << 1
	AccountsCreate       = 1 << 2
	AccountsUpdate       = 1 << 3
	AccountsUpdateStatus = 1 << 4
	AccountsActionLogs   = 1 << 5
)

// Roles Permissions (Module 18)
const (
	RolesFetchMany    = 1 << 0
	RolesFetchOne     = 1 << 1
	RolesCreate       = 1 << 2
	RolesUpdate       = 1 << 3
	RolesUpdateStatus = 1 << 4
	RolesActionLogs   = 1 << 5
)

// Org Permissions (Module 19)
const (
	OrgFetch  = 1 << 0
	OrgUpdate = 1 << 1
)

// Shift Permissions (Module 20)
const (
	ShiftFetchMany = 1 << 0
	ShiftFetchOne  = 1 << 1
	ShiftCreate    = 1 << 2
	ShiftUpdate    = 1 << 3
	ShiftSchedules = 1 << 4
	ShiftAssign    = 1 << 5
)

// Tags Permissions (Module 21)
const (
	TagsFetchMany = 1 << 0
	TagsCreate    = 1 << 1
	TagsUpdate    = 1 << 2
)

// OrderSource Permissions (Module 22)
const (
	OrderSourceFetchMany = 1 << 0
	OrderSourceCreate    = 1 << 1
	OrderSourceUpdate    = 1 << 2
)

// ReportReason Permissions (Module 23)
const (
	ReportReasonFetchMany = 1 << 0
	ReportReasonCreate    = 1 << 1
	ReportReasonUpdate    = 1 << 2
)

// CancelReason Permissions (Module 24)
const (
	CancelReasonFetchMany = 1 << 0
	CancelReasonCreate    = 1 << 1
	CancelReasonUpdate    = 1 << 2
)

// PrintNotes Permissions (Module 25)
const (
	PrintNotesFetchMany = 1 << 0
	PrintNotesCreate    = 1 << 1
	PrintNotesUpdate    = 1 << 2
)

// Translation Permissions (Module 26)
const (
	TranslationFetchMany = 1 << 0
	TranslationUpdate    = 1 << 1
)

// Fulfillment Permissions (Module 27)
const (
	FulfillmentFetchMany    = 1 << 0
	FulfillmentFetchOne     = 1 << 1
	FulfillmentIntegrate    = 1 << 2
	FulfillmentDisintegrate = 1 << 3
)

// Marketplace Permissions (Module 28)
const (
	MarketplaceFetchMany = 1 << 0
	MarketplaceUpdate    = 1 << 1
	MarketplaceCreate    = 1 << 2
)

// Currency Permissions (Module 29)
const (
	CurrencyFetchMany = 1 << 0
	CurrencyUpdate    = 1 << 1
	CurrencyCreate    = 1 << 2
)

// DeviceManagement Permissions (Module 30)
const (
	DeviceManagementFetchMany = 1 << 0
	DeviceManagementCreate    = 1 << 1
)
