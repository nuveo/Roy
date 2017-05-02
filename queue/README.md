# Queue

That is a mechanism for storing data temporarily and making it available to be consumed asynchronously by other parts of the system.

## Characteristics

This implementation works only in memory, there is no automatic persistence if the memory limit is reached.
The safest way to use the queue system is to persist the data in a database using a transaction and in the same way update that database in a transaction whenever a queue item is processed. In case there is an error, undo the change in the database.

## Basic Operation

### Put

Places an item in the queue.

###Example

```go
q, err := New()
if err != nil {
	fmt.Println(err.Error())
    return
}

b := []byte{'a', 'b', 'c'}
q.Put(b)
```

reserve

delete

list

count 

size
