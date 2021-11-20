package sqlstorage

type ucb1Db struct {
	id       int64 `db:"id"`
	slotID   int64 `db:"slot_id"`
	bannerID int64 `db:"banner_id"`
	groupID  int64 `db:"group_id"`
	hitCnt   int64 `db:"hit_cnt"`
	showCnt  int64 `db:"show_cnt"`
}
