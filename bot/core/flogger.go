package core

import (
	"encoding/json"
	"fmt"
	"os"
)

// This is a simple logger that I use for saving a JSON array of objects.

type Flogfile struct {
	outfile			*os.File
	outfilename		string
	at_start		bool
	failed			bool
}

type FlogObject struct {
	T				int			`json:"t"`
	X				int			`json:"x"`
	Y				int			`json:"y"`
	Msg				string		`json:"msg"`
}

func NewFlog(outfilename string) *Flogfile {
	return &Flogfile{
		nil,
		outfilename,
		true,
		false,
	}
}

func (self *Flogfile) Flog(t, x, y int, msg string) {

	if self == nil || self.failed {
		return
	}

	if self.outfile == nil {
		var err error
		self.outfile, err = os.Create(self.outfilename)
		if err != nil {
			self.failed = true
			return
		}
	}

	f := FlogObject{T: t, X: x, Y: y, Msg: msg}

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
	fmt.Fprintf(self.outfile, "\n]")
	self.outfile.Close()
}

// ---------------------------------------------------------------
// Methods on the Game object...

func (self *Game) StartFlog(flogfilename string) {
	self.flogfile = NewFlog(flogfilename)
}

func (self *Game) Flog(t, x, y int, msg string) {
	self.flogfile.Flog(t, x, y, msg)
}

func (self *Game) StopFlog() {
	self.flogfile.Close()
}
