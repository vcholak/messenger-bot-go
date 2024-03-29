package files

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/vcholak/messenger-bot/lib/errp"
	"github.com/vcholak/messenger-bot/storage"
)

type FileStorage struct {
	basePath string
}

const defaultPerm = 0774

// New creates a new FileStorage.
func New(basePath string) FileStorage {
	return FileStorage{basePath: basePath}
}

// Save saves a page to the storage.
func (s FileStorage) Save(_ context.Context, page *storage.Page) (err error) {
	defer func() { err = errp.WrapIfErr("can't save page", err) }()

	fPath := filepath.Join(s.basePath, page.FirstName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

// PickRandom picks a random page from the storage.
func (s FileStorage) PickRandom(f_ context.Context, firstName string) (page *storage.Page, err error) {
	defer func() { err = errp.WrapIfErr("can't pick random page", err) }()

	path := filepath.Join(s.basePath, firstName)
	log.Printf("Trying to pick up a random page from the path: %s", path)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

// Remove removes a page from the storage.
func (s FileStorage) Remove(_ context.Context, p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return errp.Wrap("can't remove page", err)
	}

	path := filepath.Join(s.basePath, p.FirstName, fileName)

	if err := os.Remove(path); err != nil {
		msg := fmt.Sprintf("can't remove page %s", path)

		return errp.Wrap(msg, err)
	}

	return nil
}

// IsExists checks if a page exists in the storage.
func (s FileStorage) IsExists(_ context.Context, p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, errp.Wrap("can't check if page exists", err)
	}

	path := filepath.Join(s.basePath, p.FirstName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := fmt.Sprintf("can't check if file %s exists", path)

		return false, errp.Wrap(msg, err)
	}

	return true, nil
}

func (s FileStorage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, errp.Wrap("can't decode page", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, errp.Wrap("can't decode page", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
