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
		if err == ERROR_NO_ITEMS_AVALIABLE && i == 3 {
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
	if err != ERROR_HASH_NOT_FOUND {
		t.Fatal("Expected error ERROR_HASH_NOT_FOUND, " + err.Error())
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
	if err != ERROR_NOT_RESERVED {
		t.Fatal("Expected error ERROR_NOT_RESERVED, " + err.Error())
	}

	err = q.Remove(hash)
	if err != ERROR_NOT_RESERVED {
		t.Fatal("Expected error ERROR_NOT_RESERVED, " + err.Error())
	}

	err = q.Renew("fake hash")
	if err != ERROR_HASH_NOT_FOUND {
		t.Fatal("Expected error ERROR_HASH_NOT_FOUND, " + err.Error())
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
	if err != ERROR_HASH_NOT_FOUND {
		t.Fatal("Expected error ERROR_HASH_NOT_FOUND, " + err.Error())
	}

	err = q.Release(hash)
	if err != nil {
		t.Fatal(err.Error())
	}

	err = q.Release(hash)
	if err != ERROR_NOT_RESERVED {
		t.Fatal("Expected error ERROR_NOT_RESERVED, " + err.Error())
	}
}
