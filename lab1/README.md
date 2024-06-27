# Lab 1: Concurrent Programming in Go

## Introduction

In Lab 1, you will learn how to build concurrent program with Go.

## Goals of this lab

- Understand concurrent programming
- Understand the memory models in Go
- Learn How to use synchronization primitives in Go

## Race Condition & Critical Section

Race Condition is a situation where two or more goroutines access a shared resource concurrently, and at least one of the goroutines modifies the resource. This can lead to unexpected behavior, such as data corruption or deadlock.

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var counter int
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            counter++
            wg.Done()
        }()
    }

    wg.Wait()
    fmt.Println(counter)
}
```

Critical Section is a section of code that accesses a shared resource and must be executed by only one goroutine at a time.

The Critical Section in the example listed above is:
```go
counter++
```

Because the `counter++` operation is not atomic or protected by a lock, multiple goroutines can access the shared resource concurrently.
Concurrent access will lead to the race condition.

## Goroutines

Goroutines are lightweight threads of execution that are managed by the Go runtime. Goroutines are used to perform concurrent operations in Go.

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    go func() {
        fmt.Println("Hello, World!")
    }()

    time.Sleep(1 * time.Second)
}
```

## Atomic Operations

In Go, atomic operations are operations that are guaranteed to be executed as a single operation without interruption. This is important in concurrent programming because it ensures that the operation is executed in a consistent state.

In fact, atomic operations are achieved by atomic instructions, provided by the hardware. These instructions are used to perform operations on shared memory in a way that is guaranteed to be atomic.

- [[University of Washington CSE378] Atomic Operations](https://courses.cs.washington.edu/courses/cse378/07au/lectures/L25-Atomic-Operations.pdf)

```go
package main

import (
    "fmt"
    "sync/atomic"
)

func main() {
    var counter int32
    atomic.AddInt32(&counter, 1)
    fmt.Println(counter)
}
```

The example above demonstrates how to use atomic operations in Go. The `AddInt32` function is an atomic operation that increments the value of the counter by 1.

## SpinLock

Spinlock is a synchronization primitive that is used to protect shared resources from concurrent access. A spinlock is used to ensure that only one goroutine can access the shared resource at a time.

Golang doesn't provide a built-in spinlock, but you can implement it using atomic operations.

```go
package main

import (
    "fmt"
    "sync/atomic"
)

type SpinLock struct {
    flag int32
}

func (s *SpinLock) Lock() {
    for !atomic.CompareAndSwapInt32(&s.flag, 0, 1) {
        // Spin until the lock is acquired
    }
}

func (s *SpinLock) Unlock() {
    atomic.StoreInt32(&s.flag, 0)
}
```

## Mutex

In Go, a mutex is a synchronization primitive that is used to protect shared resources from concurrent access. A mutex is used to ensure that only one goroutine can access the shared resource at a time.

The difference between a mutex and a spinlock is that a mutex will put the goroutine to sleep if the resource is already locked, while a spinlock will keep the goroutine busy until the resource is available.

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var counter int
    var mu sync.Mutex

    mu.Lock()
    counter++
    mu.Unlock()

    fmt.Println(counter)
}
```

Any other goroutine that tries to access the shared resource while the mutex is locked will be blocked until the mutex is unlocked.

- [Source Code of Mutex in Go](https://go.dev/src/sync/mutex.go)