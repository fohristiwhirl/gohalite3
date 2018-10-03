package core

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// ---------------------------------------

type TokenParser struct {
	scanner		*bufio.Scanner
	count		int
}

func NewTokenParser() *TokenParser {
	ret := new(TokenParser)
	ret.scanner = bufio.NewScanner(os.Stdin)
	ret.scanner.Split(bufio.ScanWords)
	return ret
}

func (self *TokenParser) Int() int {
	bl := self.scanner.Scan()
	if bl == false {
		err := self.scanner.Err()
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		} else {
			panic(fmt.Sprintf("End of input."))
		}
	}
	ret, err := strconv.Atoi(self.scanner.Text())
	if err != nil {
		panic(fmt.Sprintf("TokenReader.Int(): Atoi failed at token %d: \"%s\"", self.count, self.scanner.Text()))
	}

	self.count++
	return ret
}

func (self *TokenParser) Float() float64 {
	bl := self.scanner.Scan()
	if bl == false {
		err := self.scanner.Err()
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		} else {
			panic(fmt.Sprintf("End of input."))
		}
	}
	ret, err := strconv.ParseFloat(self.scanner.Text(), 64)
	if err != nil {
		panic(fmt.Sprintf("TokenReader.Float(): ParseFloat failed at token %d: \"%s\"", self.count, self.scanner.Text()))
	}

	self.count++
	return ret
}

func (self *TokenParser) Bool() bool {
	bl := self.scanner.Scan()
	if bl == false {
		err := self.scanner.Err()
		if err != nil {
			panic(fmt.Sprintf("%v", err))
		} else {
			panic(fmt.Sprintf("End of input."))
		}
	}
	val, err := strconv.Atoi(self.scanner.Text())
	if err != nil {
		panic(fmt.Sprintf("TokenReader.Bool(): Atoi failed at token %d: \"%s\"", self.count, self.scanner.Text()))
	}
	if val != 0 && val != 1 {
		panic(fmt.Sprintf("TokenReader.Bool(): Non-bool at token %d: \"%s\"", self.count, self.scanner.Text()))
	}

	self.count++
	if val == 0 {
		return false
	} else {
		return true
	}
}

// ---------------------------------------

func (self *Game) Parse() {
	return
}

func (self *Game) Send() {
	return
}
