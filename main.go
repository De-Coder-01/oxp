package main

import (
	"context"
	"fmt"

	"github.com/man0xff/oxp"
)

func main() {
	fmt.Println("Hello")
	nc := oxp.NewClient()
	//fmt.Println(nc)

	ctx := context.TODO()

	res, err := nc.Search(ctx, "abound")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}
