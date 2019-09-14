package process

import (
	"os"
	"reflect"
	"unsafe"
)

// NOTE: Set process name, as in the name seen in `ps`
func (self *Process) SetName(name string) {
	argv0str := (*reflect.StringHeader)(unsafe.Pointer(&os.Args[0]))
	argv0 := (*[1 << 30]byte)(unsafe.Pointer(argv0str.Data))[:argv0str.Len]
	n := copy(argv0, name)
	if n < len(argv0) {
		argv0[n] = 0
	}
}
