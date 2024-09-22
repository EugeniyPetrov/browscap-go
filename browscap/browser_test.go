package browscap

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"
)

func compileAndLoad(filename string) (*Browscap, error) {
	storage := NewMemoryBrowserStorage()
	loader := NewLoader(storage)
	err := loader.Compile(filename)
	if err != nil {
		return nil, fmt.Errorf("error compiling: %w", err)
	}

	bc, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading: %w", err)
	}

	return bc, nil
}

func TestGetBrowser(t *testing.T) {
	bc, err := compileAndLoad("fixtures/lite_php_browscap.ini")
	if err != nil {
		t.Fatal(err)
	}

	b, err := bc.GetBrowser("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, b.Pattern, "mozilla/5.0 (*mac os x*) applewebkit* (*khtml*like*gecko*) chrome/128.0*safari/*")
}
