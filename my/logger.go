package my



import (
	"fmt"
// "reflect"

)

type Logger func(...any)

func (l Logger) Err(format string, a ...any) error {
	err := fmt.Errorf(format, a...)
	l("ERROR:", err)
	return err
}

func (l Logger) Typeof(v any) {
	l(fmt.Sprintf("TYPE: %T", v ))
	// l(fmt.Sprintf("TYPE: %v", reflect.TypeOf(v)))
}

var Log Logger = func(a ...any) {
	fmt.Println(a...)
}