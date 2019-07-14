## Promise.go

Like them or not, you can't really argue that Promises have not sometimes been a useful and succesful strategy

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