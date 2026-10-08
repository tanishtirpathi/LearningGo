package main

// it don't exactly the different from the channel worker but still it work like this kinda

//Producer → Queue → Consumers

// this is something code look like

// jobs := make(chan Job, 10)

// for i := 1; i <= 3; i++ {
//     go worker(i)
// }

// for i := 1; i <= 10; i++ {
//     jobs <- Job{ID: i}
// }

// in more simple term it look like this

// Producer
//    ↓
// Channel
//    ↓
// Workers
