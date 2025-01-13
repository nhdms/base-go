package enums

// OrderStatus represents the status of an order.proto in the system
type OrderStatus int

const (
	// Draft (-2) - đơn nháp
	// Đơn tạo bất kỳ có một vài thông tin cơ bản.
	Draft OrderStatus = -2

	// New (0) - đơn mới
	// Đơn có đầy đủ các thông tin cơ bản và bắt buộc.
	New OrderStatus = 0

	// AwaitingStock (1) - chờ hàng
	// 🚦Trạng thái tự động. Đơn đã được xác nhận nhưng thiếu hàng trong kho thì hệ thống tự động chuyển sang trạng thái này.
	AwaitingStock OrderStatus = 1

	// Reconfirm (2) - xác nhận lại
	// 🚦Trạng thái tự động. Đơn sau khi được có sản phẩm đủ tồn kho, sẽ tự động chuyển sang trạng thái này
	Reconfirm OrderStatus = 2

	// Confirmed (3) - đã xác nhận
	// Chuyển danh sách đơn xác nhận sang cho kho lấy hàng đóng gói.
	Confirmed OrderStatus = 3

	// Preparing (4) - đang chuẩn bị hàng
	// Đơn được lên danh sách chờ lấy hàng ra khỏi kệ trong kho.
	Preparing OrderStatus = 4

	// HandlingOver (5) - Đang bàn giao vận chuyển
	// Đơn đã chuẩn bị hàng xong, chờ 3Pl đến lấy.
	HandlingOver OrderStatus = 5

	// InTransit (6) - Đang vận chuyển
	// Đơn đã bàn giao cho 3PL.
	InTransit OrderStatus = 6

	// InDelivery (7) - Đang giao
	// Đơn đang được giao cho khách.
	InDelivery OrderStatus = 7

	// Delivered (8) - Đã giao (chờ đối soát)
	// Đơn đã được giao tới khách.
	Delivered OrderStatus = 8

	// DeliveredCompleted (9) - Đã giao (Hoàn tất)
	// Đơn đã xong, Giao đã thanh toán xong cho Khách.
	DeliveredCompleted OrderStatus = 9

	// FailedDelivery (10) - Giao thất bại
	FailedDelivery OrderStatus = 10

	// AwaitingReturn (11) - Chờ hoàn
	// Đơn bị NN từ chối và quá số lần giao/lưu kho.
	AwaitingReturn OrderStatus = 11

	// InReturn (12) - Đang hoàn
	// Đơn đang được 3PL hoàn về.
	InReturn OrderStatus = 12

	// ReturnedStocked (13) - Đã hoàn (Chờ đối soát)
	// Đơn hoàn đã tái nhập kho.
	ReturnedStocked OrderStatus = 13

	// ReturnedCompleted (14) - Đã hoàn (Hoàn tất)
	// Đơn hoàn đã đối soát xong.
	ReturnedCompleted OrderStatus = 14

	// Damaged (15) - Hư hỏng (Chờ xử lý)
	// Đơn hư hỏng vì kho bảo quản không tốt hoặc do 3PL...
	Damaged OrderStatus = 15

	// DamagedCompleted (16) - Hư hỏng, thất lạc (Hoàn tất)
	// Đơn hư hỏng, thất lạc bởi kho hoặc 3PL đã truy thu/ xử lý tài chính xong.
	DamagedCompleted OrderStatus = 16

	// Canceled (17) - Hủy
	// Đơn được hủy trước khi bàn giao cho 3PL.
	Canceled OrderStatus = 17

	// Lost (18) - Thất lạc (Chờ xử lý)
	// Đơn thất lạc vì kho bảo quản không tốt hoặc do 3PL...
	Lost OrderStatus = 18

	// LostCompleted (19) - Thất lạc (Hoàn tất)
	// Đơn thất lạc bởi kho hoặc 3PL đã truy thu/ xử lý tài chính xong.
	LostCompleted OrderStatus = 19
)

var statusNames = []string{
	New:                "New",
	AwaitingStock:      "AwaitingStock",
	Reconfirm:          "Reconfirm",
	Confirmed:          "Confirmed",
	Preparing:          "Preparing",
	HandlingOver:       "HandlingOver",
	InTransit:          "InTransit",
	InDelivery:         "InDelivery",
	Delivered:          "Delivered",
	DeliveredCompleted: "DeliveredCompleted",
	FailedDelivery:     "FailedDelivery",
	AwaitingReturn:     "AwaitingReturn",
	InReturn:           "InReturn",
	ReturnedStocked:    "ReturnedStocked",
	ReturnedCompleted:  "ReturnedCompleted",
	Damaged:            "Damaged",
	DamagedCompleted:   "DamagedCompleted",
	Canceled:           "Canceled",
	Lost:               "Lost",
	LostCompleted:      "LostCompleted",
}

// String returns the string representation of the OrderStatus
func (s OrderStatus) String() string {
	if s < Draft || int(s) >= len(statusNames) {
		return "Unknown"
	}
	if s < 0 {
		return "Draft" // Return "Draft" for -2
	}
	return statusNames[s]
}

// IsValid checks if the status is within valid range
func (s OrderStatus) IsValid() bool {
	return s >= Draft && int(s) < len(statusNames)
}

func (s OrderStatus) EqualNumber(i int) bool {
	return int(s) == i
}
