package queue

import "sync"

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
