## Promise.go

Promises are (were?) a staple of Javascript programming, and they get a lot of hate. With the introduction of `async-await`, the age of the promise may be coming to a close.

Go `chan`s behave similarly to `async-await`, but might seem a little confusing to the outsider.

But the nice thing about `chan`s is that they provide a nice primitive to build asynchronous code on top of.

So with that said, here is a library implementing the `Promise` pattern well-known (maligned?) by web devs everywhere.

### Just like mom used to make them

```golang
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
```


### Error handling - just like the real thing!

```golang
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
```

### Real async ! Don't terminate the main thread...

```golang
    go promise.Resolve(2).
        Then(func(v interface{}) interface{} {
            time.sleep(10);
            log.Print("never gonna happen");
        })
    go promise.Resolve(2).
        Then(func(v interface{}) interface{} {
            time.sleep(10);
            log.Print("main thread will exit");
        })
```

### Spread them everywhere !

```golang
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

```

### Quickly log using the standard logging interface

```golang

    logger = log.New(os.Stdout, "LOG: ", 0)

	ch := make(chan int)
	go promise.Resolve(2).
		Then(func(v interface{}) interface{} {
			return v.(int) + 1
		}).
		Then(promise.Log(logger)).
		Then(func(v interface{}) interface{} {
			return v.(int) * 7
		}).
		Then(func(v interface{}) interface{} {
			ch <- v.(int)
			return nil
		})
	i := <-ch
	assert.Equal(t, 21, i)
```

### But wait, there's more !

You might think of promise chains like a type of abstracted callstack.

Promise.go gives you the power to control the execution of Promise chains with the introduction of paused promises.

Paused promises are just like promises, except they don't execute until you tell them to.
`Step()` will execute the next callback in the chain, and `Allocate(n int)` is like calling `Step()` `n` times.

`Continue()` returns the promise to normal behavior.

Here's an example, slightly modified from a test case:

```golang

type intChannel chan int

func (ch intChannel) reportAndIncrement(v interface{}) interface{} {
	ch <- v.(int)
	return v.(int) + 1
}
func (ch intChannel) reportAndDecrement(v interface{}) interface{} {
	ch <- v.(int)
	return v.(int) - 1
}

ch := make(intChannel)

sayChannel0 := promise.Paused(0).
	Then(ch.reportAndIncrement).
	Then(ch.reportAndIncrement).(promise.PausedPromise)

sayChannel1 := promise.Paused(100).
	Then(ch.reportAndDecrement).
	Then(ch.reportAndDecrement).
	Then(ch.reportAndDecrement).(promise.PausedPromise)

sayChannel1.Step()
i := <-ch
assert.Equal(t, 100, i)

sayChannel0.Step()
i = <-ch
assert.Equal(t, 0, i)

sayChannel1.Continue()
i = <-ch
assert.Equal(t, 99, i)
i = <-ch
assert.Equal(t, 98, i)

```