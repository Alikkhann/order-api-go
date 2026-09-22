package resp

import (
	"time"
	"math/rand"
)

func GenerateVerificationCode() int {
	src := rand.NewSource(time.Now().UnixNano()) 
	r := rand.New(src)                           
	return r.Intn(900000) + 100000              
}
