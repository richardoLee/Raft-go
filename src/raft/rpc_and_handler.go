package raft

import "strconv"

//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
    // Your data here (2A, 2B).
    Term         int
    CandidateId  int
    LastLogIndex int
    LastLogTerm  int
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
    // Your data here (2A).
    Term        int
    VoteGranted bool
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
    // Your code here (2A, 2B).
    rf.mu.Lock()
    defer rf.mu.Unlock()
    rf.resetElectionTimer()
    DPrintf("peer "+strconv.Itoa(rf.me)+" catch args %v", args)
    if args.Term < rf.currentTerm {
        reply.VoteGranted = false
    } else {
    }
    DPrintf("peer "+strconv.Itoa(rf.me)+" part 1= %v", (rf.votedFor == -1 || rf.votedFor == args.CandidateId))
    DPrintf("peer "+strconv.Itoa(rf.me)+" part 2= %v", rf.UpToDateCheck(args.LastLogIndex, args.LastLogTerm))
    if (rf.votedFor == -1 || rf.votedFor == args.CandidateId) && rf.UpToDateCheck(args.LastLogIndex, args.LastLogTerm) {
        reply.VoteGranted = true
    }
    reply.Term = rf.currentTerm
    DPrintf("peer "+strconv.Itoa(rf.me)+" reply is %v", reply)
}

func (rf *Raft) UpToDateCheck(lastLogIndex int, lastLogTerm int) bool {
    localLastLog := rf.log.getLastLog()
    // DPrintf("peer "+strconv.Itoa(rf.me)+" lastLogIndex >= localLastLog.Index %v",lastLogIndex >= localLastLog.Index)
    if localLastLog.Term == lastLogTerm {
        return lastLogIndex >= localLastLog.Index
    } else {
        return lastLogTerm > localLastLog.Term
    }
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
    ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
    return ok
}





type AppendEntriesArgs struct {
    Term         int
    LeaderId     int
    PrevLogIndex int

    PrevLogTerm  int
    Entries      []Entry
    LeaderCommit int
}

type AppendEntriesReply struct {
    Term    int
    Success bool
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply) {
    reply.Term=rf.currentTerm
    
    // if args.Term<rf.currentTerm {
    // 	reply.Success=false
    // }

    if args.Entries==nil {
        rf.resetElectionTimer()
    }
}

func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool {
    ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)
    return ok
}
