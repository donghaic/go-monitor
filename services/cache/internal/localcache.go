package internal

import (
	"gt-monitor/models"
	gocache "github.com/patrickmn/go-cache"
	"time"
)

var (
	defaultExpiration = 5 * time.Minute
	cleanupInterval   = 5 * time.Minute
)

type LocalCacheService struct {
	creative *gocache.Cache
	link     *gocache.Cache
}

func NewLocalCache() *LocalCacheService {
	return &LocalCacheService{
		creative: gocache.New(defaultExpiration, cleanupInterval),
		link:     gocache.New(defaultExpiration, cleanupInterval),
	}
}

func (c *LocalCacheService) SetCreative(creative *models.CreativeInfo) {
	c.creative.SetDefault(creative.CreativeId, creative)
}

func (c *LocalCacheService) GetCreativeById(creativeId string) (*models.CreativeInfo, error) {
	if x, found := c.creative.Get(creativeId); found {
		info := x.(*models.CreativeInfo)
		return info, nil
	}
	return nil, ErrNotFound
}

func (c *LocalCacheService) GetLinkById(linkId string) (*models.LinkSwitch, error) {
	if x, found := c.link.Get(linkId); found {
		info := x.(*models.LinkSwitch)
		return info, nil
	}
	return nil, ErrNotFound
}

func (c *LocalCacheService) SetLink(link *models.LinkSwitch) {
	c.link.SetDefault(link.AfLinkId, link)
}

func (c *LocalCacheService) Flush() {
	c.link.Flush()
	c.creative.Flush()
}
