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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

type Storage struct {
	Rqs       map[string]interface{}
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

	if len(b) == 0 {
		return nil
	}

	if err = json.Unmarshal(b, &s.Rqs); err != nil {
		return fmt.Errorf("load error - unmarshal: %w", err)
	}

	return nil
}

func NewStore(filename string) (*Storage, error) {
	storage := new(Storage)

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	storage.FileCache = file

	if err = storage.load(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) resetCacheFile() error {
	err := s.FileCache.Truncate(0)
	if err != nil {
		return fmt.Errorf("(deleteAll) truncate failed: %w", err)
	}

	if _, err = s.FileCache.Seek(0, 0); err != nil {
		return fmt.Errorf("(deleteAll) seek failed: %w", err)
	}

	if err = s.FileCache.Sync(); err != nil {
		return fmt.Errorf("(deleteAll) sync failed: %w", err)
	}

	return nil
}

func (s *Storage) DeleteAll() error {
	if err := s.resetCacheFile(); err != nil {
		return err
	}

	clear(s.Rqs)

	return nil
}

func (s *Storage) writeToCacheFile() error {
	// encode the in-memory store
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	if err := enc.Encode(s.Rqs); err != nil {
		return fmt.Errorf("writeFileError - encoding %w", err)
	}

	// reset cache
	s.resetCacheFile()

	// write encoded in-memory data to cache file
	if _, err := s.FileCache.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("writeFileError - write %w", err)
	}
	fmt.Println("Successfully cached")
	return nil
}

func (s *Storage) writeToMap(url string, resp []byte) error {
	// convert slice of bytes to interface
	var mapping interface{}
	err := json.Unmarshal(resp, &mapping)
	if err != nil {
		return fmt.Errorf("error (cacheResp - unmarshal): %w", err)
	}
	s.Rqs[url] = mapping
	fmt.Printf("stored requests %v", s.Rqs)

	return nil
}

func (s *Storage) Save(url string, resp []byte) error {
	err := s.writeToMap(url, resp)
	if err != nil {
		return err
	}

	if err := s.writeToCacheFile(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) getInMemResponse(url string) (resp interface{}, found bool) {
	resp, found = s.Rqs[url]
	return
}

func (s *Storage) getAllFromCache() (map[string]interface{}, error) {
	mapping := map[string]interface{}{}

	decoder := json.NewDecoder(s.FileCache)
	if err := decoder.Decode(&mapping); err != nil {
		return nil, fmt.Errorf("error (getAllFromCache - Decode) %w", err)
	}

	return mapping, nil
}

func (s *Storage) sync(mapping map[string]interface{}) error {
	// will implement tthis
	fmt.Printf("applied sync %v ", mapping)
	return nil
}

// check if resp for url can be found in in-mem store
// if found return,
// else check cache
// if found sync (should be a separate goroutine in order not to block)
// return resp
func (s *Storage) Get(url string) (interface{}, error) {
	if resp, found := s.getInMemResponse(url); found {
		return json.Marshal(resp)
	}

	// get all from cache in order to sync in memory data
	mapping, err := s.getAllFromCache()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	resp, found := mapping[url]
	if !found {
		return nil, nil
	}

	// found then sync
	if err = s.sync(mapping); err != nil {
		log.Print("Failed to sync")
	}

	// return resp
	return resp, nil
}

// urls -> urlkeystore

// id: {origin: "", params: "", query: "", headers: "", }
// [prod1, prod2]
// post [prod1, prod2, prod3]
// should be able to invalidate data
