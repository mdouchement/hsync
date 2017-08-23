package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	resty "gopkg.in/resty.v0"
)

const (
	endpoint = "http://localhost:5005"
	lockID   = "tk-007"
)

func main() {
	unlock(lockID)

	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)

		go func(user string) {
			defer wg.Done()
			polling(user, lockID)
		}(fmt.Sprintf("user_%d", i))
	}

	wg.Wait()
}

func polling(name, id string) {
	fmt.Printf("(%s) init %s (%s)\n", name, id, time.Now().Format("15:04:05.000"))

	for {
		if lock(id) {
			fmt.Printf("(%s) lock %s acquired (%s)\n", name, id, time.Now().Format("15:04:05.000"))
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	time.Sleep(time.Duration(rand.Intn(10)+1) * time.Second)

	unlock(id)
	fmt.Printf("(%s) lock %s released (%s)\n", name, id, time.Now().Format("15:04:05.000"))
}

func lock(id string) bool {
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

func unlock(id string) bool {
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
