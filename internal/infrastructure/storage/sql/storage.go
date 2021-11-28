package sqlstorage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql" // nolint
	"github.com/jmoiron/sqlx"
	"github.com/leksss/banner_rotator/internal/domain/entities"
	"github.com/leksss/banner_rotator/internal/domain/interfaces"
)

const (
	slotLimit = 20
	statLimit = 100

	hitField  = "hit_cnt"
	showField = "show_cnt"
)

type Storage struct {
	db  *sqlx.DB
	log interfaces.Log
}

func New(db *sqlx.DB, log interfaces.Log) *Storage {
	return &Storage{
		log: log,
		db:  db,
	}
}

func (s *Storage) AddBanner(ctx context.Context, slotID, bannerID uint64) error {
	query := `INSERT INTO slot2banner (slot_id, banner_id) VALUES (:slotID, :bannerID)`
	arg := map[string]interface{}{"slotID": slotID, "bannerID": bannerID}
	_, err := s.execContext(ctx, query, arg)
	return err
}

func (s *Storage) RemoveBanner(ctx context.Context, slotID, bannerID uint64) error {
	query := `DELETE FROM slot2banner WHERE slot_id=:slotID AND banner_id=:bannerID`
	arg := map[string]interface{}{"slotID": slotID, "bannerID": bannerID}
	_, err := s.execContext(ctx, query, arg)
	return err
}

func (s *Storage) IncrementHit(ctx context.Context, slotID, bannerID, groupID uint64) error {
	return s.incrementCounter(ctx, slotID, bannerID, groupID, hitField)
}

func (s *Storage) IncrementShow(ctx context.Context, slotID, bannerID, groupID uint64) error {
	return s.incrementCounter(ctx, slotID, bannerID, groupID, showField)
}

func (s *Storage) incrementCounter(ctx context.Context, slotID, bannerID, groupID uint64, field string) error {
	if slotID == 0 || bannerID == 0 || groupID == 0 {
		s.log.Error("Invalid params: slotID, bannerID, groupID")
		return nil
	}
	params := map[string]interface{}{
		"slotID":   slotID,
		"bannerID": bannerID,
		"groupID":  groupID,
	}
	query := `SELECT id, slot_id, banner_id, group_id, hit_cnt, show_cnt 
				FROM ucb1 
				WHERE slot_id=:slotID 
					AND banner_id=:bannerID 
					AND group_id=:groupID`
	rows, err := s.queryContext(ctx, query, params)
	if err != nil {
		return err
	}
	defer rows.Close()

	var row ucb1Row
	if !rows.Next() {
		query = `INSERT INTO ucb1 (slot_id, banner_id, group_id, hit_cnt, show_cnt) 
					VALUES (:slotID, :bannerID, :groupID, :hitCnt, :showCnt)`
		if field == hitField {
			params["hitCnt"] = 1
			params["showCnt"] = 0
		}
		if field == showField {
			params["hitCnt"] = 0
			params["showCnt"] = 1
		}
	} else {
		if err = rows.StructScan(&row); err != nil {
			return err
		}
		query = `UPDATE ucb1 SET 
						hit_cnt=:hitCnt, 
						show_cnt=:showCnt 
					WHERE slot_id=:slotID 
						AND banner_id=:bannerID 
						AND group_id=:groupID`
		if field == hitField {
			row.HitCnt++
			params["hitCnt"] = row.HitCnt
			params["showCnt"] = row.ShowCnt
		}
		if field == showField {
			row.ShowCnt++
			params["hitCnt"] = row.HitCnt
			params["showCnt"] = row.ShowCnt
		}
	}
	_, err = s.execContext(ctx, query, params)
	return err
}

func (s *Storage) GetSlotCounters(ctx context.Context, slotID, groupID uint64) (entities.BannerCounterMap, error) {
	query := `SELECT id, slot_id, banner_id, group_id, hit_cnt, show_cnt  
				FROM ucb1 
				WHERE slot_id=:slotID AND group_id=:groupID 
				LIMIT :statLimit`
	arg := map[string]interface{}{
		"slotID":    slotID,
		"groupID":   groupID,
		"statLimit": statLimit,
	}
	rows, err := s.queryContext(ctx, query, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var row ucb1Row
	counters := make(entities.BannerCounterMap)
	for rows.Next() {
		err = rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		counters[entities.BannerID(row.BannerID)] = &entities.Counter{
			SlotID:   row.SlotID,
			BannerID: row.BannerID,
			GroupID:  row.GroupID,
			HitCnt:   float64(row.HitCnt),
			ShowCnt:  float64(row.ShowCnt),
		}
	}
	return counters, nil
}

func (s *Storage) GetBannersBySlot(ctx context.Context, slotID uint64) ([]entities.BannerID, error) {
	query := `SELECT banner_id FROM slot2banner WHERE slot_id=:slotID LIMIT :slotLimit`
	arg := map[string]interface{}{
		"slotID":    slotID,
		"slotLimit": slotLimit,
	}
	rows, err := s.queryContext(ctx, query, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bannerID entities.BannerID
	var bannerIDs []entities.BannerID
	for rows.Next() {
		if err = rows.Scan(&bannerID); err != nil {
			return nil, err
		}
		bannerIDs = append(bannerIDs, bannerID)
	}
	return bannerIDs, nil
}

func (s *Storage) execContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	result, err := s.db.NamedExecContext(ctx, query, arg)
	s.logQuery(query, arg)
	return result, err
}

func (s *Storage) queryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	rows, err := s.db.NamedQueryContext(ctx, query, arg)
	s.logQuery(query, arg)
	return rows, err
}

func (s *Storage) logQuery(sql string, arg interface{}) {
	sql = strings.ReplaceAll(sql, "\n", "")
	sql = strings.ReplaceAll(sql, "\t", "")
	byteArg, _ := json.Marshal(arg)
	s.log.Info(fmt.Sprintf("query: %s params: %s", fmt.Sprintf("%q", sql), string(byteArg)))
}
