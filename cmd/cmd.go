package main

import (
	"fmt"
	"github.com/Licoy/stail"
	"os"
	"strconv"
)

func main() {
	st, err := stail.New(stail.Options{})
	if err != nil {
		panic(err)
	}
	atoi, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	si, err := st.Tail(os.Args[1], atoi, func(content string) {
		fmt.Print(content)
	})
	if err != nil {
		panic(err)
	}
	si.Watch()
}
