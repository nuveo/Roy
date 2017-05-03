# Queue

That is a mechanism for storing data temporarily and making it available to be consumed asynchronously by other parts of the system.

## Characteristics

This implementation works only in memory, there is no automatic persistence if the memory limit is reached.
The safest way to use the queue system is to persist the data in a database using a transaction and in the same way update that database in a transaction whenever a queue item is processed. In case there is an error, undo the change in the database.

## Basic Operation

### Put

Places an item in the queue.

#### Example

```go
q, err := New()
if err != nil {
	fmt.Println(err.Error())
    return
}

b := []byte{'a', 'b', 'c'}
q.Put(b)
```

### Reserve

Reserve picks up an item from the queue and makes it available for processing. The system has 30 seconds or the time in MaxReserveTime to process the item, this time can be extended using the Renew function. If the item is not deleted or released until the time runs out the system automatically returns the item to the list.

#### Example

```go
q, err := New()
if err != nil {
	fmt.Println(err.Error())
    return
}

b := []byte{'a', 'b', 'c'}
q.Put(b)

var r interface{}
_, r, err = q.Reserve()

x := r.([]byte)

fmt.Println(string(x))
```

### Renew

The Renew function reset the timer of the reserved item so that the system has more time to process the data. In long processes call Renew periodically.

#### Example

```go
q, err := New()
if err != nil {
    fmt.Println(err.Error())
    return
}

b := []byte{'a', 'b', 'c'}
q.Put(b)

var r interface{}
var hash string
hash, r, err = q.Reserve()
if err != nil {
    fmt.Println(err.Error())
    return
}

...
	
err = q.Renew(hash)
if err != nil {
    fmt.Println(err.Error())
    return
}
```

### Release

The Release function frees the reserved item by returning it to the queue and leaving it available to by used by another process. Must be used when the current instance can not process the reserved item.

#### Example

```go
q, err := New()
if err != nil {
    fmt.Println(err.Error())
    return
}

b := []byte{'a', 'b', 'c'}
q.Put(b)

var r interface{}
var hash string
hash, r, err = q.Reserve()
if err != nil {
    fmt.Println(err.Error())
    return
}

...
	
err = q.Release(hash)
if err != nil {
    fmt.Println(err.Error())
    return
}
```

### Remove

The Remove function is used to remove the reserved item from the list. Must be used when the current instance was able to process the data and it can be removed.

#### Example

```go
q, err := New()
if err != nil {
    fmt.Println(err.Error())
    return
}

b := []byte{'a', 'b', 'c'}
q.Put(b)

var r interface{}
var hash string
hash, r, err = q.Reserve()
if err != nil {
    fmt.Println(err.Error())
    return
}

...
	
err = q.Remove(hash)
if err != nil {
    fmt.Println(err.Error())
    return
}
```

### Count 

The Count function is used to know how many items still exist in the list.

#### Example

```go
q, err := New()
if err != nil {
    fmt.Println(err.Error())
    return
}

b := []byte{'a', 'b', 'c'}
q.Put(b)

fmt.Printf("Count: %d",q.Count())
```
