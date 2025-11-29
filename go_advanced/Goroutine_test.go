//编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
func main() {
    go printOddNumbers()
    go printEvenNumbers()
    var input string
    fmt.Scanln(&input)
}

func printOddNumbers() {
    for i := 1; i <= 10; i += 2 {
        fmt.Println("奇数：", i)
    }
}

func printEvenNumbers() {
    for i := 2; i <= 10; i += 2 {}
}

//设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。

func main() {
    tasks := []func(){
        task1,
        task2,
        task3,
    }

    var wg sync.WaitGroup
    wg.Add(len(tasks))

    for _, task := range tasks {
        go func(t func()) {
            start := time.Now()
            t()
            elapsed := time.Since(start)
            fmt.Printf("任务执行完成，耗时：%s\n", elapsed)
        }
    }
}

