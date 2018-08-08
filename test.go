package main

import (
	"addlib/addlib"
	"fmt"
	"time"
)

func test() {
	t1 := time.Now()
	for i:=0; i < 100; i++ {
		go func () {
			for _, p := range addlib.Provinces() {
				for _, c := range addlib.Cities(p) {
					for _, d := range addlib.Districts(c) {
						_ = d
					}
				}
			}
		}()
	}
	t2 := time.Since(t1)
	fmt.Println(t2)
}

func main() {
	test()
}

