//定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
package main

import (
    "fmt"
    "math"
)


type Shape interface {
    Area() float64
    Perimeter() float64
}


type Rectangle struct {
    Width  float64
    Height float64
}


func (r Rectangle) Area() float64 {
    return r.Width * r.Height
}


func (r Rectangle) Perimeter() float64 {
    return 2 * (r.Width + r.Height)
}


type Circle struct {
    Radius float64
}


func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
    return 2 * math.Pi * c.Radius
}

func main() {
    fmt.Println("=== 接口实现演示 ===\n")

    rect := Rectangle{Width: 5, Height: 3}
    fmt.Printf("矩形: 宽度 = %.2f, 高度 = %.2f\n", rect.Width, rect.Height)
    fmt.Printf("面积: %.2f\n", rect.Area())
    fmt.Printf("周长: %.2f\n\n", rect.Perimeter())
    

    circle := Circle{Radius: 4}
    fmt.Printf("圆形: 半径 = %.2f\n", circle.Radius)
    fmt.Printf("面积: %.2f\n", circle.Area())
    fmt.Printf("周长: %.2f\n\n", circle.Perimeter())
    

    fmt.Println("=== 使用接口类型 ===")
    demonstrateInterface()
    

    fmt.Println("\n=== 切片中的多态 ===")
    demonstratePolymorphism()
    
    fmt.Println("\n=== 类型断言 ===")
    demonstrateTypeAssertion()
}

//使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
package main

import "fmt"


type Person struct {
    Name string
    Age  int
}


func (p Person) Introduce() string {
    return fmt.Sprintf("大家好，我是%s，今年%d岁", p.Name, p.Age)
}


type Employee struct {
    Person      // 匿名嵌入，组合
    EmployeeID  string
    Department  string
    Position    string
    Salary      float64
}


func (e Employee) PrintInfo() {
    fmt.Println("=== 员工信息 ===")
    fmt.Printf("员工ID: %s\n", e.EmployeeID)
    fmt.Printf("姓名: %s\n", e.Name)        // 直接访问嵌入结构的字段
    fmt.Printf("年龄: %d\n", e.Age)         // 直接访问嵌入结构的字段
    fmt.Printf("部门: %s\n", e.Department)
    fmt.Printf("职位: %s\n", e.Position)
    fmt.Printf("薪资: %.2f\n", e.Salary)
    fmt.Printf("个人介绍: %s\n", e.Introduce()) // 直接调用嵌入结构的方法
    fmt.Println()
}


func (e Employee) Work() {
    fmt.Printf("%s 正在%s部门工作...\n", e.Name, e.Department)
}

func main() {
    fmt.Println("=== 组合结构体演示 ===\n")
    

    emp1 := Employee{
        Person: Person{
            Name: "张三",
            Age:  28,
        },
        EmployeeID: "E1001",
        Department: "技术部",
        Position:   "高级工程师",
        Salary:     15000.00,
    }
    
    emp2 := Employee{
        Person: Person{
            Name: "李四",
            Age:  32,
        },
        EmployeeID: "E1002",
        Department: "市场部",
        Position:   "市场经理",
        Salary:     18000.00,
    }
    

    emp1.PrintInfo()
    emp2.PrintInfo()
    
  
    demonstrateMoreFeatures()
    

    demonstrateSliceOperations()
}

