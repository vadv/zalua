package storage

import (
	"fmt"
	"log"
	"sync"
	"time"
)

var Box = newStorage()

type storageItem struct {
	value      string
	created_at int64
	ttl        int64
}

func (s *storageItem) valid() bool {
	return s.created_at+s.ttl > time.Now().Unix()
}

type storage struct {
	sync.Mutex
	data map[string]*storageItem
}

// создание storage
func newStorage() *storage {
	result := &storage{data: make(map[string]*storageItem, 0)}
	go result.junitor()
	return result
}

func (s *storage) junitor() {
	ticker := time.NewTicker(15 * time.Minute)
	for range ticker.C {
		log.Printf("[INFO] start compact storage")
		s.Lock()
		for key, item := range s.data {
			if !item.valid() {
				delete(s.data, key)
			}
		}
		s.Unlock()
		log.Printf("[INFO] compact storage done")
	}
}

// установить значение по ключу с default TTL
func (p *storage) Set(key, val string, ttl int64) {
	p.Lock()
	defer p.Unlock()

	p.data[key] = &storageItem{
		value:      val,
		created_at: time.Now().Unix(),
		ttl:        ttl,
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

	return item.value, true
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
			result = append(result, fmt.Sprintf("%s\t\t%s\t\t%d", key, item.value, item.created_at))
		}
	}
	return result
}
