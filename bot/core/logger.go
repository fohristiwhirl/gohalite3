package core

import (
	"encoding/json"
	"fmt"
	"os"
)

type Logfile struct {
	outfile			*os.File
	outfilename		string
	logged_once		map[string]bool
	closed			bool
}

func NewLog(outfilename string) *Logfile {
	return &Logfile{
		nil,
		outfilename,
		make(map[string]bool),
		false,
	}
}

func (self *Logfile) Log(format_string string, args ...interface{}) {

	if self == nil || self.closed {
		return
	}

	if self.outfile == nil {

		var err error

		if _, tmp_err := os.Stat(self.outfilename); tmp_err == nil {
			// File exists
			self.outfile, err = os.OpenFile(self.outfilename, os.O_APPEND|os.O_WRONLY, 0666)
		} else {
			// File needs creating
			self.outfile, err = os.Create(self.outfilename)
		}

		if err != nil {
			self.closed = true
			return
		}
	}

	fmt.Fprintf(self.outfile, format_string, args...)
}

func (self *Logfile) LogOnce(format_string string, args ...interface{}) bool {
	if self.logged_once[format_string] == false {
		self.logged_once[format_string] = true         // Note that it's format_string that is checked / saved
		self.Log(format_string, args...)
		return true
	}
	return false
}

func (self *Logfile) Close() {
	if self == nil || self.closed {
		return
	}
	self.closed = true
	self.outfile.Close()
}

// ---------------------------------------------------------------
// Methods on the Game object...

func (self *Game) StartLog(logfilename string) {
	self.logfile = NewLog(logfilename)
}

func (self *Game) Log(format_string string, args ...interface{}) {
	format_string = fmt.Sprintf("t %3d: ", self.Turn()) + format_string + "\r\n"
	self.logfile.Log(format_string, args...)
}

func (self *Game) LogOnce(format_string string, args ...interface{}) bool {
	format_string = "t %3d: " + format_string + "\r\n"
	var newargs []interface{}
	newargs = append(newargs, self.Turn())
	newargs = append(newargs, args...)
	return self.logfile.LogOnce(format_string, newargs...)
}

func (self *Game) LogWithoutTurn(format_string string, args ...interface{}) {
	self.logfile.Log(format_string + "\r\n", args...)
}

func (self *Game) StopLog() {
	self.logfile.Close()
}



// ---------------------------------------------------------------
// This is a simple logger that I use for saving a JSON array of
// objects for later interpretation by Fluorine.

type Flogfile struct {
	outfile			*os.File
	outfilename		string
	at_start		bool
	closed			bool
}

type FlogObject struct {
	T				int			`json:"t"`
	X				int			`json:"x"`
	Y				int			`json:"y"`
	Msg				string		`json:"msg,omitempty"`
	Colour			string		`json:"colour,omitempty"`
}

func NewFlog(outfilename string) *Flogfile {
	return &Flogfile{
		nil,
		outfilename,
		true,
		false,
	}
}

func (self *Flogfile) Flog(t, x, y int, msg, colour string) {

	// msg or colour can be ""

	if self == nil || self.closed {
		return
	}

	if self.outfile == nil {
		var err error
		self.outfile, err = os.Create(self.outfilename)
		if err != nil {
			self.closed = true
			return
		}
	}

	f := FlogObject{T: t, X: x, Y: y, Msg: msg, Colour: colour}

	s, _ := json.Marshal(f)

	if self.at_start {
		fmt.Fprintf(self.outfile, "[\n  ")
		self.at_start = false
	} else {
		fmt.Fprintf(self.outfile, ",\n  ")
	}

	fmt.Fprintf(self.outfile, string(s))
}

func (self *Flogfile) Close() {
	if self == nil || self.closed {
		return
	}
	fmt.Fprintf(self.outfile, "\n]")
	self.closed = true
	self.outfile.Close()
}

// ---------------------------------------------------------------
// Methods on the Game object...

func (self *Game) StartFlog(flogfilename string) {
	self.flogfile = NewFlog(flogfilename)
}

func (self *Game) Flog(x, y int, msg, colour string) {
	self.flogfile.Flog(self.Turn(), x, y, msg, colour)
}

func (self *Game) StopFlog() {
	self.flogfile.Close()
}
