package raft

import "strconv"

func (rf *Raft) leaderElection() {
	rf.currentTerm++
	rf.votedFor = rf.me
	rf.state = Candidate
	requestVoteArgs := RequestVoteArgs{
		Term:         rf.currentTerm,
		CandidateId:  rf.me,
		LastLogIndex: rf.log.getLastEntry().Index,
		LastLogTerm:  rf.log.getLastEntry().Term,
	}

	DPrintf("candidate " + strconv.Itoa(rf.me) + " Election Term " + strconv.Itoa(rf.currentTerm))
	DPrintf("requestVoteArgs is %v", requestVoteArgs)
	voteCounter := 1
	for peerNo := 0; peerNo < len(rf.peers); peerNo++ {
		if peerNo != rf.me {
			go rf.CandidateSendRequestToPeer(peerNo, &requestVoteArgs, &voteCounter)
		}
	}
}

func (rf *Raft) CandidateSendRequestToPeer(peerNo int, args *RequestVoteArgs, voteCounter *int) {

	reply := RequestVoteReply{}

	DPrintf("candidate " + strconv.Itoa(rf.me) + " Send Request To Peer " + strconv.Itoa(peerNo))
	if !rf.sendRequestVote(peerNo, args, &reply) {
		return
	}
	DPrintf("candidate catch vote result %v", reply)

	rf.mu.Lock()
	defer rf.mu.Unlock()

	if !reply.VoteGranted {
		return
	}
	if reply.Term > args.Term {
		rf.newTerm(reply.Term)
		return
	}
	if reply.Term < args.Term {
		return
	}

	*voteCounter++

	if *voteCounter > len(rf.peers)/2 && rf.state == Candidate && rf.currentTerm == args.Term {
		rf.becomeLeader()
	}

	DPrintf("candidate rf.state %v", rf.state)

}

func (rf *Raft) becomeLeader() {
	rf.state = Leader

	lastLog := rf.log.getLastEntry()
	for peerNo := range rf.peers {
		rf.nextIndex[peerNo] = lastLog.Index + 1
		rf.matchIndex[peerNo] = 0
	}
	DPrintf("candidate LeaderSendEntriesToPeer")
	rf.logReplication(true)

}
