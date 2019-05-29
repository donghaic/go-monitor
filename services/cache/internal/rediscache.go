package internal

import (
	"encoding/json"
	"gt-monitor/common/redis"
	"gt-monitor/models"
)

type RedisCacheService struct {
	pool *redis.ConnPool
}

func NewRedisCacheService(connPool *redis.ConnPool) *RedisCacheService {
	return &RedisCacheService{
		pool: connPool,
	}
}

func (r RedisCacheService) GetCreativeById(creativeId string) (*models.CreativeInfo, error) {
	jsonData, err := r.pool.GetString(creativeId)
	if err == nil {
		info := models.CreativeInfo{}
		var er = json.Unmarshal([]byte(jsonData), &info)
		if nil == er {
			return &info, nil
		}
	}
	return nil, err
}

func (r RedisCacheService) GetLinkById(linkId string) (*models.LinkSwitch, error) {
	jsonData, err := r.pool.GetString(linkId)
	if err == nil {
		info := models.LinkSwitch{}
		err = json.Unmarshal([]byte(jsonData), &info)
		if nil == err {
			return &info, nil
		}
	}
	return nil, err
}
