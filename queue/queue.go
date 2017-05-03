package queue

import (
	"crypto/rand"
	"errors"
	"sync"
	"time"
)

var ErrorHashNotFound = errors.New("Hash not found")
var ErrorItemNotReserved = errors.New("Item not reserved")
var ErrorNoItemsAvailable = errors.New("No items available")

// Item struct is the basic queue item
type Item struct {
	ReservedAt time.Time   `json:"reserved_at"`
	Value      interface{} `json:"value"`
}

// Data contains all data in the current queue
type Data struct {
	MaxReserveTime float64         `json:"max_reserve_time"`
	ItemList       map[string]Item `json:"item"`
	mutex          *sync.RWMutex
}

// New queue
func New() (q *Data, err error) {
	q = &Data{
		MaxReserveTime: 30.0,
		mutex:          &sync.RWMutex{},
		ItemList:       make(map[string]Item),
	}
	return
}

// Count items in the queue
func (q *Data) Count() int {
	q.mutex.RLock()
	defer q.mutex.RUnlock()
	return len(q.ItemList)
}

// Put data in the queue
func (q *Data) Put(v interface{}) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.ItemList[randStr()] = Item{Value: v}
}

// Reserve searches for the next available item in the queue
// If the item is not removed or the reservation time is not
// renewed, the item will returns to the queue automatically
func (q *Data) Reserve() (hash string, value interface{}, err error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for k, v := range q.ItemList {
		now := time.Now()
		diff := now.Sub(v.ReservedAt)
		if diff.Seconds() > q.MaxReserveTime {
			v.ReservedAt = now
			q.ItemList[k] = v
			value = v.Value
			hash = k
			return
		}
	}

	err = ErrorNoItemsAvailable
	return
}

// Renew the reservation of an item in the queue
func (q *Data) Renew(hash string) (err error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	var item Item
	item, err = q.getIten(hash)
	if err != nil {
		return
	}

	item.ReservedAt = time.Now()
	q.ItemList[hash] = item
	return
}

// Remove item from the queue, the item must be reserved
func (q *Data) Remove(hash string) (err error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	_, err = q.getIten(hash)
	if err != nil {
		return
	}

	delete(q.ItemList, hash)
	return
}

// Release reserved item
func (q *Data) Release(hash string) (err error) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	var item Item
	item, err = q.getIten(hash)
	if err != nil {
		return
	}

	item.ReservedAt = time.Time{}
	q.ItemList[hash] = item

	return
}

func (q *Data) getIten(hash string) (item Item, err error) {
	v, ok := q.ItemList[hash]
	if !ok {
		err = ErrorHashNotFound
		return
	}
	diff := time.Since(v.ReservedAt)
	if diff.Seconds() >= q.MaxReserveTime {
		err = ErrorItemNotReserved
		return
	}
	item = v
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
