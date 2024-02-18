package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"read-adviser-bot/lib/e"
	"read-adviser-bot/storage"
	"time"
)

type Storage struct {
	basePath string
}

func New(basePath string) *Storage {
	return &Storage{basePath: basePath}
}

func (s *Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("Не могу сохранить статью", err) }()

	filePath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(filePath, 0774); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	filePath = filepath.Join(filePath, fName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("Не могу взять рандомную статью", err) }()

	filePath := filepath.Join(s.basePath, userName)

	if _, err := s.isFolderExists(userName); err != nil {
		return nil, err
	}

	files, err := os.ReadDir(filePath)

	fmt.Println(err)

	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())

	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(filePath, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("Не могу удалить статью", err)
	}

	filePath := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(filePath); err != nil {
		return e.Wrap(fmt.Sprintf("Не могу удалить статью: %s", filePath), err)
	}

	return nil
}

func (s Storage) isFolderExists(username string) (bool, error) {
	folderPath := filepath.Join("filesStorage", username)

	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return false, e.Wrap("Вы еще не сохранили ни одной статьи", err)
	}

	return true, nil
}

func (s Storage) IsExist(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("Не могу найти статью", err)
	}

	filePath := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(filePath); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, e.Wrap(fmt.Sprintf("Не могу проверить существует ли данный файл: %s", filePath), err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("Не могу декодировать статью", err)
	}

	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("Не могу декодировать статью", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
