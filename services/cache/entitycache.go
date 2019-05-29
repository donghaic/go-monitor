package cache

import (
	"fmt"
	"gt-monitor/common/redis"
	"gt-monitor/models"
	"gt-monitor/services/cache/internal"
)

type EntityCacheService struct {
	localCache *internal.LocalCacheService
	redisCache *internal.RedisCacheService
}

func NewEntityCacheService(redisPool *redis.ConnPool) *EntityCacheService {
	return &EntityCacheService{
		localCache: internal.NewLocalCache(),
		redisCache: internal.NewRedisCacheService(redisPool),
	}
}

func (c EntityCacheService) GetCreativeById(creativeId string) (*models.CreativeInfo, error) {
	creative, err := c.localCache.GetCreativeById(creativeId)
	if err != nil {
		creative, err := c.redisCache.GetCreativeById(creativeId)
		if err == nil {
			c.localCache.SetCreative(creative)
		}
		return creative, err
	} else {
		return creative, err
	}
}

func (c EntityCacheService) GetLinkById(linkId string) (*models.LinkSwitch, error) {
	link, err := c.localCache.GetLinkById(linkId)
	if err != nil {
		link, err := c.redisCache.GetLinkById(fmt.Sprintf("link:%s", linkId))
		if err == nil {
			c.localCache.SetLink(link)
		}
	}
	return link, err
}

func (c EntityCacheService) FlushLocal() {
	c.localCache.Flush()
}
