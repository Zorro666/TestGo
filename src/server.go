package main

import (
	"fmt"
)

type request struct {
	a, b   int
	replyc chan int
}

type binOp func(a, b int) int

func run(op binOp, req *request) {
	reply := op(req.a, req.b)
	req.replyc <- reply
}

func server(op binOp, server chan *request) {
	for {
		req := <-server
		go run(op, req) // don't wait for it
	}
}

func startServer(op binOp) chan *request {
	req := make(chan *request)
	go server(op, req)
	return req
}

func main() {
	adder := startServer(func(a, b int) int { return a + b })
	const N = 100
	var reqs [N]request
	for i := 0; i < N; i++ {
		req := &reqs[i]
		req.a = i
		req.b = i + N
		req.replyc = make(chan int)
		adder <- req
	}
	for i := N - 1; i >= 0; i-- { // doesn't matter what order
		result := <-reqs[i].replyc
		fmt.Println(i, "Result=", result)
		if result != N+2*i {
			fmt.Println("fail at", i)
		}
	}
	fmt.Println("done")
}
