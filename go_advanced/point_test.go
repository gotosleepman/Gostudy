//题目 ：编写一个Go程序，定义一个函数，该函数接收一个整数指针作为参数，在函数内部将该指针指向的值增加10，然后在主函数中调用该函数并输出修改后的值。
func main() {
	var num int = 5
	var ptr *int = &num
	addTen(ptr)
	fmt.Println("修改后的值：", num)
}

func addTen(ptr *int) {
	*ptr += 10
}

//实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
func main() {
	nums := []int{1, 2, 3, 4, 5}
	doubleSlice(nums)
	fmt.Println(nums)
}

func doubleSlice(slice *[]int) {
	for i := 0; i < len(*slice); i++ {
		(*slice)[i] *= 2
	}
}