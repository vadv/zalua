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

type storage struct {
	sync.Mutex
	data map[string]*StorageItem
}

// создание storage
func newStorage() *storage {
	result := &storage{data: make(map[string]*StorageItem, 0)}
	// загрузка storage из файла
	if _, err := os.Stat(settings.StoragePath()); err == nil {
		if data, err := ioutil.ReadFile(settings.StoragePath()); err == nil {
			list := make(map[string]*StorageItem, 0)
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
func (p *storage) Set(metric, val string, tags map[string]string, ttl int64) {
	p.Lock()
	defer p.Unlock()

	item := &StorageItem{
		Value:     val,
		Tags:      tags,
		Metric:    metric,
		CreatedAt: time.Now().Unix(),
		TTL:       ttl,
	}
	p.data[item.Key()] = item
}

// отдать значение по ключу
func (p *storage) Get(metric string, tags map[string]string) (*StorageItem, bool) {
	p.Lock()
	defer p.Unlock()

	metricKey := storageKey(metric, tags)

	item := p.data[metricKey]
	if item == nil {
		return nil, false
	}
	if !item.valid() {
		delete(p.data, metricKey)
		return item, false
	}

	return item, true
}

// удалить значение по ключу
func (p *storage) Delete(key string) {
	p.Lock()
	defer p.Unlock()

	delete(p.data, key)
}

// список всех ключей и значений
func (p *storage) List() map[string]*StorageItem {
	p.Lock()
	defer p.Unlock()

	result := make(map[string]*StorageItem, 0)
	for key, item := range p.data {
		if item.valid() {
			result[key] = item
		}
	}
	return result
}
