package sqlstorage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // nolint
	"github.com/jmoiron/sqlx"
	"github.com/leksss/banner_rotator/internal/domain/interfaces"
	"github.com/leksss/banner_rotator/internal/infrastructure/logger"
)

const (
	slotLimit = 20
	statLimit = 100

	hitField  = "hit_cnt"
	showField = "show_cnt"
)

type Storage struct {
	db   *sqlx.DB
	conf interfaces.DatabaseConf
	log  logger.Log
}

func New(conf interfaces.DatabaseConf, log logger.Log) *Storage {
	return &Storage{
		conf: conf,
		log:  log,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf("%s:%s@(%s:3306)/%s?parseTime=true", s.conf.User, s.conf.Password, s.conf.Host, s.conf.Name)
	db, err := sqlx.ConnectContext(ctx, "mysql", dsn)
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	if err := s.db.Close(); err != nil {
		return err
	}
	return nil
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

func (s *Storage) IncHit(ctx context.Context, slotID, bannerID, groupID uint64) error {
	return s.incCounter(ctx, slotID, bannerID, groupID, hitField)
}

func (s *Storage) IncShow(ctx context.Context, slotID, bannerID, groupID uint64) error {
	return s.incCounter(ctx, slotID, bannerID, groupID, showField)
}

func (s *Storage) incCounter(ctx context.Context, slotID, bannerID, groupID uint64, field string) error {
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

	var ucb1 ucb1Db
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
		err = rows.StructScan(&ucb1)
		if err != nil {
			return err
		}
		query = `UPDATE ucb1 SET 
						hit_cnt=:hitCnt, 
						show_cnt=:showCnt 
					WHERE slot_id=:slotID 
						AND banner_id=:bannerID 
						AND group_id=:groupID`
		if field == hitField {
			ucb1.HitCnt++
			params["hitCnt"] = ucb1.HitCnt
			params["showCnt"] = ucb1.ShowCnt
		}
		if field == showField {
			ucb1.ShowCnt++
			params["hitCnt"] = ucb1.HitCnt
			params["showCnt"] = ucb1.ShowCnt
		}
	}
	_, err = s.execContext(ctx, query, params)
	return err
}

func (s *Storage) GetSlotCounters(ctx context.Context, slotID, groupID uint64) (map[uint64]uint64, error) {
	query := `SELECT id, banner_id, event_type 
				FROM statistics 
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

	//var bannerIDs []uint64
	//var bannerID uint64
	//for rows.Next() {
	//	err = rows.StructScan(&bannerID)
	//	if err != nil {
	//		return nil, err
	//	}
	//	bannerIDs = append(bannerIDs, bannerID)
	//}
	return map[uint64]uint64{}, nil
}

func (s *Storage) GetBannersBySlot(ctx context.Context, slotID uint64) ([]uint64, error) {
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

	var bannerIDs []uint64
	var bannerID uint64
	for rows.Next() {
		err = rows.StructScan(&bannerID)
		if err != nil {
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
	byteArg, _ := json.MarshalIndent(arg, "", "  ")
	s.log.Info(fmt.Sprintf("%s %s", sql, string(byteArg)))
}
