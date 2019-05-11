package rand

import (
	rnd "crypto/rand"
)

var (
	Reader = rnd.Reader
)

func GetBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rnd.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

const (
	chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTVWXYZ0123456789"
)

func GetString(s int) (string, error) {
	b, err := GetBytes(s + 1)
	if err != nil {
		return "", err
	}
	for i := 0; i < s; i++ {
		val := int64(b[i])*256 + int64(b[i+1])
		val *= int64(len(chars))
		val /= (256 * 256)
		b[i] = chars[val]
	}

	return string(b[0:s]), nil
}
