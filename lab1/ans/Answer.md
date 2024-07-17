## Answer
- What is atomic operation?

An atomic operation is a low-level operation that is completed in a single step relative to other threads or processes. Without locks, there are only simple CPU instructions that cannot be further subdivided into smaller steps.
- Does the example below have concurrency issues? Why?
```go
var a int64

func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    
    go func() {
        increment()
        wg.Done()
    }()

    wg.Wait()

    fmt.Println("Final value of a:", a)
}

func increment() {
    for i := 0; i < 100000000; i++ {
        go func() {
             a = a + 1    
        }()
    }
}
```
Yes, In the original code, **`a = a + 1`** is not an atomic operation. This operation is divided into multiple steps: reading the value of **`a`**, adding 1, and then writing it back to **`a`**. The **`go func()`** inside the loop creates multiple goroutines that execute these steps simultaneously, interfering with each other and leading to unpredictable results. Now we rewrote this code with an atomic operation to improve concurrency issues:
```go
//goog code
var a int64

func main() {
    var wg sync.WaitGroup
    wg.Add(1)

    go func() {
        increment()
        wg.Done()
    }()

    wg.Wait()

    fmt.Println("Final value of a:", a)
}

func increment() {
    for i := 0; i < 100000000; i++ {
        go atomic.AddInt64(&a, 1)
    }
}
```
The **`atomic.AddInt64`** ensures that the increment operation on the **`a`** is performed atomically. This means that the increment operation cannot be interrupted, ensuring the correctness and expected result even when it's accessed by multiple goroutines.
- What do the **`wg.Add(1)`** and **`wg.Done()`** do in the above statement? And what does the **`1`** repersent?
    * **`wg.Add(1)`** and **`wg.Done()`** are used in conjunction with a sync.WaitGroup to synchronize the execution of the main goroutine with the completion of the incrementing goroutine
    * The argument **`1`** indicates that 1 goroutine is being added to the count of goroutines that the WaitGroup will wait for.
- Please define a Counter struct with an integer field and a sync.Mutex, then implement a function to increment the counter safely.
```go
// Counter struct with a mutex and a value
type Counter struct {
	mu    sync.Mutex
	value int
}

func (c *Counter) Increment() {
	c.mu.Lock()   // Lock the mutex before modifying the value
	defer c.mu.Unlock() // Ensure the mutex is unlocked after the function completes
	c.value++
}
```
- Why is there a fetal error in following code?
```go=
func main() {
    var intChan chan int
    fmt.Println(intChan)
    intChan <- 10
}
```
Without using **`make`** to allocate space, the channel is used directly. It will cause fetal error: all goroutines are asleep - deadlock!

- In Topic: **Select**-**Setup Timeout mechanism**, We offered a sample code that showed how to setup **`timeout`** with **`select`**. Actually, there is a potentail error (Hint: memory leak) because of **`time.After`** usage. Please describe the reason for the error and how to fix it.

The previous sample code avoids memory leaks because **`time.After`** executes first. This ensures there's no resource buildup. However, if the **`go fun`**'s waiting time is shorter than **`time.After`**'s trigger time, **`case <-timeout`** would execute first. This could prevent **`time.After`** from correctly returning a value, potentially causing a memory leak. Here are possible solutions:
    * Pull **`time.After`** outside of **`select`** and then assign its return value to a **`variable`** before using it in select
    * Alternatively, use **`time.NewTimer`** and **`timer.Stop()`** instead of **`time.After()`** to more flexibly manage timer stops and memory release, thereby avoiding potential memory leaks associated with **`time.After()`**:
    
```go
    func main() {
    timeout := make(chan bool, 1)
    go func() {
        time.Sleep(1 * time.Second)
        timeout <- true
    }()
    ch := make(chan int)

    timer := time.NewTimer(time.Second * 2)

    select {
    case <-ch:
    case <-timeout:
        fmt.Println("Open5GS")
    case <-timer.C:
        fmt.Println("free5GC")
    }

    if !timer.Stop() {
        <-timer.C 
        }
    }

```
* reference: [Memory Leak in Go](https://arangodb.com/2020/09/a-story-of-a-memory-leak-in-go-how-to-properly-use-time-after/)