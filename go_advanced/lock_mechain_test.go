// 编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// 共享计数器结构
type Counter struct {
	value int
	mu    sync.Mutex
}

// 安全的递增方法
func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// 获取当前值
func (c *Counter) GetValue() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

func main() {
	fmt.Println("=== 使用 sync.Mutex 保护共享计数器 ===\n")

	counter := &Counter{}
	var wg sync.WaitGroup

	goroutineCount := 10
	incrementsPerGoroutine := 1000

	fmt.Printf("启动 %d 个协程，每个协程递增 %d 次\n", goroutineCount, incrementsPerGoroutine)
	fmt.Printf("期望最终值: %d\n\n", goroutineCount*incrementsPerGoroutine)

	// 启动多个协程并发递增计数器
	startTime := time.Now()

	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go worker(i, counter, incrementsPerGoroutine, &wg)
	}

	// 等待所有协程完成
	wg.Wait()

	duration := time.Since(startTime)

	// 输出结果
	finalValue := counter.GetValue()
	expectedValue := goroutineCount * incrementsPerGoroutine

	fmt.Printf("\n=== 执行结果 ===\n")
	fmt.Printf("最终计数器值: %d\n", finalValue)
	fmt.Printf("期望计数器值: %d\n", expectedValue)
	fmt.Printf("执行时间: %v\n", duration)

	if finalValue == expectedValue {
		fmt.Println("✅ 结果正确！所有递增操作都成功执行")
	} else {
		fmt.Printf("❌ 结果错误！丢失了 %d 次递增\n", expectedValue-finalValue)
	}
}

// 工作协程：执行指定次数的递增操作
func worker(id int, counter *Counter, times int, wg *sync.WaitGroup) {
	defer wg.Done()

	// 每个协程执行指定次数的递增
	for i := 0; i < times; i++ {
		counter.Increment()
	}

	fmt.Printf("协程 %d 完成 %d 次递增\n", id, times)
}

func main() {
	fmt.Println("=== 使用 sync/atomic 实现无锁计数器 ===\n")

	// 演示1：原子操作计数器
	demonstrateAtomicCounter()

	// 演示2：与互斥锁性能对比
	demonstratePerformanceComparison()

	// 演示3：更多原子操作示例
	demonstrateMoreAtomicOperations()
}

// 演示1：原子操作计数器
func demonstrateAtomicCounter() {
	fmt.Println("1. 原子操作计数器演示:")

	var counter int64 // 必须使用 int64 类型
	var wg sync.WaitGroup
	goroutineCount := 10
	incrementsPerGoroutine := 1000

	wg.Add(goroutineCount)

	start := time.Now()

	for i := 0; i < goroutineCount; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				// 使用原子操作递增计数器
				atomic.AddInt64(&counter, 1)
			}
			fmt.Printf("  协程 %d 完成 %d 次递增\n", id, incrementsPerGoroutine)
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	expected := int64(goroutineCount * incrementsPerGoroutine)
	finalValue := atomic.LoadInt64(&counter)

	fmt.Printf("  最终计数器: %d (期望值: %d)\n", finalValue, expected)
	fmt.Printf("  总耗时: %v\n", elapsed)

	if finalValue == expected {
		fmt.Printf("  ✅ 原子操作正确！\n")
	} else {
		fmt.Printf("  ❌ 操作错误！差值: %d\n", expected-finalValue)
	}
	fmt.Println()
}

// 演示2：与互斥锁性能对比
func demonstratePerformanceComparison() {
	fmt.Println("2. 原子操作 vs 互斥锁性能对比:")

	iterations := 1000000 // 每个协程的操作次数
	goroutineCount := 10

	// 测试原子操作
	var atomicCounter int64
	var wg1 sync.WaitGroup
	wg1.Add(goroutineCount)

	atomicStart := time.Now()
	for i := 0; i < goroutineCount; i++ {
		go func() {
			defer wg1.Done()
			for j := 0; j < iterations; j++ {
				atomic.AddInt64(&atomicCounter, 1)
			}
		}()
	}
	wg1.Wait()
	atomicTime := time.Since(atomicStart)

	// 测试互斥锁
	var mutexCounter int64
	var mutex sync.Mutex
	var wg2 sync.WaitGroup
	wg2.Add(goroutineCount)

	mutexStart := time.Now()
	for i := 0; i < goroutineCount; i++ {
		go func() {
			defer wg2.Done()
			for j := 0; j < iterations; j++ {
				mutex.Lock()
				mutexCounter++
				mutex.Unlock()
			}
		}()
	}
	wg2.Wait()
	mutexTime := time.Since(mutexStart)

	fmt.Printf("  原子操作 - 计数器: %d, 耗时: %v\n", atomic.LoadInt64(&atomicCounter), atomicTime)
	fmt.Printf("  互斥锁  - 计数器: %d, 耗时: %v\n", mutexCounter, mutexTime)
	fmt.Printf("  性能提升: %.2fx\n", float64(mutexTime)/float64(atomicTime))
	fmt.Println()
}

