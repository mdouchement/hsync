package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/phayes/freeport"

	resty "gopkg.in/resty.v0"
)

var endpoint string

func init() {
	socket := fmt.Sprintf("localhost:%d", freeport.GetPort())
	endpoint = fmt.Sprintf("http://%s", socket)

	go server(socket)
}

func TestHSync(t *testing.T) {
	var wg sync.WaitGroup
	var m = map[string]*metric{}

	for i := 0; i < 3; i++ {
		wg.Add(1)

		go func(user string) {
			defer wg.Done()
			polling(user, "my-lock", m)
		}(fmt.Sprintf("user_%d", i))
	}

	wg.Wait()

	for k1, v1 := range m {
		for k2, v2 := range m {
			if k1 == k2 {
				continue // Same entry
			}

			if v1.LockedAt.Before(v2.LockedAt) {
				if !v1.UnlockedAt.Before(v2.LockedAt) {
					t.Errorf("%s (%v) is locked before %s (%v) is unlocked", k2, v2.LockedAt, k1, v1.UnlockedAt)
				}
			} else {
				if !v2.UnlockedAt.Before(v1.LockedAt) {
					t.Errorf("%s (%v) is locked before %s (%v) is unlocked", k1, v1.UnlockedAt, k2, v2.LockedAt)
				}
			}
		}
	}
}

//===================================================//
//                                                   //
//  Helpers                                          //
//                                                   //
//===================================================//

type (
	metric struct {
		User       string
		LockedAt   time.Time
		UnlockedAt time.Time
	}
)

func polling(user, id string, m map[string]*metric) {
	fmt.Printf("(%s) init %s (%s)\n", user, id, time.Now().Format("15:04:05.000"))

	for {
		if acquire(id) {
			m[user] = &metric{User: user, LockedAt: time.Now()}
			fmt.Printf("(%s) lock %s acquired (%s)\n", user, id, time.Now().Format("15:04:05.000"))
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	time.Sleep(time.Duration(rand.Intn(4)+1) * time.Second)

	m[user].UnlockedAt = time.Now()
	release(id)
	fmt.Printf("(%s) lock %s released (%s)\n", user, id, time.Now().Format("15:04:05.000"))
}

func acquire(id string) bool {
	r, err := resty.R().
		SetBody(map[string]interface{}{"id": id}).
		Post(fmt.Sprintf("%s/locks", endpoint))

	check(err)

	switch r.RawResponse.StatusCode {
	case http.StatusCreated:
		return true
	case http.StatusLocked:
		return false
	default:
		panic(fmt.Errorf("Unsupported HTTP status %s", r.RawResponse.Status))
	}
}

func release(id string) bool {
	r, err := resty.R().
		Delete(fmt.Sprintf("%s/locks/%s", endpoint, id))

	check(err)

	switch r.RawResponse.StatusCode {
	case http.StatusNoContent:
		return true
	case http.StatusNotFound:
		// panic(fmt.Errorf("Unwanted HTTP status %s", r.RawResponse.Status))
		return true
	default:
		panic(fmt.Errorf("Unsupported HTTP status %s", r.RawResponse.Status))
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
