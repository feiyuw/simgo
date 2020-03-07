package client

import (
	"errors"
	"sort"
	"sync"
	"sync/atomic"

	"simgo/protocols"
)

var (
	clientStorage = &clients{M: map[uint64]*Client{}}
	nextClientID  uint64
)

type Client struct {
	Id        uint64                 `json:"id"`
	Protocol  string                 `json:"protocol"`
	Server    string                 `json:"server"`
	Options   map[string]interface{} `json:"options"`
	RpcClient protocols.RpcClient
}

type clients struct {
	sync.RWMutex

	M map[uint64]*Client
}

func (c *clients) Add(value *Client) (uint64, error) {
	c.Lock()
	defer c.Unlock()
	value.Id = atomic.AddUint64(&nextClientID, 1)
	c.M[value.Id] = value
	return value.Id, nil
}

func (c *clients) Remove(key uint64) error {
	c.Lock()
	defer c.Unlock()

	client, exists := c.M[key]
	if exists {
		if client.RpcClient != nil {
			if err := client.RpcClient.Close(); err != nil { // NOTE: Close may cause some time, make c.M hung
				return err
			}
		}
		delete(c.M, key)
	}
	return nil
}

func (c *clients) FindAll() ([]*Client, error) {
	c.RLock()
	defer c.RUnlock()
	items := make([]*Client, len(c.M))
	idx := 0
	for _, v := range c.M {
		items[idx] = v
		idx++
	}
	// sort clients with ID order
	sort.Slice(items, func(idx1, idx2 int) bool {
		return items[idx1].Id < items[idx2].Id
	})

	return items, nil
}

func (c *clients) FindOne(key uint64) (*Client, error) {
	c.RLock()
	defer c.RUnlock()
	if v, exists := c.M[key]; exists {
		return v, nil
	}
	return nil, errors.New("not found")
}
