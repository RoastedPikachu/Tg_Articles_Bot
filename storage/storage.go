package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"read-adviser-bot/lib/e"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExist(p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("Нет сохраненных статей")

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("Не могу вычислить хеш", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("Не могу вычислить хеш", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
