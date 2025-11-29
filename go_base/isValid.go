func isValid(s string) bool {
	stack := []rune{}

	bracketMap := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	for _, char := range s {
		if matchingLeft, isRight := bracketMap[char]; isRight {
			if len(stack) == 0 || stack[len(stack)-1] != matchingLeft {
				return false
			}
			stack = stack[:len(stack)-1]

		} else {
			stack = append(stack, char)
		}
	}
	return len(stack) == 0
}