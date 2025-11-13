package main

import (
	"flag"
	order "pizza/cmd/order"
)

// import ""

func main() {
	mode := flag.String("mode", "", "select the mode")
	flag.Parse()
	switch *mode {
	case "order-service":
		order.Main()
	case "kitchen-worker":
	case "tracking-service":
	case "notification-subscriber":
	default:
		flag.Usage()
	}
}
