package shortcode

import "math/rand"

// GenerateShortCode string

type ShortCode struct {
	length int
}

const chars = "1234567890abcdefghijklmnopqrstuvwxyz"

func NewShortCode(n int) *ShortCode {
	return &ShortCode{
		length: n,
	}
}

func (s *ShortCode) GenerateShortCode() string {
	length := len(chars)
	result := make([]byte, s.length)

	for i := 0; i < s.length; i++ {
		result[i] = chars[rand.Intn(length)]
	}
	return string(result)
}
