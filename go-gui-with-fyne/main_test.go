package main

import (
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/paulbellamy/ratecounter"
)

func TestXxx(t *testing.T) {
	a := float64(0.0001)
	// str := fmt.Sprintf("%f", a)
	// fmt.Println(str)
	if hasDecimal(a) {
		fmt.Println("has decimal")
	}
}
func hasDecimal(value float64) bool {
	return value != math.Trunc(value)
}

func Test2(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(6)
	counter := ratecounter.NewRateCounter(1 * time.Second)
	go func(c *ratecounter.RateCounter) {
		time.Sleep(time.Second)
		defer wg.Done()
		fmt.Println("job 1", time.Now())
		call(c)
		fmt.Println("number of concurrent: ", c.Rate())
	}(counter)
	go func(c *ratecounter.RateCounter) {
		defer wg.Done()
		fmt.Println("job 2", time.Now())
		call(c)
		fmt.Println("number of concurrent: ", c.Rate())
	}(counter)
	go func(c *ratecounter.RateCounter) {
		defer wg.Done()
		fmt.Println("job 3", time.Now())
		call(c)
		fmt.Println("number of concurrent: ", c.Rate())
	}(counter)
	go func(c *ratecounter.RateCounter) {
		time.Sleep(2 * time.Second)
		defer wg.Done()
		fmt.Println("job 4", time.Now())
		call(c)
		fmt.Println("number of concurrent: ", c.Rate())
	}(counter)
	go func(c *ratecounter.RateCounter) {
		time.Sleep(3 * time.Second)
		defer wg.Done()
		fmt.Println("job 5", time.Now())
		call(c)
		fmt.Println("number of concurrent: ", c.Rate())
	}(counter)

	go func(c *ratecounter.RateCounter) {
		defer wg.Done()
		fmt.Println("job 6", time.Now())
		call(c)
		fmt.Println("number of concurrent: ", c.Rate())
	}(counter)

	wg.Wait()
}

func call(c *ratecounter.RateCounter) {
	c.Incr(1)
}
