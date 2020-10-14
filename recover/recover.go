package recover

import (
	"fmt"
	"go-dog/log"
	"runtime/debug"
)

//Recover  Recover
func Recover(v ...interface{}) {
	if err := recover(); err != nil {
		if len(v) > 0 {
			s := fmt.Sprintln(v...)
			log.Errorf("%v\n%s\n%s\n", err, s, string(debug.Stack()))
		} else {
			log.Errorf("%v\n%s\n", err, string(debug.Stack()))
		}
	}
}
