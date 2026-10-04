package main

import (
	"fmt"
	"time"
)

func processTask(name string) {
	start := time.Now()

	fmt.Println("Starting:", name)

	time.Sleep(2 * time.Second)

	fmt.Println("Completed:", name)

	elapsed := time.Since(start)
	fmt.Println("Time taken:", elapsed)
	fmt.Println()
}

func main() {
	tasks := []string{
		"task1",
		"task2",
		"task3",
		"task4",
		"task5",
	}

	fmt.Println("========== WITHOUT GO ==========")

	for _, task := range tasks {
		processTask(task)
	}

	fmt.Println("========== WITH GO ==========")

	for _, task := range tasks {
		go processTask(task)
	}

	time.Sleep(3 * time.Second)
}
