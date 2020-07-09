package test

import (
	"fmt"
	"log"
	"math/rand"
	"shopping/utils"
	"testing"
	"time"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().Unix()))
}

// RandString 生成随机字符串
func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}

func TestHash(t *testing.T) {
	consistent := utils.NewConsistent(20)
	consistent.Add("127.0.0.1")
	consistent.Add("127.0.0.2")
	consistent.Add("127.0.0.3")
	consistent.Add("127.0.0.4")

	for i := 0; i < 10; i++ {
		s := RandString(10)
		ip, err := consistent.Get(s)
		if err != nil {
			log.Panic(err.Error())
		}
		fmt.Println(s, "对应的ip应该为->", ip)
	}
	for i := 0; i < 10; i++ {
		s := RandString(10)
		ip, err := consistent.Get(s)
		if err != nil {
			log.Panic(err.Error())
		}
		fmt.Println(s, "对应的ip应该为->", ip)
	}

}
