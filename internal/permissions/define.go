package permissions

type Permission int

const (
	Dashboard        Permission = iota // 0
	Order                              // 1
	Customer                           // 2
	Product                            // 3
	Supplier                           // 4
	Inventory                          // 5
	Inbound                            // 6
	Outbound                           // 7
	Stocktaking                        // 8
	ReturnHandling                     // 9
	Telesales                          // 10
	CarePage                           // 11
	BotManagement                      // 12
	FanPages                           // 13
	Campaigns                          // 14
	Projects                           // 15
	Countries                          // 16
	Accounts                           // 17
	Roles                              // 18
	Org                                // 19
	Shift                              // 20
	Tags                               // 21
	OrderSource                        // 22
	ReportReason                       // 23
	CancelReason                       // 24
	PrintNotes                         // 25
	Translation                        // 26
	Fulfillment                        // 27
	Marketplace                        // 28
	Currency                           // 29
	DeviceManagement                   // 30
)
