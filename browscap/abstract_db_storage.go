package browscap

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/jmoiron/sqlx"
	"iter"
	"strconv"
	"sync/atomic"
)

type MetaQuoter interface {
	QuoteMeta(string) string
}

type Inserter interface {
	InsertIgnore(table string, columns []string, data any) error
}

type TableExistenceChecker interface {
	TableExists(table string) (bool, error)
}

type TableCreator interface {
	CreateTable(name string, columns string) error
}

type PlaceholderMaker interface {
	MakePlaceholder(from int) string
}

type FixedPlaceholderMaker struct {
	char string
}

func NewFixedPlaceholderMaker(char string) *FixedPlaceholderMaker {
	return &FixedPlaceholderMaker{char: char}
}

func (f *FixedPlaceholderMaker) MakePlaceholder(from int) string {
	return f.char
}

type NumberedPlaceholderMaker struct{}

func NewNumberedPlaceholderMaker() *NumberedPlaceholderMaker {
	return &NumberedPlaceholderMaker{}
}

func (n *NumberedPlaceholderMaker) MakePlaceholder(from int) string {
	return "$" + strconv.Itoa(from+1)
}

type AbstractDBStorage struct {
	db                    *sqlx.DB
	incrementCounter      int32
	inserter              Inserter
	metaQuoter            MetaQuoter
	tableExistenceChecker TableExistenceChecker
	tableCreator          TableCreator
	placeholderMaker      PlaceholderMaker
}

func NewAbstractDBStorage(
	db *sqlx.DB,
	inserter Inserter,
	metaQuoter MetaQuoter,
	tableExistenceChecker TableExistenceChecker,
	tableCreator TableCreator,
	placeholderMaker PlaceholderMaker,
) *AbstractDBStorage {
	return &AbstractDBStorage{
		db:                    db,
		inserter:              inserter,
		metaQuoter:            metaQuoter,
		tableExistenceChecker: tableExistenceChecker,
		tableCreator:          tableCreator,
		placeholderMaker:      placeholderMaker,
	}
}

func (s *AbstractDBStorage) Patterns() iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		from := 0
		batchSize := 500

		for {
			rows, err := s.db.Queryx(
				fmt.Sprintf(
					`SELECT id, pattern FROM browser WHERE id > %s ORDER BY id LIMIT %s`,
					s.placeholderMaker.MakePlaceholder(0),
					s.placeholderMaker.MakePlaceholder(1),
				),
				from,
				batchSize,
			)
			if err != nil {
				yield("", fmt.Errorf("error getting patterns: %w", err))
				return
			}
			defer rows.Close()

			count := 0
			for rows.Next() {
				var id int
				var pattern string

				err = rows.Scan(&id, &pattern)
				if err != nil {
					yield("", fmt.Errorf("error scanning pattern: %w", err))
					return
				}

				if !yield(pattern, nil) {
					return
				}

				from = id
				count++
			}

			err = rows.Err()
			if err != nil {
				yield("", fmt.Errorf("error iterating patterns: %w", err))
				return
			}

			if count < batchSize {
				break
			}
		}
	}
}

func (s *AbstractDBStorage) hash(pattern string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(pattern)))
}

func (s *AbstractDBStorage) SaveVersion(ver *Version) error {
	_, err := s.db.NamedExec(`INSERT INTO version (version, type) VALUES (:version, :type)`, ver)
	if err != nil {
		return fmt.Errorf("error saving version: %w", err)
	}

	return nil
}

func (s *AbstractDBStorage) Get(pattern string) (*BrowserNode, error) {
	hash := s.hash(pattern)

	node := new(BrowserNode)
	err := s.db.Get(
		node,
		fmt.Sprintf(`
			SELECT
				id, parent, pattern, comment, browser, browser_type, browser_bits, browser_maker, browser_modus,
				version, major_ver, minor_ver, platform, platform_version, platform_description, platform_bits,
				platform_maker, alpha, beta, win16, win32, win64, frames, iframes, tables, cookies, background_sounds,
				javascript, vbscript, java_applets, activex_controls, is_mobile_device, is_tablet,
				is_syndication_reader, crawler, is_fake, is_anonymized, is_modified, css_version, aol_version,
				device_name, device_maker, device_type, device_pointing_method, device_code_name, device_brand_name,
				rendering_engine_name, rendering_engine_version, rendering_engine_description, rendering_engine_maker
			FROM browser WHERE hash = %s`,
			s.placeholderMaker.MakePlaceholder(0),
		),
		hash,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting node: %w", err)
	}

	return node, nil
}

func (s *AbstractDBStorage) columnsToSql(columns []string) string {
	buf := bytes.NewBufferString("")
	for i, column := range columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(s.metaQuoter.QuoteMeta(column))
	}

	return buf.String()
}

