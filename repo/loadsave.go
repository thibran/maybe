package repo

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"thibaut/maybe/ratedfolder"
)

// LoadData from dataDir or create directory.
func (r *Repo) LoadData(dataDir string) {
	if err := r.Load(); err != nil {
		if err != errNoFile {
			log.Fatalln(err)
		}
		// create data dir, if not existent
		if err := os.MkdirAll(dataDir, 0770); err != nil {
			log.Fatalf("main - create data dir: %s\n", err)
		}
	}
}

// Save repo map to dataPath.
func (r *Repo) Save() error {
	f, err := os.Create(r.dataPath)
	if err != nil {
		log.Fatalf("could not save filerepo: %s %v\n", r.dataPath, err)
	}
	defer f.Close()
	return saveGzip(f, r.m)
}

func saveGzip(w io.Writer, data ratedfolder.Map) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(data)
	wg := gzip.NewWriter(w)
	defer wg.Close()
	wg.Write(b.Bytes())
	return nil
}

// Load repo map from dataPath.
func (r *Repo) Load() error {
	f, err := os.Open(r.dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			return errNoFile
		}
		return err
	}
	defer f.Close()
	m, err := loadGzip(f)
	if err != nil {
		return err
	}
	r.m = m
	return nil
}

func loadGzip(r io.Reader) (ratedfolder.Map, error) {
	gr, err := gzip.NewReader(r)
	defer gr.Close()
	if err != nil {
		return nil, err
	}
	var m ratedfolder.Map
	dec := gob.NewDecoder(gr)
	if err := dec.Decode(&m); err != nil {
		return nil, fmt.Errorf("could not decode: %v", err)
	}
	return m, nil
}
