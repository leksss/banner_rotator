package errors

import "github.com/pkg/errors"

var (
	// ErrInvalidRequestSlotBannerRequired validation error.
	ErrInvalidRequestSlotBannerRequired = errors.New("Invalid request. slotID and bannerID are required")
	// ErrInvalidRequestSlotBannerGroupRequired validation error.
	ErrInvalidRequestSlotBannerGroupRequired = errors.New("Invalid request. slotID, bannerID and groupID are required")
	// ErrInvalidRequestSlotGroupRequired validation error.
	ErrInvalidRequestSlotGroupRequired = errors.New("Invalid request. slotID and groupID are required")

	// ErrBannerNotFound logic error.
	ErrBannerNotFound = errors.New("Banner not found")
	// ErrNoAvailableBannersInSlot logic error.
	ErrNoAvailableBannersInSlot = errors.New("No available banners in slot")
)
