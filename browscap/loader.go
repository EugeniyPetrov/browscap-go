package browscap

import (
	"errors"
	"fmt"
	ini "github.com/eugeniypetrov/ini-reader"
	radix "github.com/eugeniypetrov/radix-tree"
	"github.com/go-viper/mapstructure/v2"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

const versionSectionName = "GJK_Browscap_Version"

type Version struct {
	Version int    `db:"version"`
	Type    string `db:"type"`
}

var ErrEmptyCache = fmt.Errorf("cache is empty")

type Loader struct {
	browserStorage BrowserStorage
}

func NewLoader(browserStorage BrowserStorage) *Loader {
	return &Loader{
		browserStorage: browserStorage,
	}
}

func (l *Loader) parseVersion(s *ini.Section) (*Version, error) {
	if s.Name != versionSectionName {
		return nil, fmt.Errorf("no version section")
	}

	if s.Properties["Format"] != "php" {
		return nil, fmt.Errorf("expected php format, but %s found", s.Properties["Format"])
	}

	verNum, ok := s.Properties["Version"].(int64)
	if !ok {
		return nil, fmt.Errorf("invalid version")
	}

	typeS, ok := s.Properties["Type"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid type")
	}

	return &Version{
		Version: int(verNum),
		Type:    typeS,
	}, nil
}

func (l *Loader) checkCache(ver *Version) error {
	cachedVer, err := l.browserStorage.GetVersion()
	if err != nil {
		return fmt.Errorf("error getting version from cache: %w", err)
	}

	if *ver != *cachedVer {
		return fmt.Errorf("version mismatch: %v != %v", *ver, *cachedVer)
	}

	return nil
}

func (l *Loader) normalizePattern(pattern string) string {
	return strings.ToLower(pattern)
}

func (l *Loader) browserNode(s *ini.Section) (*BrowserNode, error) {
	res := &BrowserNode{}
	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToBoolHookFunc(),
			mapstructure.StringToIntHookFunc(),
			Int64ToStringHookFunc(),
		),
		Result: res,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating decoder: %w", err)
	}

	err = d.Decode(s.Properties)
	if err != nil {
		return nil, fmt.Errorf("error decoding: %w", err)
	}

	res.Pattern = l.normalizePattern(s.Name)
	res.Parent = l.normalizePattern(res.Parent)

	return res, nil
}

func (l *Loader) storeCache(node *BrowserNode) error {
	err := l.browserStorage.Save(node)
	if err != nil {
		return fmt.Errorf("error saving cache: %w", err)
	}

	return nil
}

func (l *Loader) storeVersion(ver *Version) error {
	err := l.browserStorage.SaveVersion(ver)
	if err != nil {
		return fmt.Errorf("error saving version: %w", err)
	}

	return nil
}

func (l *Loader) makeCache(r *ini.Reader, ver *Version) error {
	err := l.browserStorage.Prepare()
	if err != nil {
		return fmt.Errorf("error preparing cache: %w", err)
	}

	for r.Next() {
		s := r.Section()

		if s.Name == versionSectionName {
			continue
		}

		node, err := l.browserNode(s)
		if err != nil {
			return fmt.Errorf("error creating browser node: %w", err)
		}

		err = l.storeCache(node)
		if err != nil {
			return fmt.Errorf("error storing cache: %w", err)
		}
	}

	err = r.Err()
	if err != nil {
		return fmt.Errorf("error reading ini: %w", err)
	}

	err = l.storeVersion(ver)
	if err != nil {
		return fmt.Errorf("error storing version: %w", err)
	}

	return nil
}

func (l *Loader) Compile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer f.Close()

	r := ini.NewReader(f)

	h := r.Next()
	if !h {
		return fmt.Errorf("unable to read ini file. %w", r.Err())
	}

	ver, err := l.parseVersion(r.Section())
	if err != nil {
		return fmt.Errorf("error parsing version: %w", err)
	}

	err = l.checkCache(ver)
	if err == nil {
		// already compiled
		return nil
	}

	if !errors.Is(err, ErrEmptyCache) {
		return fmt.Errorf("invalid cache: %w", err)
	}

	err = l.makeCache(r, ver)
	if err != nil {
		return fmt.Errorf("error making cache: %w", err)
	}

	return nil
}

func (l *Loader) Load() (*Browscap, error) {
	tree := radix.NewRadix()

	for pattern, err := range l.browserStorage.Patterns() {
		if err != nil {
			return nil, fmt.Errorf("error getting pattern: %w", err)
		}

		tree.Add(pattern)
	}

	return NewBrowscap(tree, l.browserStorage), nil
}
