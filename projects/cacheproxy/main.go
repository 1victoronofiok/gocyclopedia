package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
)

var origin *string
var store *Storage

func handleGetRequests(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Path[1:]
	query := r.URL.RawQuery

	url := fmt.Sprintf("%s%s?%s", *origin, params, query)
	fmt.Println(url)

	if cacheRes, err := store.Get(url); err == nil && cacheRes != nil {
		// fmt.Fprintf(w, "%s", cacheRes)

		b, err := json.Marshal(cacheRes)
		if err != nil {
			http.Error(w, "Failed to marshal cache response", http.StatusInternalServerError)
			return
		}
		fmt.Println("Response from cache")
		w.Header().Set("X-Cache", "HIT") // headers should be set before write header else they might not be applied
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode == 200 {
		if err = store.Save(url, b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// fmt.Fprintf(w, "%v", string(b))
	w.Header().Set("X-Cache", "MISS")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(b)
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
