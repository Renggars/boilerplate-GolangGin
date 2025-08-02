package utils

import (
	"crypto/rand"
	"math/big"
)

func GenerateOTP() string {
	n := 6
	otp := ""
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			panic(err)
		}
		otp += num.String()
	}
	return otp
}
