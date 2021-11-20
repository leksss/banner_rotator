package sqlstorage

type ucb1Db struct {
	ID       int64 `db:"id"`
	SlotID   int64 `db:"slot_id"`
	BannerID int64 `db:"banner_id"`
	GroupID  int64 `db:"group_id"`
	HitCnt   int64 `db:"hit_cnt"`
	ShowCnt  int64 `db:"show_cnt"`
}
