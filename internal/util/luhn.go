package util

func Valid(number uint64) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number uint64) uint64 {
	var out uint64

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		out += cur
		number = number / 10
	}
	return out % 10
}
