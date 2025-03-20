package main

import "fmt"

func main() {
	c := make(chan int)
	go func() {
		defer close(c)
		for i := 0; i < 100; i++ {
			c <- i
		}
	}()

	for {
		if i, ok := <-c; ok {
			fmt.Println(i)
		} else {
			break
		}
	}
}
