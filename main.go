package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	sleeping = iota
	cutting
)

type Barber struct {
	name string
	sync.Mutex
	state int // sleeping or cutting
}

type Customer struct {
	name string
}

func NewBarber() *Barber {
	return &Barber{
		name:  "The Barber",
		state: sleeping,
	}
}

func BarberWork(waitingCustomers chan *Customer, b *Barber) {

	for {
		b.Lock()
		defer b.Unlock()
		b.state = sleeping

		fmt.Printf("Currently, there are %v clients in the waiting room", len(waitingCustomers))

		select {
		case c := <-waitingCustomers:
			DoesHairCut(c, b)

		default:
			fmt.Println("The barber is sleeping")
			b.state = sleeping
			b.Unlock()

		}

	}

}

func DoesHairCut(c *Customer, b *Barber) {
	b.state = cutting
	fmt.Printf("Cutting hair to %v", c.name)
	time.Sleep(time.Millisecond * 100)
	fmt.Printf("Hair cut to %v is finished", c.name)

}

func main() {

}
