package browscap

import (
	"fmt"
	"github.com/zeebo/xxh3"
	"iter"
)

type MemoryBrowserStorage struct {
	version  *Version
	browsers map[uint64]*BrowserNode
}

func NewMemoryBrowserStorage() *MemoryBrowserStorage {
	return &MemoryBrowserStorage{
		browsers: make(map[uint64]*BrowserNode),
	}
}

func (s *MemoryBrowserStorage) Prepare() error {
	return nil
}

func (s *MemoryBrowserStorage) GetVersion() (*Version, error) {
	if s.version == nil {
		return nil, ErrEmptyCache
	}

	return s.version, nil
}

func (s *MemoryBrowserStorage) SaveVersion(ver *Version) error {
	s.version = ver
	return nil
}

func (s *MemoryBrowserStorage) hash(pattern string) uint64 {
	return xxh3.Hash([]byte(pattern))
}

func (s *MemoryBrowserStorage) Save(node *BrowserNode) error {
	hash := s.hash(node.Pattern)
	s.browsers[hash] = node
	return nil
}

func (s *MemoryBrowserStorage) Get(pattern string) (*BrowserNode, error) {
	hash := s.hash(pattern)
	node, ok := s.browsers[hash]
	if !ok {
		return nil, fmt.Errorf("pattern not found")
	}

	return node, nil
}

func (s *MemoryBrowserStorage) Patterns() iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		for _, b := range s.browsers {
			if !yield(b.Pattern, nil) {
				return
			}
		}
	}
}
