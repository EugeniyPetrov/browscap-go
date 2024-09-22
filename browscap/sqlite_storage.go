package browscap

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type SqliteBrowserStorage struct {
	*AbstractDBStorage
}

func NewSqliteBrowserStorage(db *sqlx.DB) *SqliteBrowserStorage {
	s := &SqliteBrowserStorage{}
	s.AbstractDBStorage = NewAbstractDBStorage(
		db, s, s, s, s,
		NewFixedPlaceholderMaker("?"),
	)
	return s
}

func (s *SqliteBrowserStorage) CreateTable(name string, columns string) error {
	query := fmt.Sprintf(`
		CREATE TABLE %s (%s)`,
		s.QuoteMeta(name),
		columns,
	)

	_, err := s.db.Exec(query)
	return err
}

func (s *SqliteBrowserStorage) QuoteMeta(m string) string {
	return m
}

func (s *SqliteBrowserStorage) InsertIgnore(table string, columns []string, data any) error {
	_, err := s.db.NamedExec(
		fmt.Sprintf("INSERT OR IGNORE INTO %s (%s) VALUES (%s)",
			table,
			s.columnsToSql(columns),
			s.columnsToPlaceholders(columns),
		),
		data,
	)
	return err
}

func (s *SqliteBrowserStorage) TableExists(table string) (bool, error) {
	var tableExists bool
	err := s.db.Get(
		&tableExists,
		`SELECT COUNT(*) > 0 FROM sqlite_master WHERE type = 'table' AND name = ?`,
		table,
	)

	return tableExists, err
}
