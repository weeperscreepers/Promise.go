package promise_test

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/weeperscreepers/promise"
)

func TestThen(t *testing.T) {
	ch := make(chan int)
	go promise.Resolve(2).
		Then(func(v interface{}) interface{} {
			return v.(int) + 1
		}).
		Then(func(v interface{}) interface{} {
			return v.(int) * 7
		}).
		Then(func(v interface{}) interface{} {
			ch <- v.(int)
			return nil
		})
	i := <-ch
	assert.Equal(t, i, 21)
}
func TestCatch(t *testing.T) {

	EXPECTED := "We recovered from the error"

	ch := make(chan string)
	go promise.Reject(errors.New("This is an error")).
		Catch(func(e error) interface{} {
			log.Print("There was an expected error: ", e)
			return EXPECTED
		}).
		Then(func(v interface{}) interface{} {
			ch <- v.(string)
			return nil
		})
	i := <-ch
	assert.Equal(t, i, EXPECTED)
}

func TestNewResolve(t *testing.T) {
	ch := make(chan int)
	go promise.New(
		func(res promise.Resolver, rej promise.Rejecter) {
			log.Print("Fetching from the 'database' ;)")
			time.Sleep(3) // do some async stuff
			res(200)
		},
	).
		Then(func(v interface{}) interface{} {
			ch <- v.(int)
			return nil
		})
	i := <-ch
	assert.Equal(t, i, 200)
}

func TestNewReject(t *testing.T) {
	ch := make(chan int)
	go promise.New(
		func(res promise.Resolver, rej promise.Rejecter) {
			log.Print("Fetching from the 'database' ;)")
			time.Sleep(3) // do some async stuff
			// res(200)
			rej(errors.New("Error 400"))
		},
	).
		Then(func(v interface{}) interface{} {
			log.Print("we thened")
			// t.Error("This should never be executed")
			ch <- v.(int)
			return nil
		}).
		Catch(func(e error) interface{} {
			log.Print("There was an expected error: ", e)
			ch <- 400
			return nil
		})
	i := <-ch
	assert.Equal(t, i, 400)
}
