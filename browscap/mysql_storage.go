package browscap

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type MysqlBrowserStorage struct {
	*AbstractDBStorage
}

func NewMysqlBrowserStorage(db *sqlx.DB) *MysqlBrowserStorage {
	s := &MysqlBrowserStorage{}
	s.AbstractDBStorage = NewAbstractDBStorage(
		db, s, s, s, s,
		NewFixedPlaceholderMaker("?"),
	)
	return s
}

func (s *MysqlBrowserStorage) CreateTable(name string, columns string) error {
	query := fmt.Sprintf(`
		CREATE TABLE %s (%s) ENGINE=InnoDB ROW_FORMAT=COMPRESSED`,
		s.QuoteMeta(name),
		columns,
	)

	_, err := s.db.Exec(query)
	return err
}

func (s *MysqlBrowserStorage) QuoteMeta(m string) string {
	return fmt.Sprintf("`%s`", m)
}

func (s *MysqlBrowserStorage) InsertIgnore(table string, columns []string, data any) error {
	_, err := s.db.NamedExec(
		fmt.Sprintf("INSERT IGNORE INTO %s (%s) VALUES (%s)",
			table,
			s.columnsToSql(columns),
			s.columnsToPlaceholders(columns),
		),
		data,
	)
	return err
}

func (s *MysqlBrowserStorage) TableExists(table string) (bool, error) {
	var tableExists bool
	err := s.db.Get(
		&tableExists,
		`SELECT COUNT(*) > 0 FROM information_schema.TABLES WHERE TABLE_SCHEMA = SCHEMA() AND TABLE_NAME = ?`,
		table,
	)

	return tableExists, err
}
