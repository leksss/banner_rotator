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

const slotLimit = 20

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

func (s *Storage) HitBanner(ctx context.Context, slotID, bannerID, groupID uint64) error {
	// отправляем событие в очередь, увеличиваем счетчик переходов
	return nil
}

func (s *Storage) GetBanner(ctx context.Context, slotID, groupID uint64) (uint64, error) {
	query := `SELECT banner_id FROM slot2banner WHERE slot_id=:slotID LIMIT :slotLimit`
	arg := map[string]interface{}{
		"slotID":    slotID,
		"slotLimit": slotLimit,
	}
	rows, err := s.queryContext(ctx, query, arg)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var banners []uint64
	var bannerID uint64
	for rows.Next() {
		err = rows.Scan(&bannerID)
		if err != nil {
			return 0, err
		}
		banners = append(banners, bannerID)
	}

	if len(banners) == 0 {
		return 0, nil
	}

	return banners[0], nil
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
