package browscap

import "iter"

type BrowserStorage interface {
	Prepare() error
	GetVersion() (*Version, error)
	SaveVersion(ver *Version) error
	Save(node *BrowserNode) error
	Get(pattern string) (*BrowserNode, error)
	Patterns() iter.Seq2[string, error]
}
