package json

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_parse(t *testing.T) {
	lx, err := lex([]byte(`{"hello": ["htello", "world", false, true, ["nested array!!", "woh"], 1235]}`))
	fmt.Println(lx, err)
	px, err := parse(lx)
	fmt.Printf("%#v %v\n", px, err)
	s, _ := json.MarshalIndent(px, "", "\t");
	fmt.Println(string(s))
}