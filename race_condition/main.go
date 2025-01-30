package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

type Airplane struct {
	Seats SeatsConfinement
}

func (a *Airplane) MakeAirplane(numSeats int) {
	a.Seats = MakeSeatsConfinement(numSeats)
}

func (a *Airplane) BookSeat() bool {
	return a.Seats.BookSeat()
}

type SeatsConfinement struct {
	bookSeatChan 		chan<- any
	bookSeatSuccessChan <-chan bool
	numSeatsBookedChan 	<-chan int
	termination 		chan any
}

func (s SeatsConfinement) BookSeat() bool {
	s.bookSeatChan <- struct{}{}
	return <-s.bookSeatSuccessChan
}

func (s SeatsConfinement) GetNumSeatsBooked() int {
	return <-s.numSeatsBookedChan
}

func (s SeatsConfinement) Done() {
	close(s.termination)
}

func MakeSeatsConfinement(numSeats int) SeatsConfinement {
	bookSeatChan 		:= make(chan any)
	bookSeatSuccessChan := make(chan bool)
	numSeatsBookedChan 	:= make(chan int)
	termination 		:= make(chan any)

	s := SeatsConfinement{
		bookSeatChan: bookSeatChan,
		bookSeatSuccessChan: bookSeatSuccessChan,
		numSeatsBookedChan: numSeatsBookedChan,
		termination: termination,
	}

	go func() {
		seats := make([]bool, numSeats)
		numSeatsBooked := 0
		for {
			select {
			case <-bookSeatChan:
				// Iterates through Seats and books the first available seat
				success := false
				for i := 0; i < len(seats); i++ {
					if !seats[i] {
						seats[i] = true
						fmt.Printf("Seat %d booked\n", i)
						numSeatsBooked++
						success = true
						break
					}
				}
				bookSeatSuccessChan <- success
			case numSeatsBookedChan <- numSeatsBooked: 
			case <-termination:
				return
			}
		}
	} ()

	return s
}

func main() {
	usage := "Usage: go run main.go [num seats] [num customers]"
	if len(os.Args) != 3 {
		fmt.Println(usage)
		return
	} 

	numSeats, seatErr := strconv.Atoi(os.Args[1])
	numCustomers, customerErr := strconv.Atoi(os.Args[2])
	if seatErr != nil || customerErr != nil {
		fmt.Println(usage)
		return
	}

	a := Airplane{}
	a.MakeAirplane(numSeats)

	var bookingWG sync.WaitGroup
	bookingWG.Add(numCustomers)
	for i := 0; i < numCustomers; i++ {
		go func() {
			defer bookingWG.Done()
			a.BookSeat()
		} ()
	}

	bookingWG.Wait()
	
	// To consistently expose race condition, use 1000 seats and 10000 customers
	fmt.Printf("Number of seats on airplane: %d\n", numSeats)
	fmt.Printf("Number of seats booked: %d\n", a.Seats.GetNumSeatsBooked())
	a.Seats.Done()
}