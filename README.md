## Promise.go

Like them or not, you can't really argue that Promises have not been a useful and succesful strategy in quite a few contexts.

### Just like mom used to make them

```golang
go Promise.Resolve(2).
    then(func(v interface{}) interface{} {
        return v.(int) + 1
    }).
    then(func(v interface{}) interface{} {
        return v.(int) * 7
    }).
    then(func(v interface{}) interface{} {
        log.Print("Answer: ", v)
        return nil
    })
```


### Error handling - just like the real thing!

```golang
go Promise.Reject(errors.New("This is an error")).
    catch(func(e error) interface{} {
        log.Print("There was an error: ", e)
        return "We recovered from the error"
    }).
    then(func(v interface{}) interface{} {
        log.Print("Finally: ", v)
        return nil
    })
```

### Make sure you don't terminate the main thread though...

```golang
go Promise.Resolve(2).
    then(func(v interface{}) interface{} {
        time.sleep(10);
        log.Print("never gonna happen");
    })
go Promise.Resolve(2).
    then(func(v interface{}) interface{} {
        time.sleep(10);
        log.Print("main thread will exit");
    })
```