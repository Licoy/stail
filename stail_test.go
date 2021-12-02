package stail

import (
	"fmt"
	"testing"
	"time"
)

func TestSTail(t *testing.T) {
	st, err := New(Options{})
	if err != nil {
		t.Error(err)
		return
	}
	var si1, si2 STailItem
	go func() {
		si1, err = st.Tail("D:/stail.log", -1, func(content string) {
			fmt.Print(fmt.Sprintf("%s", content))
		})
		if err != nil {
			t.Error(err)
		} else {
			si1.Watch()
		}
	}()
	time.AfterFunc(time.Second*5, func() {
		cErr := si1.Close()
		if err != nil {
			t.Error(cErr)
		}
		fmt.Println("stail 1 stopped")
	})
	si2, err = st.Tail("D:/stail2.log", -1, func(content string) {
		fmt.Print(fmt.Sprintf("%s", content))
	})
	if err != nil {
		t.Error(err)
	} else {
		si2.Watch()
	}
}
