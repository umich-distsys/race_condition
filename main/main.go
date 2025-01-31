package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"umich.edu/eecs491/d4/race_condition/value"
)

type Airplane struct {
	Seats 			[]bool // Seats[i] is true if seat i is booked, false otherwise
	NumSeatsBooked	value.Value // Value confinement from lecture      
}

func (a *Airplane) MakeAirplane(numSeats int) {
	a.Seats = make([]bool, numSeats)
	a.NumSeatsBooked = value.MakeValue()
}

func (a *Airplane) GetNumSeatsBooked() int {
	return a.NumSeatsBooked.Current()
}

func (a *Airplane) BookSeat() bool {
	// Iterates through Seats and books the first available seat
	for i := 0; i < len(a.Seats); i++ {
		if !a.Seats[i] {
			a.Seats[i] = true
			fmt.Printf("Seat %d booked\n", i)
			a.NumSeatsBooked.Add(1)
			return true
		}
	}
	return false
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
	fmt.Printf("Number of seats booked: %d\n", a.GetNumSeatsBooked())
}