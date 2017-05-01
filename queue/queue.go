package queue

import (
	"crypto/rand"
	"errors"
	"sync"
	"time"
)

const MaxReserTime = 30

var ERROR_HASH_NOT_FOUND = errors.New("Hash not found")
var ERROR_ALREADY_RESERVED = errors.New("Item already reserved")
var ERROR_NOT_RESERVED = errors.New("Item not reserved")

// Item struct is the basic queue item
type Item struct {
	ReservedAt time.Time `json:"reserved_at"`
	Tag        string    `json:"tag"`
	Value      []byte    `json:"value"`
}

// Data contains all data in the current queue
type Data struct {
	Tag      string          `json:"tag"`
	ItemList map[string]Item `json:"item"`
	sync.Mutex
}

// New queue
func New() (q *Data, err error) {
	q = &Data{
		ItemList: make(map[string]Item),
	}
	return
}

// Count items in the queue
func (q *Data) Count() int {
	return len(q.ItemList)
}

// Put data in the queue
func (q *Data) Put(b []byte) {
	q.Lock()
	defer q.Unlock()
	q.ItemList[randStr()] = Item{Value: b}
}

// Reserve an item to be processed.
// If the item is not removed or the reservation time is not
// renewed, the item will returns to the queue automatically
func (q *Data) Reserve(hash string) (item Item, err error) {
	q.Lock()
	defer q.Unlock()
	v, ok := q.ItemList[hash]
	if !ok {
		err = ERROR_HASH_NOT_FOUND
		return
	}
	now := time.Now()
	diff := now.Sub(v.ReservedAt)
	if diff.Seconds() < MaxReserTime {
		err = ERROR_ALREADY_RESERVED
		return
	}
	v.ReservedAt = now
	item = v
	q.ItemList[hash] = v
	return
}

// Remove item from the queue, the item must be reserved.
func (q *Data) Remove(hash string) (err error) {
	q.Lock()
	defer q.Unlock()
	v, ok := q.ItemList[hash]
	if !ok {
		err = ERROR_HASH_NOT_FOUND
		return
	}
	diff := time.Since(v.ReservedAt)
	if diff.Seconds() >= MaxReserTime {
		err = ERROR_NOT_RESERVED
		return
	}
	delete(q.ItemList, hash)
	return
}

// Renew the reservation of an item in the queue.
func (q *Data) Renew(hash string) (err error) {
	q.Lock()
	defer q.Unlock()
	v, ok := q.ItemList[hash]
	if !ok {
		err = ERROR_HASH_NOT_FOUND
		return
	}
	now := time.Now()
	diff := now.Sub(v.ReservedAt)
	if diff.Seconds() >= MaxReserTime {
		err = ERROR_NOT_RESERVED
		return
	}
	v.ReservedAt = now
	q.ItemList[hash] = v
	return
}

func randStr() (ret string) {
	const charList = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	const clLen byte = 62
	const size = 36
	var bytes = make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}
	for k, v := range bytes {
		bytes[k] = charList[v%clLen]
	}
	ret = string(bytes)
	return
}
