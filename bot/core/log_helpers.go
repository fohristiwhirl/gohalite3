package core

import (
	"fmt"
	"../logging"
)



func (self *Frame) Log(format_string string, args ...interface{}) {
	format_string = fmt.Sprintf("t %3d: ", self.Turn()) + format_string
	logging.Log(format_string, args...)
}

func (self *Frame) LogOnce(format_string string, args ...interface{}) bool {
	format_string = "t %3d: " + format_string
	var newargs []interface{}
	newargs = append(newargs, self.Turn())
	newargs = append(newargs, args...)
	return logging.LogOnce(format_string, newargs...)
}

func (self *Frame) LogWithoutTurn(format_string string, args ...interface{}) {
	logging.Log(format_string, args...)
}



func (self *Frame) Flog(x, y int, msg, colour string) {
	logging.Flog(self.Turn(), x, y, msg, colour)
}
