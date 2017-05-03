package queue

import "testing"

func TestQueue(t *testing.T) {
	q, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}

	b := []byte{'a', 'b', 'c'}
	q.Put(b)
	q.Put(b)
	q.Put(b)
	if q.Count() != 3 {
		t.Fatal("The queue should contain three items")
	}

	for i := 0; i < 4; i++ {
		var r interface{}
		var hash string
		hash, r, err = q.Reserve()
		if err == ErrorNoItemsAvailable && i == 3 {
			break
		}

		if err != nil {
			t.Fatal(err.Error())
		}

		x := r.([]byte)

		if len(x) != 3 || x[0] != 'a' || x[1] != 'b' || x[2] != 'c' {
			t.Fatalf("Corrupted Payload, len:%d, %#v", len(x), x)
		}

		err = q.Remove(hash)
		if err != nil {
			t.Fatal(err.Error())
		}
	}

	err = q.Remove("fake hash")
	if err != ErrorHashNotFound {
		t.Fatal("Expected error ErrorHashNotFound, " + err.Error())
	}

}

func TestRenew(t *testing.T) {
	q, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}

	b := []byte{'a', 'b', 'c'}
	q.Put(b)
	if q.Count() != 1 {
		t.Fatal("The queue should contain one item")
	}

	var r interface{}
	var hash string
	hash, r, err = q.Reserve()

	x := r.([]byte)

	if len(x) != 3 || x[0] != 'a' || x[1] != 'b' || x[2] != 'c' {
		t.Fatalf("Corrupted Payload, len:%d, %#v", len(x), x)
	}

	err = q.Renew(hash)
	if err != nil {
		t.Fatal(err.Error())
	}

	q.MaxReserveTime = 0

	err = q.Renew(hash)
	if err != ErrorItemNotReserved {
		t.Fatal("Expected error ErrorItemNotReserved, " + err.Error())
	}

	err = q.Remove(hash)
	if err != ErrorItemNotReserved {
		t.Fatal("Expected error ErrorItemNotReserved, " + err.Error())
	}

	err = q.Renew("fake hash")
	if err != ErrorHashNotFound {
		t.Fatal("Expected error ErrorHashNotFound, " + err.Error())
	}

}

func TestRelease(t *testing.T) {
	q, err := New()
	if err != nil {
		t.Fatal(err.Error())
	}

	b := []byte{'a', 'b', 'c'}
	q.Put(b)
	if q.Count() != 1 {
		t.Fatal("The queue should contain one item")
	}

	var r interface{}
	var hash string
	hash, r, err = q.Reserve()

	x := r.([]byte)

	if len(x) != 3 || x[0] != 'a' || x[1] != 'b' || x[2] != 'c' {
		t.Fatalf("Corrupted Payload, len:%d, %#v", len(x), x)
	}

	err = q.Release("fake hash")
	if err != ErrorHashNotFound {
		t.Fatal("Expected error ErrorHashNotFound, " + err.Error())
	}

	err = q.Release(hash)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = q.Release(hash)
	if err != ErrorItemNotReserved {
		t.Fatal("Expected error ErrorItemNotReserved, " + err.Error())
	}
}
