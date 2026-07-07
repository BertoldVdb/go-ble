package blesmp

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io/fs"
	"os"
	"sync"
)

// KeyStore is the backing store SMP uses to persist its long-term key
// database. The data argument is an opaque byte blob produced by the
// SMP layer; implementations only need to round-trip the bytes between
// Save and Load. Load should return (nil, nil) when nothing has been
// stored yet (e.g. first run) — a non-nil error aborts SMP startup.
//
// SMP serialises calls to Load and Save itself, so implementations do
// not need their own concurrency control.
type KeyStore interface {
	Load() ([]byte, error)
	Save(data []byte) error
}

// FileKeyStore is the default disk-backed KeyStore. It writes the blob
// to Path atomically (via a sibling .tmp + rename) with mode 0600. An
// empty Path turns the store into a no-op, matching the historical
// behaviour of leaving StoredKeysPath unset (keys live only in memory
// for the current process).
type FileKeyStore struct {
	Path string
}

// Load reads the previously persisted blob, returning (nil, nil) if the
// file does not yet exist.
func (f *FileKeyStore) Load() ([]byte, error) {
	if f.Path == "" {
		return nil, nil
	}
	data, err := os.ReadFile(f.Path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	return data, err
}

// Save atomically writes the blob to Path.
func (f *FileKeyStore) Save(data []byte) error {
	if f.Path == "" {
		return nil
	}
	tmp := f.Path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, f.Path)
}

// smpKeyPersist wraps a KeyStore with the in-memory map and the mutex
// that the SMP code paths take while reading or mutating storedKeys.
// Save serialises the current map via gob (matching the historical
// on-disk format produced by gobpersist) and hands the bytes to the
// underlying KeyStore.
type smpKeyPersist struct {
	sync.Mutex
	target *map[smpStoredLTKMapKey]smpStoredLTK
	store  KeyStore
}

func (p *smpKeyPersist) Load() error {
	blob, err := p.store.Load()
	if err != nil || len(blob) == 0 {
		return err
	}
	return gob.NewDecoder(bytes.NewReader(blob)).Decode(p.target)
}

func (p *smpKeyPersist) Save() error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(p.target); err != nil {
		return err
	}
	return p.store.Save(buf.Bytes())
}
