package main

// create cache
// write -- save, update one, delete all cache
// read from cache

/** Writing
- on application start, load data from cache into a map
- on save, update, delete;
	- save to cache
	- update map
- essentially a map should be updated from the cache on each write
*/

/** On read
- on application start, load data from cache into a map
- on read;
	- get data from map
- every write ensures to sync data between cache and map
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Storage struct {
	Rqs map[string]interface{}
	FileCache *os.File
}

func (s *Storage) load() error {

	// after loading data into in-memory storage, close file
	// only open when writing and close after (not sure if this is efficient?)
	// defer s.fileCache.Close() // inefficient

	// read from cache into map
	_, err := s.FileCache.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("load error - seek: %w", err)
	}

	// buffer := make([]byte, 64) // what if this is larger than 64bytes
	// _, err = s.fileCache.Read(buffer)
	// if err != nil {
	// 	return fmt.Errorf("load error - : %w", err)
	// }

	// readall file content into memory
	b, err := io.ReadAll(s.FileCache)
	if err != nil {
		return fmt.Errorf("load error - readall: %w", err)
	}

	s.Rqs = map[string]interface{}{}

	if len(b) == 0 { return nil }

	if err = json.Unmarshal(b, &s.Rqs); err != nil {
		return fmt.Errorf("load error - unmarshal: %w", err)
	}

	return nil
}

func NewStore(filename string) (*Storage, error) {
	storage := new(Storage)

	file, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	storage.FileCache = file

	if err = storage.load(); err != nil {
		return nil, err
	}

	return storage, nil
} 

func (s *Storage) DeleteAll() error {
	
	s.clearCacheFile()
	clear(s.Rqs)

	return nil
}

func (s *Storage) clearCacheFile() error {
	err := s.FileCache.Truncate(0)
	if err != nil {
		return fmt.Errorf("(deleteAll) truncate failed: %w", err)
	}

	if _, err = s.FileCache.Seek(0,0); err != nil {
		return fmt.Errorf("(deleteAll) seek failed: %w", err)
	}

	if err = s.FileCache.Sync(); err != nil {
		return fmt.Errorf("(deleteAll) sync failed: %w", err)
	}

	return nil
}