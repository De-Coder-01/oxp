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
	resp, _ := p.Search(context.Background(), "abound")


	var entries = []*Entry(resp.([]*Entry))

	fmt.Println(entries[0].Senses[0].Def)
	fmt.Println(entries[0].Senses[0].Examples[0])

}
