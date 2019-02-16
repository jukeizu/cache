package cache

import (
	"testing"
)

type KeyTest struct {
	A string
	B string
}

type ValueTest struct {
	C string
	D string
}

func TestSetGet(t *testing.T) {
	cacheConfig := Config{
		Address: DefaultRedisAddress,
		Version: "0.0.1",
	}

	key := KeyTest{"A", "B"}
	value := ValueTest{"C", "D"}

	c := New(cacheConfig)

	err := c.Set(key, value, 0)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	gotValue := ValueTest{}

	err = c.Get(key, &gotValue)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	t.Log(gotValue)

	if gotValue.C != "C" {
		t.Errorf("Should be C")
	}
}
