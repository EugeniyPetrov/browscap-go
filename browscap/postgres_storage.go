package browscap

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type PostgresBrowserStorage struct {
	*AbstractDBStorage
}

func NewPostgresBrowserStorage(db *sqlx.DB) *PostgresBrowserStorage {
	s := &PostgresBrowserStorage{}
	s.AbstractDBStorage = NewAbstractDBStorage(
		db, s, s, s, s,
		NewNumberedPlaceholderMaker(),
	)
	return s
}

func (s *PostgresBrowserStorage) CreateTable(name string, columns string) error {
	query := fmt.Sprintf(`
		CREATE TABLE %s (%s)`,
		s.QuoteMeta(name),
		columns,
	)

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresBrowserStorage) QuoteMeta(m string) string {
	return fmt.Sprintf("\"%s\"", m)
}

func (s *PostgresBrowserStorage) InsertIgnore(table string, columns []string, data any) error {
	_, err := s.db.NamedExec(
		fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON CONFLICT DO NOTHING",
			table,
			s.columnsToSql(columns),
			s.columnsToPlaceholders(columns),
		),
		data,
	)
	return err
}

func (s *PostgresBrowserStorage) TableExists(table string) (bool, error) {
	var tableExists bool
	err := s.db.Get(
		&tableExists,
		`SELECT count(*) > 0 FROM information_schema.tables WHERE table_schema = current_schema AND table_name = $1`,
		table,
	)

	return tableExists, err
}
