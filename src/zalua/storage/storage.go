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

type StorageItem struct {
	ItemValue *StorageItemValue `json:"item_value"`
	CreatedAt int64             `json:"created_at"`
	TTL       int64             `json:"ttl"`
}

type StorageItemValue struct {
	Value string            `json:"value"`
	Tags  map[string]string `json:"tags"`
}

func (s *StorageItem) valid() bool {
	return s.CreatedAt+s.TTL > time.Now().Unix()
}

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
func (p *storage) Set(key, val string, tags map[string]string, ttl int64) {
	p.Lock()
	defer p.Unlock()

	metricKey := key
	if len(tags) > 0 {
		data, err := json.Marshal(&tags)
		if err == nil {
			metricKey = metricKey + string(data)
		}
	}

	p.data[metricKey] = &StorageItem{
		ItemValue: &StorageItemValue{
			Value: val,
			Tags:  tags,
		},
		CreatedAt: time.Now().Unix(),
		TTL:       ttl,
	}
}

// отдать значение по ключу
func (p *storage) Get(key string, tags map[string]string) (*StorageItemValue, bool) {
	p.Lock()
	defer p.Unlock()

	metricKey := key
	if len(tags) > 0 {
		data, err := json.Marshal(&tags)
		if err == nil {
			metricKey = metricKey + string(data)
		}
	}

	item := p.data[metricKey]
	if item == nil || item.ItemValue == nil {
		return nil, false
	}
	if !item.valid() {
		delete(p.data, metricKey)
		return item.ItemValue, false
	}

	return item.ItemValue, true
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