func (*AbstractDBStorage) columnsToPlaceholders(columns []string) string {
	buf := bytes.NewBufferString("")
	for i, column := range columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(":")
		buf.WriteString(column)
	}

	return buf.String()
}

func (s *AbstractDBStorage) Save(node *BrowserNode) error {
	id := atomic.AddInt32(&s.incrementCounter, 1)

	n := &struct {
		ID   int    `db:"id"`
		Hash string `db:"hash"`
		*BrowserNode
	}{
		ID:          int(id),
		Hash:        s.hash(node.Pattern),
		BrowserNode: node,
	}

	err := s.inserter.InsertIgnore(
		"browser",
		[]string{"id", "hash", "parent", "pattern", "comment", "browser", "browser_type", "browser_bits",
			"browser_maker", "browser_modus", "version", "major_ver", "minor_ver", "platform", "platform_version",
			"platform_description", "platform_bits", "platform_maker", "alpha", "beta", "win16", "win32", "win64",
			"frames", "iframes", "tables", "cookies", "background_sounds", "javascript", "vbscript", "java_applets",
			"activex_controls", "is_mobile_device", "is_tablet", "is_syndication_reader", "crawler", "is_fake",
			"is_anonymized", "is_modified", "css_version", "aol_version", "device_name", "device_maker", "device_type",
			"device_pointing_method", "device_code_name", "device_brand_name", "rendering_engine_name",
			"rendering_engine_version", "rendering_engine_description", "rendering_engine_maker",
		},
		n,
	)
	if err != nil {
		return fmt.Errorf("error saving node: %w", err)
	}

	return nil
}

func (s *AbstractDBStorage) GetVersion() (*Version, error) {
	versionExists, err := s.tableExistenceChecker.TableExists("version")
	if err != nil {
		return nil, fmt.Errorf("error checking version table: %w", err)
	}

	if !versionExists {
		return nil, ErrEmptyCache
	}

	ver := new(Version)
	err = s.db.Get(ver, "SELECT version, type FROM version")
	if err != nil {
		return nil, fmt.Errorf("error getting version: %w", err)
	}

	return ver, nil
}

func (s *AbstractDBStorage) createVersionTable() error {
	err := s.tableCreator.CreateTable(
		"version",
		`version INT NOT NULL,
		type VARCHAR(255) NOT NULL`,
	)
	if err != nil {
		return fmt.Errorf("error creating version table: %w", err)
	}

	return err
}

func (s *AbstractDBStorage) createBrowserTable() error {
	err := s.tableCreator.CreateTable(
		"browser",
		`id INT NOT NULL PRIMARY KEY,
		hash CHAR(32) NOT NULL UNIQUE,
		parent VARCHAR(255) NOT NULL,
		pattern VARCHAR(255) NOT NULL,
		comment VARCHAR(255),
		browser VARCHAR(255),
		browser_type VARCHAR(255),
		browser_bits INT,
		browser_maker VARCHAR(255),
		browser_modus VARCHAR(255),
		version VARCHAR(255),
		major_ver VARCHAR(255),
		minor_ver VARCHAR(255),
		platform VARCHAR(255),
		platform_version VARCHAR(255),
		platform_description VARCHAR(255),
		platform_bits INT,
		platform_maker VARCHAR(255),
		alpha BOOL,
		beta BOOL,
		win16 BOOL,
		win32 BOOL,
		win64 BOOL,
		frames BOOL,
		iframes BOOL,
		tables BOOL,
		cookies BOOL,
		background_sounds BOOL,
		javascript BOOL,
		vbscript BOOL,
		java_applets BOOL,
		activex_controls BOOL,
		is_mobile_device BOOL,
		is_tablet BOOL,
		is_syndication_reader BOOL,
		crawler BOOL,
		is_fake BOOL,
		is_anonymized BOOL,
		is_modified BOOL,
		css_version INT,
		aol_version INT,
		device_name VARCHAR(255),
		device_maker VARCHAR(255),
		device_type VARCHAR(255),
		device_pointing_method VARCHAR(255),
		device_code_name VARCHAR(255),
		device_brand_name VARCHAR(255),
		rendering_engine_name VARCHAR(255),
		rendering_engine_version VARCHAR(255),
		rendering_engine_description VARCHAR(255),
		rendering_engine_maker VARCHAR(255)`,
	)
	if err != nil {
		return fmt.Errorf("error creating browser table: %w", err)
	}

	return err
}

func (s *AbstractDBStorage) Prepare() error {
	var err error

	err = s.createVersionTable()
	if err != nil {
		return fmt.Errorf("error creating version table: %w", err)
	}

	err = s.createBrowserTable()
	if err != nil {
		return fmt.Errorf("error creating browser table: %w", err)
	}

	return nil
}
