package enums

type LeadType int

const (
	LeadNew LeadType = iota
	LeadAfterSale
)

// CareState is the custom enum type for care states
type CareState int

const (
	CareState_New CareState = iota
	CareState_UnassignAttempted
	CareState_Assigned
	CareState_NoAttempt
	CareState_Attempted
	CareState_Potential
	CareState_AwaitingStock
	CareState_Reconfirm
	CareState_Confirmed
	CareState_Failed
	CareState_Lost
	CareState_Junk
	CareState_Temp
)

// String method to provide string representation of the CareState enum
func (cs CareState) String() string {
	return [...]string{
		"new",                // Chưa tiếp nhận (Mới)
		"unassign_attempted", // Chưa tiếp nhận (Đã xử lý)
		"assigned",           // Đã tiếp nhận
		"no_attempt",         // Đang xử lý (Chưa xử lý)
		"attempted",          // Đang xử lý (Đã xử lý)
		"potential",          // Đang xử lý (Tiềm năng)
		"awaiting_stock",     // Chốt đơn (Chờ hàng)
		"reconfirm",          // Chốt đơn (Xác nhận lại)
		"confirmed",          // Chốt đơn (Đã chốt)
		"failed",             // Thất bại
		"lost",               // Hủy (Hủy)
		"junk",               // Hủy (Rác)
		"temp",               // Chốt đơn (Chốt tạm)
	}[cs]
}
