package oxp_test

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"github.com/man0xff/oxp"
)

func TestParser(t *testing.T) {
	p := oxp.NewClient()
	resp, err := p.Search(context.Background(), "home")
	if err != nil {
		t.Fatal(err)
	}
	spew.Dump(resp)
}
