package raftsqlite

import (
	"database/sql"
	"encoding/binary"
	"fmt"
	"github.com/hashicorp/raft"

	// sql lite driver import
	_ "github.com/mattn/go-sqlite3"
)

var (
	ErrKeyNotFound = fmt.Errorf("requested key not found")
)

const (
	driverName = "sqlite3"
	dbName     = "raft-sqlite.db"
)

// Queries
const (
	firstIndexQuery     = `SELECT ifnull(min(l_index), 0) from r_log`
	lastIndexQuery      = `SELECT ifnull(max(l_index), 0) from r_log`
	getLogForIndexQuery = `SELECT l_index, term, type, data FROM r_log WHERE l_index = ?`
	storeLogQuery       = `REPLACE INTO r_log (l_index, term, type, data) VALUES (?, ?, ?, ?)`
	deleteRangeQuery    = `DELETE FROM r_log WHERE l_index >= ? AND l_index <= ?`
	setQuery            = `REPLACE INTO r_store (key, value) VALUES (?, ?)`
	getQuery            = `SELECT min(value) FROM r_store WHERE key = ?`
)

// SQL schema
var schemaQueries = []string{
	`
		CREATE TABLE IF NOT EXISTS r_log (
			l_index integer,
			term bigint not null,
			type int not null,
			data blob,
			PRIMARY KEY (l_index)
		)
	`,
	`
		CREATE TABLE IF NOT EXISTS r_store (
			id integer,
			key varbinary(512) not null,
			value blob not null,
			PRIMARY KEY (id)
		)
	`,
	`
		CREATE INDEX IF NOT EXISTS r_store_key_idx ON r_store(key)
	`,
}

type SQLStore struct {
	db   *sql.DB
	path string
}

func NewSQLStore(path string) (*SQLStore, error) {
	dbFullPath := fmt.Sprintf("%s/%s", path, dbName)
	db, err := newDB(dbFullPath)
	if err != nil {
		return nil, err
	}

	store := SQLStore{
		db:   db,
		path: dbFullPath,
	}

	for _, query := range schemaQueries {
		_, err := db.Exec(query)
		if err != nil {
			return nil, err
		}
	}

	return &store, nil
}

func newDB(path string) (*sql.DB, error) {
	// TODO: expose the sqlite PRAGMA to outside env
	path = fmt.Sprintf("%s?%s", path, "_locking_mode=EXCLUSIVE")
	db, err := sql.Open(driverName, path)
	if err != nil {
		return nil, err
	}

	// TODO: any room to optimize ?
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	return db, nil
}

func (s SQLStore) FirstIndex() (uint64, error) {
	var idx uint64
	err := s.db.QueryRow(firstIndexQuery).Scan(&idx)
	return idx, err
}

func (s SQLStore) LastIndex() (uint64, error) {
	var idx uint64
	err := s.db.QueryRow(lastIndexQuery).Scan(&idx)
	return idx, err
}

func (s SQLStore) GetLog(index uint64, log *raft.Log) error {
	err := s.db.QueryRow(getLogForIndexQuery, index).Scan(&log.Index, &log.Term, &log.Type, &log.Data)
	if err == sql.ErrNoRows {
		return raft.ErrLogNotFound
	}
	return err
}

func (s SQLStore) StoreLog(log *raft.Log) error {
	return s.StoreLogs([]*raft.Log{log})
}

func (s SQLStore) StoreLogs(logs []*raft.Log) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(storeLogQuery)
	if err != nil {
		return err
	}

	for _, log := range logs {
		_, err := stmt.Exec(log.Index, log.Term, log.Type, log.Data)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s SQLStore) DeleteRange(min, max uint64) error {
	_, err := s.db.Exec(deleteRangeQuery, min, max)
	return err
}

func (s SQLStore) Set(key, val []byte) error {
	tx, err := s.db.Begin()

	_, err = tx.Exec(setQuery, string(key), val)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s SQLStore) Get(key []byte) ([]byte, error) {
	var value []byte
	err := s.db.QueryRow(getQuery, string(key)).Scan(&value)
	if err != nil {
		return nil, err
	}

	if value == nil {
		return nil, ErrKeyNotFound
	}
	return value, nil
}

func (s SQLStore) SetUint64(key []byte, val uint64) error {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, val)
	return s.Set(key, b)
}

func (s SQLStore) GetUint64(key []byte) (uint64, error) {
	b, err := s.Get(key)
	if err != nil {
		return 0, err
	}
	if len(b) == 0 {
		return 0, ErrKeyNotFound
	}
	val := binary.LittleEndian.Uint64(b)
	return val, nil
}

func (s SQLStore) Close() error {
	return s.db.Close()
}
