package enum

import (
	"testing"
)

// type BitFlag int
//
// const (
// 	// iota为0，1左移0位 = 1
// 	Active BitFlag = 1 << iota
// 	// Send <=> Active <=> 1 << iota，此时iota为1，1左移1位 = 2
// 	Send
// 	// Receive <=> Send <=> 1 << iota，此时iota为2，1左移2位 = 4
// 	Receive
// )
//
// func TestIota1(t *testing.T) {
// 	fmt.Println(Active, Send, Receive)
// }

func TestIota2(t *testing.T) {
	t.Log(LayoutTypeLink, LayoutTypeDirectory, LayoutTypeFile, LayoutTypeWorkspace, LayoutTypeMainBox)
}
