package storage

import (
	"encoding/json"
	"time"
)

type StorageItem struct {
	Metric    string            `json:"metric"`
	Value     string            `json:"value"`
	Tags      map[string]string `json:"tags"`
	CreatedAt int64             `json:"created_at"`
	TTL       int64             `json:"ttl"`
}

func (s *StorageItem) GetMetric() string {
	return s.Metric
}

func (s *StorageItem) GetValue() string {
	return s.Value
}

func (s *StorageItem) GetTags() map[string]string {
	return s.Tags
}

func (s *StorageItem) GetCreatedAt() int64 {
	return s.CreatedAt
}

func storageKey(metric string, tags map[string]string) string {
	result := metric
	if len(tags) > 0 {
		data, err := json.Marshal(&tags)
		if err == nil {
			result = result + string(data)
		}
	}
	return result
}

func (s *StorageItem) Key() string {
	return storageKey(s.Metric, s.Tags)
}

func (s *StorageItem) valid() bool {
	return s.CreatedAt+s.TTL > time.Now().Unix()
}
