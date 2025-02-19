package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

var origin *string;
var store *Storage


func writeToFile(mapping map[string]interface{}) error {

	// add to the 
	// b, err := json.MarshalIndent(mapping, "", " ")
	// if err != nil {
	// 	return fmt.Errorf("writeFileError - marshalIndent %w", err)
	// }
	// fmt.Printf("marshalled indent data %v, %v", b, mapping)

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", " ")
	if err := enc.Encode(mapping); err != nil {
		return fmt.Errorf("writeFileError - encoding %w", err)
	}

	store.clearCacheFile()
	
	if _, err := store.FileCache.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("writeFileError - write %w", err)
	}
	fmt.Println("Successfully cached")
	return nil
}

func cacheResp(url string, resp []byte) error {

	var mapping interface{} 
	err := json.Unmarshal(resp, &mapping)
	if err != nil {
		return fmt.Errorf("error (cacheResp - unmarshal): %w", err)
	}
	store.Rqs[url] = mapping
	fmt.Printf("stored requests %v", store.Rqs)

	if err := writeToFile(store.Rqs); err != nil {
		return err
	}

	return nil
}

func handleGetRequests(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Path[1:]
	query := r.URL.RawQuery

	url := fmt.Sprintf("%s%s?%s", *origin, params, query)
	fmt.Println(url)

	// before making a request, get response from cache 
	
	resp, err := http.Get(url)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	// b, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println(err)
	// } 
	// fmt.Printf("BODY %s", string(b))

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// result := map[string]interface{}{}
	// err = json.NewDecoder(resp.Body).Decode(result)
	

	if err = cacheResp(url, b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%v", string(b))
}

func handleClearCache(w http.ResponseWriter, r *http.Request) {
	if err := store.DeleteAll(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Cache cleared successfully")
}

func main() {
	port := flag.String("port", "3000", "port to start server on")
	origin = flag.String(
		"origin", "https://dummyjson.com/", "origin to forward requests to")

	flag.Parse()

	*port = fmt.Sprintf(":%s", *port)

	if *origin == "" {
		log.Fatal("no origin url provided")
	}

	ns, err := NewStore("store.json")
	// check for file error err.Is and close file
	if err != nil {
		log.Fatal(err)
	} else {
		store = ns
		fmt.Println(store)
	}
	
	// fmt.Printf("port: %s, origin: %s", *port, *origin)
	fmt.Println("Server listening...")

	http.HandleFunc("/products", handleGetRequests)
	http.HandleFunc("/clear", handleClearCache)

	err = http.ListenAndServe(*port, nil)

	if err != nil {
		log.Fatal(err)
	}
}