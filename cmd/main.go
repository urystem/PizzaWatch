package main

import (
	"fmt"
	"os"
	"strings"

	"pizza/cmd/kitchen"
	"pizza/cmd/notification"
	order "pizza/cmd/order"
	"pizza/cmd/tracking"
)

func main() {
	os.Args = os.Args[1:]
	if len(os.Args) == 0 {
		return
	}
	if len(os.Args) > 1 && (os.Args[0] == "--mode" || os.Args[0] == "-mode") {
		os.Args = os.Args[1:]
	} else if ind := strings.Index(os.Args[0], "="); ind != -1 && (os.Args[0][:ind] == "--mode" || os.Args[0][:ind] == "-mode") {
		os.Args[0] = os.Args[0][ind+1:]
	}
	switch os.Args[0] {
	case "order-service":
		order.Main()
	case "kitchen-worker":
		kitchen.Main()
	case "tracking-service":
		tracking.Main()
	case "notification-subscriber":
		notification.Main()
	default:
		fmt.Println(os.Args[0])
	}
}
