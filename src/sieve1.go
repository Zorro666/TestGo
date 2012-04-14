package main

import (
	"fmt"
)

// Send the sequence 2, 3, 4, ... to the created and returned channel 
func generate() chan int {
	ch := make(chan int)
	go func() {
		for i := 2; ; i++ {
			ch <- i // Send 'i' to channel 'ch'
		}
	}()
	return ch
}

// Copy the values from channel 'in' to the created and returned channel
// removing those divisible by 'prime'
func filter(in chan int, prime int) chan int {
	out := make(chan int)
	go func() {
		for {
			i := <-in // Receive value of new variable 'i' from 'in'.
			if i%prime != 0 {
				out <- i // Send 'i' to channel 'out'.
			}
		}
	}()
	return out
}

func sieve() chan int {
	out := make(chan int)
	go func() {
		ch := generate()
		for {
			prime := <-ch
			out <- prime
			ch = filter(ch, prime)
		}
	}()
	return out
}

func main() {
	primes := sieve()
	for {
		fmt.Println(<-primes)
	}
}
