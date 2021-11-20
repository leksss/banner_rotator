package sqlstorage

type ucb1Row struct {
	ID       uint64 `db:"id"`
	SlotID   uint64 `db:"slot_id"`
	BannerID uint64 `db:"banner_id"`
	GroupID  uint64 `db:"group_id"`
	HitCnt   uint64 `db:"hit_cnt"`
	ShowCnt  uint64 `db:"show_cnt"`
}
