package discovery

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/allegro/bigcache"
	"github.com/hashicorp/consul/api"
)

type Discovery interface {
	GetService(name string) ([]*api.ServiceEntry, error)
	GetServiceAddress(name string) (string, error)
}

type ConsulDiscovery struct {
	client *api.Client
	cache  *bigcache.BigCache
}

func NewConsul(addr string) (*ConsulDiscovery, error) {
	config := api.DefaultConfig()
	config.Address = addr

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	cacheConfig := bigcache.DefaultConfig(30 * time.Second)
	cacheConfig.MaxEntrySize = 2048   // cukup besar untuk JSON service
	cacheConfig.HardMaxCacheSize = 64 // MB, sesuaikan kebutuhan
	cache, _ := bigcache.NewBigCache(cacheConfig)

	return &ConsulDiscovery{
		client: client,
		cache:  cache,
	}, nil
}

func (c *ConsulDiscovery) GetService(name string) ([]*api.ServiceEntry, error) {
	cacheKey := "consul:" + name

	// 1️⃣ Coba dari cache
	if cached, err := c.cache.Get(cacheKey); err == nil {
		var services []*api.ServiceEntry
		if err := json.Unmarshal(cached, &services); err == nil {
			log.Printf("🔁 BigCache hit for service: %s", name)
			return services, nil
		}
	}

	// 2️⃣ Fallback ke Consul
	services, _, err := c.client.Health().Service(name, "", true, nil)
	if err != nil || len(services) == 0 {
		log.Printf("❌ Consul lookup failed or no healthy instance: %s", name)
		return nil, fmt.Errorf("no healthy instance for %s", name)
	}

	// 3️⃣ Cache kembali ke BigCache
	data, _ := json.Marshal(services)
	_ = c.cache.Set(cacheKey, data)

	log.Printf("✅ Service %s resolved from Consul and cached", name)
	return services, nil
}

func (c *ConsulDiscovery) GetServiceAddress(name string) (string, error) {
	services, err := c.GetService(name)
	if err != nil || len(services) == 0 {
		return "", fmt.Errorf("no available instance for service: %s", name)
	}
	selected := services[time.Now().UnixNano()%int64(len(services))].Service
	return fmt.Sprintf("%s:%d", selected.Address, selected.Port), nil
}
