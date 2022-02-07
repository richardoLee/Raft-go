package raft

import (
    "log"
    "time"
    "math/rand"
)

// Debugging
const Debug = true

func DPrintf(format string, a ...interface{}) (n int, err error) {
    if Debug {
        log.Printf(format, a...)
    }
    return
}

func (rf *Raft) resetElectionTimer() {
    now := time.Now()
    timeOut := time.Duration(150+rand.Intn(150)) * time.Millisecond
    rf.electionTime = now.Add(timeOut)
}

type Log struct {
    Entries []Entry
}

// log entries; each entry contains command
// for state machine, and term when entry
// was received by leader (first index is 1)
type Entry struct {
    Command interface{}
    Term    int
    Index   int
}

func makeEmptyLog() Log {
    emptyLog := Log{Entries: make([]Entry, 0)}

    initialEntry := Entry{Index: 0}
    emptyLog.appendLog(initialEntry)
    return emptyLog
}

func (l *Log) getLastLog() Entry {
    return l.Entries[len(l.Entries)-1]
}

func (l *Log) appendLog(entry Entry) {
    l.Entries = append(l.Entries, entry)
}
