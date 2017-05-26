package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"zalua/settings"
)

var Box = newStorage()

type storageItem struct {
	Value     string `json:"value"`
	CreatedAt int64  `json:"created_at"`
	TTL       int64  `json:"ttl"`
}

func (s *storageItem) valid() bool {
	return s.CreatedAt+s.TTL > time.Now().Unix()
}

type storage struct {
	sync.Mutex
	data map[string]*storageItem
}

// создание storage
func newStorage() *storage {
	result := &storage{data: make(map[string]*storageItem, 0)}
	// загрузка storage из файла
	if _, err := os.Stat(settings.StoragePath()); err == nil {
		if data, err := ioutil.ReadFile(settings.StoragePath()); err == nil {
			list := make(map[string]*storageItem, 0)
			if json.Unmarshal(data, &list); err == nil {
				for key, item := range list {
					if !item.valid() {
						delete(list, key)
					}
				}
				result.data = list
			} else {
				fmt.Fprintf(os.Stderr, "[ERROR] Storage unmarshal: %s\n", err.Error())
			}
		}
	}
	go result.junitor()
	go result.saver()
	return result
}

func (s *storage) junitor() {
	ticker := time.NewTicker(15 * time.Minute)
	for range ticker.C {
		log.Printf("[INFO] Start compact storage\n")
		s.Lock()
		for key, item := range s.data {
			if !item.valid() {
				delete(s.data, key)
			}
		}
		s.Unlock()
		log.Printf("[INFO] Compact storage done\n")
	}
}

func (s *storage) saver() {
	ticker := time.NewTicker(time.Minute)
	for range ticker.C {
		s.Lock()
		if data, err := json.Marshal(s.data); err == nil {
			tmpFile := settings.StoragePath() + ".tmp"
			if err := ioutil.WriteFile(tmpFile, data, 0644); err == nil {
				if err := os.Rename(tmpFile, settings.StoragePath()); err == nil {
					log.Printf("[INFO] Storage file saved\n")
				}
			} else {
				log.Printf("[ERROR] Storage write: %s\n", err.Error())
			}
		}
		s.Unlock()
	}
}

// установить значение по ключу с default TTL
func (p *storage) Set(key, val string, ttl int64) {
	p.Lock()
	defer p.Unlock()

	p.data[key] = &storageItem{
		Value:     val,
		CreatedAt: time.Now().Unix(),
		TTL:       ttl,
	}
}

// отдать значение по ключу
func (p *storage) Get(key string) (string, bool) {
	p.Lock()
	defer p.Unlock()

	item := p.data[key]
	if item == nil {
		return "", false
	}
	if !item.valid() {
		delete(p.data, key)
		return "", false
	}

	return item.Value, true
}

// удалить значение по ключу
func (p *storage) Delete(key string) {
	p.Lock()
	defer p.Unlock()

	delete(p.data, key)
}

// список всех ключей и значений
func (p *storage) List() []string {
	p.Lock()
	defer p.Unlock()

	result := []string{}
	for key, item := range p.data {
		if item.valid() {
			result = append(result, fmt.Sprintf("%s\t\t%s\t\t%d", key, item.Value, item.CreatedAt))
		}
	}
	return result
}
