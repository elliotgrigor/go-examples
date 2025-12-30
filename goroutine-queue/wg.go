package main

import "fmt"

type CacheRecord struct {
	Key      string
	Value    string
	Response chan bool
}

type Cache struct {
	Store []CacheRecord
	Queue chan CacheRecord
}

var cache = Cache{
	Store: []CacheRecord{},
	Queue: make(chan CacheRecord, 10),
}

func (c *Cache) Worker() {
	for record := range c.Queue {
		c.Store = append(c.Store, CacheRecord{
			Key:   record.Key,
			Value: record.Value,
		})
		// Feed the channel
		record.Response <- true
		fmt.Println(record.Key, record.Value)
	}
}

func (c *Cache) Put(key, value string) {
	response := make(chan bool)
	c.Queue <- CacheRecord{
		Key:      key,
		Value:    value,
		Response: response,
	}
	<-response // Wait for response to be fed into the channel
}
