package main

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestCreateOrChannel(t *testing.T) {
	randInts := [5]int{rand.Intn(5) + 1, rand.Intn(5) + 1, rand.Intn(5) + 1, rand.Intn(5) + 1, rand.Intn(5) + 1}
	lowest := 6
	startTime := time.Now()
	for _, val := range randInts {
		if val < lowest {
			lowest = val
		}
	}
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	<-createOrChannel(sig(time.Duration(randInts[0])*time.Second),
		sig(time.Duration(randInts[1])*time.Second),
		sig(time.Duration(randInts[2])*time.Second),
		sig(time.Duration(randInts[3])*time.Second),
		sig(time.Duration(randInts[4])*time.Second))
	finishTime := math.Round(time.Since(startTime).Seconds())
	if int(finishTime) != lowest {
		t.Fail()
	}
	t.Log(finishTime)
	t.Log(lowest)
}
