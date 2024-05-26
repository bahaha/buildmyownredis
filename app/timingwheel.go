package main

import (
	"container/list"
	"time"
)

// TODO: make my own timing wheel
// References:
// - http://russellluo.com/2018/10/golang-implementation-of-hierarchical-timing-wheels.html

type TimingWheel struct {
	duration time.Duration
	size     int
	slots    []*list.List
}

func NewTimingWheel(duration time.Duration, size int) *TimingWheel {
	tw := &TimingWheel{
		duration: duration,
		size:     size,
		slots:    make([]*list.List, size),
	}

	for i := range tw.slots {
		tw.slots[i] = list.New()
	}

	go tw.run()

	return tw
}

func (tw *TimingWheel) run() {

}
