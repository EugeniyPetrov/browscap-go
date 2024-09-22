package browscap

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru/v2"
	"iter"
)

type LRUCachedStorage struct {
	cache   *lru.Cache[string, *BrowserNode]
	storage BrowserStorage
}

func NewLRUCachedStorage(storage BrowserStorage, cacheSize int) (*LRUCachedStorage, error) {
	cache, err := lru.New[string, *BrowserNode](cacheSize)
	if err != nil {
		return nil, fmt.Errorf("unable to create cache. %w", err)
	}

	return &LRUCachedStorage{
		cache:   cache,
		storage: storage,
	}, nil
}

func (s *LRUCachedStorage) Prepare() error {
	return s.storage.Prepare()
}

func (s *LRUCachedStorage) GetVersion() (*Version, error) {
	return s.storage.GetVersion()
}

func (s *LRUCachedStorage) SaveVersion(ver *Version) error {
	return s.storage.SaveVersion(ver)
}

func (s *LRUCachedStorage) Save(node *BrowserNode) error {
	return s.storage.Save(node)
}

func (s *LRUCachedStorage) Get(pattern string) (*BrowserNode, error) {
	if node, ok := s.cache.Get(pattern); ok {
		return node, nil
	}

	node, err := s.storage.Get(pattern)
	if err != nil {
		return nil, err
	}

	s.cache.Add(pattern, node)

	return node, nil
}

func (s *LRUCachedStorage) Patterns() iter.Seq2[string, error] {
	return s.storage.Patterns()
}
