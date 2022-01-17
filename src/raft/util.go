package raft

import "log"

// Debugging
const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

// log entries; each entry contains command
// for state machine, and term when entry
// was received by leader (first index is 1)
type Log struct{
	Command interface{}
	Term int
	Index int
}

func makeEmptyLog()  []Log{
	emptyLog := make([]Log, 0);
	return emptyLog
}
