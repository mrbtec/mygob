package mygob

import (
	"encoding/gob"
	"fmt"
	"os"
	"sync"
)

type GobStore struct {
	sync.RWMutex
	data     map[string]interface{}
	filename string
}

func NewGobStore(path string) *GobStore {
	return &GobStore{
		data:     make(map[string]interface{}),
		filename: path,
	}
}

func (s *GobStore) Set(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = value
}

func (s *GobStore) SetAndSave(key string, value interface{}) {
	s.Lock()
	defer s.Unlock()
	s.data[key] = value

	if err := s.SaveToFile(); err != nil {
		fmt.Println("Err SetAndSave ", err)
	}
}

func (s *GobStore) SetMulti(items map[string]interface{}) {
	s.Lock()
	defer s.Unlock()
	for key, value := range items {
		s.Set(key, value)
	}
	if err := s.SaveToFile(); err != nil {
		fmt.Println(err)
	}
}

func (s *GobStore) GetAll() map[string]interface{} {
	s.Lock()
	defer s.Unlock()
	result := make(map[string]interface{})
	for k, v := range s.data {
		result[k] = v
	}
	return result
}

func (s *GobStore) Get(key string) (interface{}, bool) {
	s.RLock()
	defer s.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *GobStore) GetOrSet(key string, value interface{}) interface{} {
	var v interface{}
	var ok bool
	if v, ok = s.Get(key); !ok {
		s.Set(key, value)
		if err := s.SaveToFile(); err != nil {
			fmt.Println("GetOrSet", err)
		}
	} else {
		return v
	}
	return v
}

func (s *GobStore) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	delete(s.data, key)
}

func (s *GobStore) SaveToFile() error {
	s.RLock()
	defer s.RUnlock()

	file, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(s.data); err != nil {
		return err
	}

	return nil
}

func (s *GobStore) LoadFromFile() error {
	s.Lock()
	defer s.Unlock()

	file, err := os.Open(s.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&s.data); err != nil {
		return err
	}

	return nil
}

func (s *GobStore) CreatePath(path string) error {
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		// Diretório não existe, criar
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("falha ao criar o diretório: %s", err)
		}
	}
	return nil
}

func (s *GobStore) CreateFile() error {
	if _, err := os.Stat(s.filename); os.IsNotExist(err) {
		file, err := os.Create(s.filename)
		if err != nil {
			return err
		}
		file.Close()
	}
	return nil
}

func (s *GobStore) Count() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.data)
}
func (s *GobStore) HasKey(key string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.data[key]
	return ok
}
func (s *GobStore) ListKeys() ([]string, error) {
	s.RLock()
	defer s.RUnlock()

	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}

	return keys, nil
}
func (s *GobStore) DeleteAll() error {
	s.Lock()
	defer s.Unlock()
	s.data = make(map[string]interface{})
	if err := s.SaveToFile(); err != nil {
		return err
	}
	return nil
}

func (s *GobStore) getCurrentDir() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return wd
}
