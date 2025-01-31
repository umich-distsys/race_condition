package value

// Implementation of Value confinement from lecture
type Value struct {
	adder 	chan<- int
	current <-chan int
	end 	chan any
}

func (v Value) Add(delta int) {
	v.adder <- delta
}

func (v Value) Current() int {
	return <-v.current
}

func (v Value) Done() {
	close(v.end)
}

func MakeValue() Value {
	adder 	:= make(chan int)
	current := make(chan int)
	end 	:= make(chan any)
	
	val := Value{
		adder: 	adder,
		current: 	current,
		end: 		end,
	}

	go func() {
		val := 0
		for {
			select {
			case delta := <-adder:
				val += delta
			case current <- val:
			case <-end:
				return
			}
		}
	} ()

	return val
}