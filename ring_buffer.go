package main

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	rb := NewRingBuffer(5)
	fmt.Println("Empty Test:")
	spew.Dump(rb.Emit())

	currentRune := 'a' - 1
	for i := 0; i < 10; i++ {
		currentRune++
		rb.Insert(Data{
			Stamp: time.Now(),
			Value: string(currentRune),
		})
	}
	fmt.Println("Empty Test:")
	spew.Dump(rb.Emit())
}

type RingBuffer struct {
	data         []*Data
	size         int
	lastInserted int
	nextRead     int
	emitTime     time.Time
}

type Data struct {
	Stamp time.Time
	Value string
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		size:         size,
		lastInserted: -1,
		data:         make([]*Data, size),
	}
}

func (r *RingBuffer) Insert(input Data) {
	r.lastInserted = (r.lastInserted + 1) % r.size
	r.data[r.lastInserted] = &input

	if r.nextRead == r.lastInserted {
		r.nextRead = (r.nextRead + 1) % r.size
	}
}

func (r *RingBuffer) Emit() []*Data {
	var output []*Data
	for {
		if r.data[r.nextRead] != nil {
			output = append(output, r.data[r.nextRead])
			r.data[r.nextRead] = nil
		}
		if r.nextRead == r.lastInserted || r.lastInserted == -1 {
			break
		}
		r.nextRead = (r.nextRead + 1) % r.size
	}
	return output
}
