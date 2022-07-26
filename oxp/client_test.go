package oxp_test

import (
	"context"
	"fmt"
	"testing"

	//"github.com/gdexlab/go-render/render"
	"github.com/man0xff/oxp"
	. "github.com/man0xff/oxp"
)

func TestParser(t *testing.T) {
	p := oxp.NewClient()
	resp, err := p.Search(context.Background(), "abound")
	if err != nil {
		t.Fatal(err)
	}
	//spew.Dump(resp)

	//fmt.Printf("%v\n", []*Entry(resp.([]*Entry)))
	var entries = []*Entry(resp.([]*Entry))
	// fmt.Println(entries)
	for _, val := range entries {
		for _, s := range val.Senses {
			fmt.Println(s)
		}
	}
	//output := render.AsCode(resp)
	//fmt.Println(output)
}
