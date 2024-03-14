package common

import (
	"log"
	"testing"

	"golang.org/x/exp/slices"
)

func Test_Slices(t *testing.T) {
	s := make([]int, 0)

	s = append(s, 1, 2, 3, 4, 5)

	log.Println(s)
	s = slices.Delete(s, 0, 1)
	log.Println(s)

	s1 := slices.Delete(s, 3, 4)
	log.Printf("s:%v s1:%v\n", s, s1)

}
