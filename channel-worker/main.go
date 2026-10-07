package main

import (
	"fmt"
	"time"
)

type Job struct {
	id int
}

var jobs = make(chan Job, 10)

func worker(id int) {
	for job := range jobs {

		fmt.Printf("Worker %d processing Job %d\n", id, job.id)

		time.Sleep(2 * time.Second)

		fmt.Printf("Worker %d finished Job %d\n", id, job.id)
	}
}

func main() {
	for i := 1; i <= 3; i++ {
		go worker(i)
	}

	for j := 1; j <= 10; j++ {
		jobs <- Job{id: j}
	}

	fmt.Println("all job done added ")

	time.Sleep(10 * time.Second)
}

// okay so we have some question here

// why the hell we use channel-worker or why not we use goroutine or semaphore here ?

// let me explain u every single thing here

//what go routine is ?

// imagine u run a hotel and u have 100 order came  and every order take 2 minute
// so what u can do is take one order at a time and it will take all to process 200 min or what u can do is
// u can just hire 100 cook and then it will take everything to process in 2 min
// but here is a catch this is called goroutine in go

// it seem good for the small work-load Imagine  now u have 10k order now ur cpu will blow up if u use go process(Tasks) here

// so what is a semaphore ? now we can use that instead

// its also the same as the channel but it control the number of things can hapepn at the same time
// like we can't control the goroutine in the same way it do everything at same time but here
// we can control that thing

// it not necceryly have the queues but it can used to control the number of concurrent task running actively

// so what channel worker do

// it basically creeate worker which take the job from the queue and then work on that

// imagine u have 100 order and now u are saying go process(tasks) now it will do everything at the same time like 100 at same time concurrently but

// if its big then ur cpu will fucked up

// so where we come the worker as we have 100 order we will create 3 worker which do the task but also when a order complete it will ask for the another order from the channel
// and channel is something where we are taking the order from like it will store the order

// so what's happening is we are creating 3 worker and they are taking the order from the channel and when they complete the order they will ask for another order from the channel and this way we can control the number of concurrent tasks running at the same time and also we can control the load on our CPU and also we can control the memory usage as well
