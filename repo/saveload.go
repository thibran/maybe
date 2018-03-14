package repo

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/thibran/maybe/rated"
)

// Save repo map to dataPath.
func (r *Repo) Save() error {
	f, err := os.Create(r.dataDir)
	if err != nil {
		log.Fatalf("could not save filerepo: %s %v\n", r.dataDir, err)
	}
	defer f.Close()
	return saveGzip(f, r.m)
}

func saveGzip(w io.Writer, data rated.Map) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	enc.Encode(data)
	wg := gzip.NewWriter(w)
	defer wg.Close()
	wg.Write(b.Bytes())
	return nil
}

// Load from dir, or create directory if not existent.
func (r *Repo) Load(datadir string) {
	if err := r.loadFile(); err != nil {
		if err != errNoFile {
			log.Fatalln(err)
		}
		// create data dir, if not existent
		if err := os.MkdirAll(datadir, 0770); err != nil {
			log.Fatalf("main - create data dir: %s\n", err)
		}
	}
}

// load repo map from dataPath.
func (r *Repo) loadFile() error {
	f, err := os.Open(r.dataDir)
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

func loadGzip(r io.Reader) (rated.Map, error) {
	gr, err := gzip.NewReader(r)
	defer gr.Close()
	if err != nil {
		return nil, err
	}
	var m rated.Map
	dec := gob.NewDecoder(gr)
	if err := dec.Decode(&m); err != nil {
		return nil, fmt.Errorf("could not decode: %v", err)
	}
	return m, nil
}
