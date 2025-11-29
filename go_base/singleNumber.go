func singleNumber(nums []int) int {
	result := 0
	for _, char := range nums {
		result ^= char
	}
	return result
}