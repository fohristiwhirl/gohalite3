package core

func (self *Game) SetPid(pid int) {

	// For simulation purposes, it's simplest just to have
	// each AI set its PID at the start of its turn...

	self.pid = pid
}
