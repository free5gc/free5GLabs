# Lab 1: Concurrent Programming in Go

## Introduction

In Lab 1, you will learn how to build concurrent program with Go.

## Goals of this lab

- Understand concurrent programming
- Understand the memory models in Go
- Learn How to use synchronization primitives in Go

## Atomic Operations

In Go, atomic operations are operations that are guaranteed to be executed as a single operation without interruption. This is important in concurrent programming because it ensures that the operation is executed in a consistent state.

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

## Mutex