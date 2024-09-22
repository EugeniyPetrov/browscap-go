package browscap

import (
	"fmt"
	radix "github.com/eugeniypetrov/radix-tree"
	"reflect"
	"sort"
	"strings"
)

const DefaultPatternName = "defaultproperties"

var ErrNotFound = fmt.Errorf("browser not found")

type Browscap struct {
	tree           *radix.Node
	browserStorage BrowserStorage
}

func NewBrowscap(tree *radix.Node, browserStorage BrowserStorage) *Browscap {
	return &Browscap{
		tree:           tree,
		browserStorage: browserStorage,
	}
}

func (b *Browscap) mergeBrowsers(src, dest *BrowserNode) {
	srcVal := reflect.ValueOf(src).Elem()
	destVal := reflect.ValueOf(dest).Elem()

	for i := 0; i < destVal.NumField(); i++ {
		destField := destVal.Field(i)
		srcField := srcVal.Field(i)

		notSet := false
		// for pointers, we only set if the destination is nil
		if destField.Kind() == reflect.Ptr && destField.IsNil() && !srcField.IsNil() {
			notSet = true
		}

		// for values, we only set if the destination is the zero value
		if destField.Kind() != reflect.Ptr && destField.IsZero() {
			notSet = true
		}

		if notSet {
			destField.Set(srcField)
		}
	}
}

func (b *Browscap) loadBrowserRecursive(pattern string) (*Browser, error) {
	res := &BrowserNode{}

	for {
		browser, err := b.browserStorage.Get(pattern)
		if err != nil {
			return nil, fmt.Errorf("error getting browser for pattern %s: %w", pattern, err)
		}

		b.mergeBrowsers(browser, res)

		if pattern == DefaultPatternName {
			break
		}

		pattern = browser.Parent
	}

	return res.ToBrowser(), nil
}

func (b *Browscap) GetBrowser(ua string) (*Browser, error) {
	ua = strings.ToLower(ua)
	patterns := b.tree.Find(ua)
	sort.Sort(Patterns(patterns))

	for _, p := range patterns {
		return b.loadBrowserRecursive(p)
	}

	return &Browser{}, ErrNotFound
}