// 演示3：更多原子操作示例
func demonstrateMoreAtomicOperations() {
	fmt.Println("3. 更多原子操作示例:")

	// 1. 加载和存储操作
	var value int64 = 100
	oldValue := atomic.SwapInt64(&value, 200)
	fmt.Printf("  Swap: 旧值=%d, 新值=%d\n", oldValue, atomic.LoadInt64(&value))

	// 2. 比较并交换 (CAS)
	var casValue int64 = 300
	success := atomic.CompareAndSwapInt64(&casValue, 300, 400)
	fmt.Printf("  CAS(300->400): 成功=%t, 当前值=%d\n", success, atomic.LoadInt64(&casValue))

	success = atomic.CompareAndSwapInt64(&casValue, 300, 500)
	fmt.Printf("  CAS(300->500): 成功=%t, 当前值=%d\n", success, atomic.LoadInt64(&casValue))

	// 3. 原子计数器示例
	demonstrateAtomicCounterAdvanced()
}

// 高级原子计数器示例
func demonstrateAtomicCounterAdvanced() {
	fmt.Println("\n4. 高级原子计数器示例:")

	type AtomicCounter struct {
		value int64
	}

	counter := &AtomicCounter{}
	var wg sync.WaitGroup
	goroutineCount := 8
	operationsPerGoroutine := 500

	wg.Add(goroutineCount * 2) // 递增和递减协程

	// 递增协程
	for i := 0; i < goroutineCount; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				atomic.AddInt64(&counter.value, 1)
				time.Sleep(time.Microsecond)
			}
			fmt.Printf("  递增协程 %d 完成\n", id)
		}(i)
	}

	// 递减协程（稍后启动）
	time.Sleep(2 * time.Millisecond)
	for i := 0; i < goroutineCount; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operationsPerGoroutine; j++ {
				atomic.AddInt64(&counter.value, -1)
				time.Sleep(time.Microsecond)
			}
			fmt.Printf("  递减协程 %d 完成\n", id)
		}(i)
	}

	// 监控协程
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(100 * time.Millisecond)
			current := atomic.LoadInt64(&counter.value)
			fmt.Printf("  监控: 当前计数器值 = %d\n", current)
		}
	}()

	wg.Wait()

	finalValue := atomic.LoadInt64(&counter.value)
	fmt.Printf("  最终计数器值: %d (期望: 0)\n", finalValue)

	if finalValue == 0 {
		fmt.Printf("  ✅ 原子操作正确完成！\n")
	} else {
		fmt.Printf("  ❌ 计数器值异常！\n")
	}
}

// 实用的原子计数器类型
type AtomicCounter struct {
	value int64
}

func (c *AtomicCounter) Increment() int64 {
	return atomic.AddInt64(&c.value, 1)
}

func (c *AtomicCounter) Decrement() int64 {
	return atomic.AddInt64(&c.value, -1)
}

func (c *AtomicCounter) Add(delta int64) int64 {
	return atomic.AddInt64(&c.value, delta)
}

func (c *AtomicCounter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}

func (c *AtomicCounter) Swap(newValue int64) int64 {
	return atomic.SwapInt64(&c.value, newValue)
}

func (c *AtomicCounter) CompareAndSwap(old, new int64) bool {
	return atomic.CompareAndSwapInt64(&c.value, old, new)
}

// 使用封装好的原子计数器
func demonstrateEncapsulatedCounter() {
	fmt.Println("\n5. 封装的原子计数器使用:")

	counter := &AtomicCounter{}
	var wg sync.WaitGroup
	goroutineCount := 5

	wg.Add(goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				counter.Increment()
				if j%50 == 0 {
					fmt.Printf("  协程%d: 当前值=%d\n", id, counter.Value())
				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("  最终值: %d\n", counter.Value())
}
