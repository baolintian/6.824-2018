package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	crand "crypto/rand"
	"fmt"
	"labrpc"
	"log"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

// import "bytes"
// import "labgob"

//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in Lab 3 you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh; at that point you can add fields to
// ApplyMsg, but set CommandValid to false for these other uses.
//

const candidate = "candidate"
const follower = "follower"
const master = "master"

const electionTime = time.Millisecond * 1000
const appendTime = time.Millisecond * 50

type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.
	currentTerm    int
	votedFor       int
	state          string
	electionTimer  *time.Timer
	heartbeatTimer *time.Timer
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	// Your code here (2A).
	if rf.state == master {
		isleader = true
	} else {
		isleader = false
	}
	term = rf.currentTerm

	return term, isleader

}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}

//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term        int
	CandidateId int
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Term        int
	VoteGranted bool

	RequestId int
}

type AppendEntryArgs struct {
	Term     int
	LeaderId int
}

type AppendEntryReply struct {
	Term    int
	Success bool
}

func (rf *Raft) resetTimer(d time.Duration, ty int) {

	if ty == 1 {
		rf.electionTimer.Stop()
		rf.electionTimer.Reset(d)
	} else {
		rf.heartbeatTimer.Stop()
		rf.heartbeatTimer.Reset(d)
	}
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
	//上锁有问题
	// 不是锁有问题，是生成的时间不是随机数字，，，

	rf.mu.Lock()
	defer rf.mu.Unlock()
	if args.Term < rf.currentTerm || (args.Term == rf.currentTerm && rf.votedFor != -1 && rf.votedFor != args.CandidateId) {
		reply.VoteGranted, reply.Term = false, args.Term
		return
	}
	//fmt.Println("requestVote, ", args.Term, ", ", rf.currentTerm, "")
	//fmt.Printf("In requestVote: %d -> %d\n", args.CandidateId, rf.me)
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term
		rf.stepDown(follower)
	}

	rf.votedFor = args.CandidateId
	reply.Term = rf.currentTerm
	reply.VoteGranted = true

	rf.resetTimer(randDuration(electionTime), 1)

	return

}

func (rf *Raft) AppendEntry(args *AppendEntryArgs, reply *AppendEntryReply) {
	rf.mu.Lock()
	defer rf.mu.Unlock()

	if args.Term < rf.currentTerm {
		reply.Success, reply.Term = false, rf.currentTerm
		return
	}
	if args.Term > rf.currentTerm {
		rf.currentTerm = args.Term

		rf.stepDown(follower)
	}
	//fmt.Printf("In Append: reset electionTimer, %d -> %d\n", args.LeaderId, rf.me)
	rf.resetTimer(randDuration(electionTime), 1)
	reply.Success = true
	reply.Term = rf.currentTerm
	return
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
	reply.RequestId = server
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

func (rf *Raft) sendAppendEntry(server int, args *AppendEntryArgs, reply *AppendEntryReply) bool {

	ok := rf.peers[server].Call("Raft.AppendEntry", args, reply)
	return ok
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election. even if the Raft instance has been killed,
// this function should return gracefully.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true

	// Your code here (2B).

	return index, term, isLeader
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
}

func init() {
	//labgob.Register(LogEntry{})
	max := big.NewInt(int64(1) << 62)
	bigx, _ := crand.Int(crand.Reader, max)
	seed := bigx.Int64()
	rand.Seed(seed)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

// generate random time duration that is between minDuration and 2x minDuration
func randDuration(minDuration time.Duration) time.Duration {
	// fmt.Println(rand.Int63() % 1000)
	extra := time.Duration(rand.Int63()) % minDuration
	return time.Duration(minDuration + extra)
}

// func randDuration(minTime time.Duration) time.Duration {
// 	rand.Seed(time.Now().Unix())
// 	fmt.Println(rand.Int63()%1000, ", ", time.Now().Unix())

// 	return time.Duration(rand.Int63())%minTime + minTime
// }

func (rf *Raft) requestBeLeader() {
	if rf.state == "master" {
		return
	}

	rf.currentTerm++
	fmt.Printf("ID: %d Increase Term: %d\n", rf.me, rf.currentTerm)

	rf.votedFor = rf.me
	voteCount := 1

	var req RequestVoteArgs
	req.Term = rf.currentTerm
	req.CandidateId = rf.me
	//req.lastLogTerm = rf.log[len(rf.log)-1].LogTerm
	//req.lastLogIndex = rf.log[len(rf.log)-1].LogIndex

	//fmt.Println("server numbers: ", len(rf.peers))
	//之前设计一种死循环的轮询方式不对，询问完一遍没有结果，就应该等待超时，进行新一轮的查询
	for i := 0; i < len(rf.peers); i++ {
		if i != rf.me {
			go func(index int) {
				var reply RequestVoteReply
				// fmt.Printf("sendrequestvote: %d -> %d\n", rf.me, index)
				if rf.sendRequestVote(index, &req, &reply) {
					if reply.VoteGranted {
						voteCount++
						if voteCount > len(rf.peers)/2 && rf.state == candidate {
							fmt.Printf("In term %d: %d get selected\n", rf.currentTerm, rf.me)
							rf.stepDown(master)
						}
					} else {
						if rf.currentTerm < reply.Term {
							rf.currentTerm = reply.Term
							rf.stepDown(follower)
						}
					}
				}
			}(i)
		}
	}

}

func (rf *Raft) stepDown(s string) {
	if rf.state == s {
		return
	}
	switch s {
	case follower:
		rf.votedFor = -1
		rf.resetTimer(randDuration(electionTime), 1)
		rf.heartbeatTimer.Stop()

	case master:
		rf.electionTimer.Stop()
		// fmt.Println("In stepDown Broadcast")
		rf.Broadcast()
		rf.resetTimer(randDuration(appendTime), 2)

	case candidate:
		rf.requestBeLeader()
	}
	rf.state = s
}

func (rf *Raft) Broadcast() {
	for i := range rf.peers {
		if i == rf.me {
			continue
		}
		go func(index int) {
			rf.mu.Lock()
			if rf.state != master {
				rf.mu.Unlock()
				return
			}
			args := AppendEntryArgs{
				Term:     rf.currentTerm,
				LeaderId: rf.me,
			}
			var reply AppendEntryReply
			if rf.sendAppendEntry(index, &args, &reply) {

			}
			rf.mu.Unlock()
		}(i)
	}
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here (2A, 2B, 2C).

	rf.currentTerm = 0
	rf.votedFor = -1

	rf.state = candidate

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	rf.electionTimer = time.NewTimer(randDuration(electionTime))

	rf.heartbeatTimer = time.NewTimer(randDuration(electionTime))
	rf.heartbeatTimer.Stop()

	go func() {
		for {
			select {
			case <-rf.electionTimer.C:
				rf.mu.Lock()
				//fmt.Printf("in electionTimer reset, raft id: %d\n", rf.me)
				rf.resetTimer(randDuration(electionTime), 1)
				if rf.state == follower {
					rf.stepDown(candidate)
				} else {
					rf.requestBeLeader()
				}
				rf.mu.Unlock()
			case <-rf.heartbeatTimer.C:
				rf.mu.Lock()
				//fmt.Printf("In timeout Broadcast, leader ID: %d\n", rf.me)
				rf.Broadcast()
				//rf.heartbeatTimer = time.NewTimer(randDuration(appendTime))
				rf.resetTimer(randDuration(appendTime), 2)

				rf.mu.Unlock()
			}

		}
	}()

	return rf
}
