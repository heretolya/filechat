package cache

import (
	"os"
	"path/filepath"
)

type Service struct {
	prefix string
}

func New(prefix string) *Service {
	return &Service{prefix}
}

func (s *Service) cacheDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, s.prefix)
	err = os.MkdirAll(dir, 0755)
	return dir, err
}

func (s *Service) Get(key string) ([]byte, error) {
	dir, err := s.cacheDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(dir, key)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return data, err
}

func (s *Service) Set(key string, data []byte) error {
	dir, err := s.cacheDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, key)
	return os.WriteFile(path, data, 0644)
}

func (s *Service) Drop() error {
	base, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(base, s.prefix)
	return os.RemoveAll(dir)
}
