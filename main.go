package main

import (
	"fmt"

	polygonClient "github.com/PacmanHQ/teleport/client/polygon_client"
	polygonListener "github.com/PacmanHQ/teleport/listener/polygon_listener"
	"github.com/PacmanHQ/teleport/order"
)

func main() {
	a, err := polygonClient.NewPolygonClient("", "", "", 0)
	if err != nil {
		panic(err)
	}
	c := make(chan order.Order)
	b := polygonListener.NewPolygonListener(1, *a, &c)

	go b.Start()

for{
	fmt.Println(<-c)
}

}
