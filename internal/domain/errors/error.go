package errors

type BannerError string

var (
	ErrInvalidRequestSlotAndBannerAreRequired         = BannerError("Invalid request. slotID and bannerID are required")
	ErrInvalidRequestSlotAndBannerAndGroupAreRequired = BannerError("Invalid request. slotID, bannerID and groupID are required")
	ErrInvalidRequestSlotAndGroupAreRequired          = BannerError("Invalid request. slotID and groupID are required")

	ErrBannerNotFound           = BannerError("Banner not found")
	ErrNoAvailableBannersInSlot = BannerError("No available banners in slot")
)

func (ee BannerError) Error() string {
	return string(ee)
}
