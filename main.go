package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	sleeping = iota
	cutting
	checking
)

var wg *sync.WaitGroup

type Barber struct {
	name                  string
	customerBeingAttended *Customer
	sync.Mutex
	state int // sleeping, cutting or checking
}

type Customer struct {
	name string
}

func NewBarber() *Barber {
	return &Barber{
		name:  "Bob the Barber",
		state: sleeping,
	}
}

func barberWork(waitingCustomers chan *Customer, b *Barber) {

	for {
		b.Lock()
		defer b.Unlock()
		b.state = checking
		b.customerBeingAttended = nil

		fmt.Printf("Currently, there are %v clients in the waiting room \n", len(waitingCustomers))
		time.Sleep(time.Millisecond * 100)

		select {
		case c := <-waitingCustomers:
			doesHairCut(c, b)
			b.Unlock()
		default:
			fmt.Println("The barber is sleeping because the waiting room is empty")
			b.state = sleeping
			b.customerBeingAttended = nil
			b.Unlock()

		}

	}

}

func doesHairCut(c *Customer, b *Barber) {
	b.state = cutting

	time.Sleep(time.Millisecond * 100)

	b.customerBeingAttended = c

	b.Unlock()

	fmt.Printf("Cutting hair to %v \n", c.name)
	time.Sleep(time.Millisecond * 100)
	fmt.Printf("Hair cut to %v is finished \n", c.name)

	b.Lock()
	b.customerBeingAttended = nil
	wg.Done()
}

func customerArrivesToBarberShop(waitingCustomers chan<- *Customer, c *Customer, b *Barber) {

	b.Lock()

	defer b.Unlock()

	switch b.state {

	case sleeping:
		select {
		case waitingCustomers <- c:
		default:
			wg.Done()
		}

	case cutting:
		select {
		case waitingCustomers <- c:
		default:
			fmt.Println("Waiting room is full, customer called ", c.name, " is leaving")
			wg.Done()
		}

	case checking:
		fmt.Println("Customer called ", c.name, " bumped into the barber while he was checking the waiting room, he will wait")
		waitingCustomers <- c
	}

}

func main() {

	barber := NewBarber()

	waitingRoomChairs := 5

	waitingRoom := make(chan *Customer, waitingRoomChairs)

	go barberWork(waitingRoom, barber)

	wg = new(sync.WaitGroup)

	for _, c := range customerGenerator(5) {

		wg.Add(1)

		go customerArrivesToBarberShop(waitingRoom, c, barber)
	}

	wg.Wait()

	fmt.Println("All customers were sorted")

}

func customerGenerator(customerAmount int) []*Customer {

	customerSliceToReturn := make([]*Customer, customerAmount)

	for i := 0; i < customerAmount; i++ {
		customerSliceToReturn[i] = &Customer{name: fmt.Sprintf("Customer-%v", i)}
	}

	return customerSliceToReturn
}
