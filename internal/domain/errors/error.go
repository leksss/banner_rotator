package errors

type EventError string

var (
	ErrInvalidRequestSlotAndBannerAreRequired         = EventError("Invalid request. slotID and bannerID are required")
	ErrInvalidRequestSlotAndBannerAndGroupAreRequired = EventError("Invalid request. slotID, bannerID and groupID are required")
	ErrInvalidRequestSlotAndGroupAreRequired          = EventError("Invalid request. slotID and groupID are required")
)

func (ee EventError) Error() string {
	return string(ee)
}
