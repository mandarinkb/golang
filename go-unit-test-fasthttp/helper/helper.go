package helper

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// for unit test
func Equals(tb testing.TB, expected, actual interface{}) {

	// not true
	if !reflect.DeepEqual(expected, actual) { //!reflect.DeepEqual(expected, actual)
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\texp: %#v\n\tgot: %#v\033[39m\n", filepath.Base(file), line, expected, actual)
		tb.FailNow()
	}
}
