package mygob
import (
	"encoding/gob"
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

func (s *GobStore) SetMultiGob(items map[string]interface{}) {
	for key, value := range items {
		s.Set(key, value)
		s.SaveToFile()
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
