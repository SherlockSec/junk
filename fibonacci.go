package main

import (
	"fmt"
	//"strconv"
	//"math"
	"math/big"
)

func main() {
	//var nInput uint64
	//nInput = 10000000
	//var x uint64
	
	//for x = 0; x <= nInput; x++ {
	
		//fmt.Println(fibonacci(x))
		
	//}
	
	//binFibonacci()
	primeFibonacci()
}

func fibonacci(n uint64) uint64 {
	
	slice := make([]uint64, n+1, n+2)
	if n < 2 {
		slice = slice[0:2]	
	}
	slice[0] = 0; slice[1] = 1;
	var i uint64
	for i = 2; i <= n; i++ {
		slice[i] = slice[i-1] + slice[i-2]
	}
	
	return slice[n]

}


func binFibonacci() {
	a := big.NewInt(0)
	b := big.NewInt(1)
	var limit big.Int
	limit.Exp(big.NewInt(10), big.NewInt(999), nil)

	for a.Cmp(&limit) < 0 {
		a.Add(a, b)
		a, b = b, a
		
		fmt.Println(a)
		//fmt.Println(a.ProbablyPrime(20))
		//fmt.Println("\n")
	}
	//fmt.Println(a.ProbablyPrime(20))

}

func primeFibonacci() {
	a := big.NewInt(0)
	b := big.NewInt(1)
	var limit big.Int
	limit.Exp(big.NewInt(10), big.NewInt(999), nil)

	for a.Cmp(&limit) < 0 {
		a.Add(a, b)
		a, b = b, a
		
		//fmt.Println(a)
		//fmt.Println(a.ProbablyPrime(20))
		//fmt.Println("\n")
		
		
		if test := a.ProbablyPrime(20); test != false {
			fmt.Println(a)
			//fmt.Println(a.ProbablyPrime(20))
			//fmt.Println("\n")
		}
	}
	//fmt.Println(a.ProbablyPrime(20))

}
