package queue

import (
	"crypto/rand"
	"sync"
)

// Item struct is the basic queue item
type Item struct {
	Tag   string `json:"tag"`
	Value []byte `json:"value"`
}

// Data contains all data in the current queue
type Data struct {
	Tag      string          `json:"tag"`
	ItemList map[string]Item `json:"item"`
	control  queueControl
}

type queueControl struct {
	mutex *sync.Mutex
}

// New queue
func New() (q *Data, err error) {
	q = &Data{
		ItemList: make(map[string]Item),
		control:  queueControl{mutex: &sync.Mutex{}},
	}
	return
}

// Count items in the queue
func (q *Data) Count() int {
	return len(q.ItemList)
}

// Put data in the queue
func (q *Data) Put(b []byte) {
	q.ItemList[randStr()] = Item{Value: b}
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
