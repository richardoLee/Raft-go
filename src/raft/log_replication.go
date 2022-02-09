package raft

func (rf *Raft) logReplication(heartbeat bool) {
    DPrintf("leader %v entries %v",rf.me,rf.log.Entries)
    DPrintf("leader %v nextIndex %v",rf.me,rf.nextIndex)
    DPrintf("leader %v matchIndex %v",rf.me,rf.matchIndex)
    DPrintf("leader %v commitIndex %v",rf.me,rf.commitIndex)

	lastLog := rf.log.getLastEntry()

	for peerNo := 0; peerNo < len(rf.peers); peerNo++ {
		if peerNo == rf.me {
			rf.resetElectionTimer()
		} else {
			if lastLog.Index >= rf.nextIndex[peerNo] || heartbeat {
				nextIndex := rf.nextIndex[peerNo]
				prevLog := rf.log.Entries[nextIndex-1]

				appendEntriesArgs := AppendEntriesArgs{
					Term:         rf.currentTerm,
					LeaderId:     rf.me,
					PrevLogIndex: prevLog.Index,
					PrevLogTerm:  prevLog.Term,
					Entries:      rf.log.getAppenedEntries(nextIndex),
					LeaderCommit: rf.commitIndex,
				}
                DPrintf("peerNo %v appendEntriesArgs Entries are %v",peerNo,appendEntriesArgs.Entries)

				go rf.LeaderSendEntriesToPeer(peerNo, &appendEntriesArgs)
			}
		}
	}
}

func (rf *Raft) LeaderSendEntriesToPeer(peerNo int, args *AppendEntriesArgs) {
	reply := AppendEntriesReply{}
	if !rf.sendAppendEntries(peerNo, args, &reply) {
		return
	}
    DPrintf("peerNo %v reply is %v",peerNo,reply)
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if reply.Term > rf.currentTerm {
		rf.newTerm(reply.Term)
		return
	}
	if args.Term == rf.currentTerm {
		if reply.Success {
			matchIdx := args.PrevLogIndex + len(args.Entries)
			rf.matchIndex[peerNo] = max(rf.matchIndex[peerNo], matchIdx)
			nextIdx := matchIdx + 1
			rf.nextIndex[peerNo] = max(rf.nextIndex[peerNo], nextIdx)
		} else if rf.nextIndex[peerNo] > 1 {
			rf.nextIndex[peerNo]--
		}
		rf.LeaderCommit()
	}
}

func (rf *Raft) LeaderCommit() {
	if rf.state != Leader {
		return
	}

	for N := rf.commitIndex + 1; N <= rf.log.getLastEntry().Index; N++ {

		if rf.log.Entries[N].Term != rf.currentTerm {
			continue
		}

		counter := 1
		for peerNo := 0; peerNo < len(rf.peers); peerNo++ {
			if peerNo != rf.me && rf.matchIndex[peerNo] >= N {
				counter++
			}
			if counter > len(rf.peers)/2 {
				rf.commitIndex = N
				rf.apply()
				break
			}
		}
	}
}
