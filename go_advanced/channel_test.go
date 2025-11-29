//编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
package main

import (
    "fmt"
    "time"
)

func main() {
    fmt.Println("=== 通道通信基本演示 ===\n")
    

    ch := make(chan int)
    

    go producer(ch)
    

    go consumer(ch)
    

    time.Sleep(2 * time.Second)
    fmt.Println("程序执行完毕")
}


func producer(ch chan<- int) {
    fmt.Println("生产者启动...")
    
    for i := 1; i <= 10; i++ {
        fmt.Printf("生产者发送: %d\n", i)
        ch <- i // 发送数据到通道
        time.Sleep(100 * time.Millisecond)
    }
    
    close(ch) 
    fmt.Println("生产者完成")
}


func consumer(ch <-chan int) {
    fmt.Println("消费者启动...")
    
    for {
        num, ok := <-ch 
        if !ok {
            fmt.Println("通道已关闭，消费者退出")
            return
        }
        fmt.Printf("消费者接收: %d\n", num)
        time.Sleep(150 * time.Millisecond) 
    }
}


//实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
package main

import (
    "fmt"
    "sync"
    "time"
)

func main() {
    fmt.Println("=== 带缓冲通道的生产者-消费者演示 ===\n")
    

    bufferSize := 20
    ch := make(chan int, bufferSize)
    
    var wg sync.WaitGroup
    

    wg.Add(1)
    go producer(ch, &wg)
    
 
    wg.Add(1)
    go consumer(ch, &wg)
    

    wg.Wait()
    fmt.Println("\n所有任务完成！")
}


func producer(ch chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    defer close(ch) 
    
    totalNumbers := 100
    fmt.Printf("生产者启动，将生成 %d 个整数，缓冲区大小: %d\n", totalNumbers, cap(ch))
    
    for i := 1; i <= totalNumbers; i++ {

        time.Sleep(50 * time.Millisecond)
        

        ch <- i
        

        bufferUsage := len(ch)
        bufferCapacity := cap(ch)
        fmt.Printf("生产者发送: %3d | 缓冲区: %2d/%2d | 使用率: %3.0f%%\n", 
                   i, bufferUsage, bufferCapacity, float64(bufferUsage)/float64(bufferCapacity)*100)
    }
    
    fmt.Printf("生产者完成，共发送 %d 个整数\n", totalNumbers)
}


func consumer(ch <-chan int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    fmt.Println("消费者启动，开始接收数据...")
    receivedCount := 0
    
    for num := range ch {
     
        time.Sleep(80 * time.Millisecond)
        
        receivedCount++
        
      
        bufferUsage := len(ch)
        fmt.Printf("消费者接收: %3d | 已接收: %3d | 缓冲区剩余: %2d\n", 
                   num, receivedCount, bufferUsage)
    }
    
    fmt.Printf("消费者完成，共接收 %d 个整数\n", receivedCount)
}