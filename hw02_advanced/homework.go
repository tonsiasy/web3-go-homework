package homework02

import (
	"fmt"
	"sync"
	// "sync/atomic"
	"time"
)

// Q1: 指针 - 增加值
// 编写一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10
func AddTen(val *int) {
	*val += 10
}

// Q2: 指针 - 切片元素乘2
// 实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2
func DoubleSlice(slice *[]int) {
	for i, v := range *slice {
		(*slice)[i] = 2 * v
	}
}

// Q3: Goroutine - 奇偶数打印
// 编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
// 提示：为了测试能看到输出，可以使用 sync.WaitGroup 确保主程序等待协程结束
func PrintOddEven() {
	var wg sync.WaitGroup

	wg.Go(func() {
		for i := 0; i <= 10; i++ {
			if i&0x01 == 1 {
				fmt.Printf("Odd:%v\n", i)
			}
		}
	})

	wg.Go(func() {
		for i := 0; i <= 10; i++ {
			if i&0x01 == 0 {
				fmt.Printf("Even:%v\n", i)
			}
		}
	})

	wg.Wait()
}

// Q4: Goroutine - 任务调度
// 设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
func TaskScheduler(tasks []func()) {
	var wg sync.WaitGroup
	for i, task := range tasks {
		id := i
		t := task
		wg.Go(func() {
			start := time.Now()
			t()
			fmt.Printf("任务ID: %d 耗时: %v\n", id, time.Since(start))
		})
	}

	wg.Wait()
}

// Q5: 面向对象 - 接口
// 定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Width  float64
	Height float64
}

// TODO: 为 Rectangle 实现 Shape 接口
func (r *Rectangle) Area() float64 {
	// TODO: 实现代码
	return 0
}

func (r *Rectangle) Perimeter() float64 {
	// TODO: 实现代码
	return 0
}

type Circle struct {
	Radius float64
}

// TODO: 为 Circle 实现 Shape 接口
func (c *Circle) Area() float64 {
	// TODO: 实现代码
	return 0
}

func (c *Circle) Perimeter() float64 {
	// TODO: 实现代码
	return 0
}

// Q6: 面向对象 - 组合
// 使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段
type Person struct {
	Name string
	Age  int
}

// 再创建一个 Employee 结构体，组合 Person 结构体并添加 EmployeeID 字段
type Employee struct {
	Person
	EmployeeID string
}

// TODO: 为 Employee 结构体实现一个 PrintInfo() string 方法，返回员工的信息 (格式自定，包含Name, Age, EmployeeID)
func (e Employee) PrintInfo() string {
	// TODO: 实现代码
	return ""
}

// Q7: Channel - 生产者消费者
// 编写一个程序，使用通道实现两个协程之间的通信。
// 一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
// 为了方便测试，请将接收到的数字返回
func ProducerConsumer() []int {
	var result []int
	// TODO: 实现代码
	return result
}

// Q8: Channel - 缓冲通道
// 实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
func BufferedChannel() {
	// TODO: 实现代码
}

// Q9: 锁机制 - Mutex
// 编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。
// 启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
func MutexCounter() int {
	var count int
	// var mu sync.Mutex
	// TODO: 实现代码
	return count
}

// Q10: 锁机制 - Atomic
// 使用原子操作（ sync/atomic 包）实现一个无锁的计数器。
// 启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
func AtomicCounter() int32 {
	var count int32
	// TODO: 实现代码
	return count
}
