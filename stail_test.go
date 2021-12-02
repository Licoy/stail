package stail

import (
	"fmt"
	"testing"
)

func TestSTail(t *testing.T) {
	st, err := New(Options{})
	if err != nil {
		t.Error(err)
		return
	}
	err = st.Tail("D:/stail.log", -1, func(content string) {
		fmt.Print(fmt.Sprintf("%s", content))
	})
	if err != nil {
		t.Error(err)
		return
	}
}
