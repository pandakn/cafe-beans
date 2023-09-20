package utils

// The number of bits is based on the number of roles in db
// this function for convert decimal to binary
func BinaryConverter(number int, bits int) []int {
	factor := number
	result := make([]int, bits)

	for factor >= 0 && number > 0 {
		factor = number % 2
		number /= 2
		result[bits-1] = factor
		bits--
	}
	return result
}
