package browscap

import (
	"fmt"
	radix "github.com/eugeniypetrov/radix-tree"
	"reflect"
	"sort"
	"strings"
)

const DefaultPatternName = "defaultproperties"

var DefaultBrowser = &BrowserNode{
	Comment:                    StringPtr("Default Browser"),
	Browser:                    StringPtr("Default Browser"),
	BrowserType:                StringPtr("unknown"),
	BrowserBits:                IntPtr(0),
	BrowserMaker:               StringPtr("unknown"),
	BrowserModus:               StringPtr("unknown"),
	Version:                    StringPtr("0.0"),
	MajorVer:                   StringPtr("0"),
	MinorVer:                   StringPtr("0"),
	Platform:                   StringPtr("unknown"),
	PlatformVersion:            StringPtr("unknown"),
	PlatformDescription:        StringPtr("unknown"),
	PlatformBits:               IntPtr(0),
	PlatformMaker:              StringPtr("unknown"),
	Alpha:                      BoolPtr(false),
	Beta:                       BoolPtr(false),
	Win16:                      BoolPtr(false),
	Win32:                      BoolPtr(false),
	Win64:                      BoolPtr(false),
	Frames:                     BoolPtr(false),
	Iframes:                    BoolPtr(false),
	Tables:                     BoolPtr(false),
	Cookies:                    BoolPtr(false),
	BackgroundSounds:           BoolPtr(false),
	Javascript:                 BoolPtr(false),
	VBScript:                   BoolPtr(false),
	JavaApplets:                BoolPtr(false),
	ActiveXControls:            BoolPtr(false),
	IsMobileDevice:             BoolPtr(false),
	IsTablet:                   BoolPtr(false),
	IsSyndicationReader:        BoolPtr(false),
	Crawler:                    BoolPtr(false),
	IsFake:                     BoolPtr(false),
	IsAnonymized:               BoolPtr(false),
	IsModified:                 BoolPtr(false),
	CSSVersion:                 IntPtr(0),
	AolVersion:                 IntPtr(0),
	DeviceName:                 StringPtr("unknown"),
	DeviceMaker:                StringPtr("unknown"),
	DeviceType:                 StringPtr("unknown"),
	DevicePointingMethod:       StringPtr("unknown"),
	DeviceCodeName:             StringPtr("unknown"),
	DeviceBrandName:            StringPtr("unknown"),
	RenderingEngineName:        StringPtr("unknown"),
	RenderingEngineVersion:     StringPtr("unknown"),
	RenderingEngineDescription: StringPtr("unknown"),
	RenderingEngineMaker:       StringPtr("unknown"),
}

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

	b.mergeBrowsers(DefaultBrowser, res)

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
