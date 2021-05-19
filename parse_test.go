package json

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

func Test_parse(t *testing.T) {
	const dirName = "testdata"
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		var requiredOutcome string
		{
			sp := strings.Split(f.Name(), "_")
			requiredOutcome = sp[0]
		}
		var okay, status bool
		fcont, err := ioutil.ReadFile(dirName + "/" + f.Name())
		if err != nil {
			t.Fatal(err)
		}
		fmt.Print(f.Name() + " ")
		_, err = Load(fcont)
		switch requiredOutcome {
		case "y":
			okay = err == nil
		case "n":
			okay = err != nil
		case "i":
			okay = true
		}
		status = err == nil
		fmt.Print("okay: ", okay, " status: ", status)
		if !okay {
			fmt.Println("", err)
			t.Fatal("Failed")
		} else {
			fmt.Println()
		}
	}

}