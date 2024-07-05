package sas

import (
	"container/list"
)

var LruCache *Cache

// Cache 代表LRU缓存实现，暂时未考虑线程安全
type Cache struct {
	// MaxEntries 表示缓存容量的最大值，0表示是一个空缓存
	MaxEntries int

	ll    *list.List
	cache map[string]*list.Element
}

// entry 表示一个缓存键值对
type entry struct {
	key   string
	value interface{}
}

// New 函数新建一个LRU缓存对象
func New(max int) *Cache {
	return &Cache{
		MaxEntries: max,
		ll:         list.New(),
		cache:      make(map[string]*list.Element), //value是双向链表的节点的指针
	}
}

// Set 函数添加一个缓存项到Cache对象中
func (c *Cache) Set(key string, value interface{}) {
	if c.cache == nil {
		c.cache = make(map[string]*list.Element)
		c.ll = list.New()
	}
	// 如果缓存已经存在于Cache中，那么将该缓存项移到双向链表的最前端
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		ee.Value.(*entry).value = value
		return
	}

	// 将新添加的缓存项放入双向链表的最前端
	ele := c.ll.PushFront(&entry{key, value})
	c.cache[key] = ele

	// 如果超出缓存容量，那么移除双向链表中的最后一项
	if c.MaxEntries != 0 && c.ll.Len() > c.MaxEntries {
		c.RemoveOldest()
	}
}

// Get 方法获取具有指定键的缓存项
func (c *Cache) Get(key string) (value interface{}, ok bool) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry).value, true
	}
	return
}

// Remove 方法移除具有指定键的缓存
func (c *Cache) Remove(key string) {
	if c.cache == nil {
		return
	}
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// RemoveOldest 移除双向链表中访问时间最远的那一项
func (c *Cache) RemoveOldest() {
	if c.cache == nil {
		return
	}
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *Cache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry)
	delete(c.cache, kv.key)
}

// Len 方法获取Cache中包含的缓存项个数
func (c *Cache) Len() int {
	if c.cache == nil {
		return 0
	}
	return c.ll.Len()
}

// Clear 清除整个Cache对象
func (c *Cache) Clear() {
	c.ll = nil
	c.cache = nil
}
