package json

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_parse(t *testing.T) {
	lx, err := lex([]byte(`{"hello": "\uD834\uDD1E"}`))
	fmt.Println(err)
	px, err := parse(lx)
	fmt.Println(err)
	s, _ := json.MarshalIndent(px, "", "\t")
	fmt.Println(string(s))
}
