// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package bloombits

import (
	"bytes"
	"context"
	"errors"
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common/bitutil"
	"github.com/ethereum/go-ethereum/crypto"
)

// bloomIndexes represents the bit indexes inside the bloom filter that belong
// to some key.
type bloomIndexes [3]uint //因为每个元素映射到三个位置

// calcBloomIndexes returns the bloom filter bit indexes belonging to the given key.
func calcBloomIndexes(b []byte) bloomIndexes {
	b = crypto.Keccak256(b)

	var idxs bloomIndexes
	for i := 0; i < len(idxs); i++ {
		idxs[i] = (uint(b[2*i])<<8)&2047 + uint(b[2*i+1])
	}
	return idxs
}

//部分匹配。因为一次过滤可能不止一个条件，我们假设有三个条件 A, B, C，对单个条件的匹配叫做子匹配或者部分匹配。
//bitset 表示这对单个条件的匹配结果的向量，它会在后面通过和其他条件的 bitset 取与，达到同时满足多个条件的效果。

// partialMatches with a non-nil vector represents a section in which some sub-
// matchers have already found potential matches. Subsequent sub-matchers will
// binary AND their matches with this vector. If vector is nil, it represents a
// section to be processed by the first sub-matcher.
type partialMatches struct {
	section uint64
	bitset  []byte
}

//匹配中用于给 request 和 response 传递结果的结构，表示一次检索任务。

// Retrieval represents a request for retrieval task assignments for a given
// bit with the given number of fetch elements, or a response for such a request.
// It can also have the actual results set to be used as a delivery data struct.
//
// The contest and error fields are used by the light client to terminate matching
// early if an error is encountered on some path of the pipeline.
type Retrieval struct {
	Bit      uint //bit 表示检索的位，调度器也是按位安排检索的。
	Sections []uint64
	Bitsets  [][]byte //bitsets 表示 Sections 中每个检索结果向量构成的矩阵。

	//用于 eth 客户端终止过滤的匹配操作
	Context context.Context
	Error   error //外部调用时的出错提示
}

//匹配器，通过流水线的方式，承担二进制向量取与、或的任务，并且保留可能包含检索内容的区块。

// Matcher is a pipelined system of schedulers and logic matchers which perform
// binary AND/OR operations on the bit-streams, creating a stream of potential
// blocks to inspect for data content.
type Matcher struct {
	//段的大小，默认 4096 个区块
	sectionSize uint64 // Size of the data batches to filter on
	//这个流水线处理的过滤操作，第一个维度表示需要检索第几位，第二个维度表示这需要检索一位的元素
	filters [][]bloomIndexes // Filter the system is matching for
	//一次匹配工作包括多个调度器，因为调度器是按照一个位检索的，一次匹配至少检索 3 个位
	schedulers map[uint]*scheduler // Retrieval schedulers for loading bloom bits

	//当需要检索的位置分配好了后，传递检索任务
	retrievers chan chan uint // Retriever processes waiting for bit allocations
	//当一次检索任务完成时，传递当前完成的任务数量
	counters chan chan uint // Retriever processes waiting for task count reports
	//当检索任务分配好后，传递检索任务
	retrievals chan chan *Retrieval // Retriever processes waiting for task allocations
	//当检索完成后，传递任务的结果 response
	deliveries chan *Retrieval // Retriever processes waiting for task response deliveries
	// 指定是否运行的原子变量
	running uint32 // Atomic flag whether a session is live or not
}

//新建流水线，通过给定 topic 结构(即智能合约中事件带有 indexed 标识的变量）或者地址在 bloom 比特流中筛选。
//filters 第一个维度表示对某一位的检索，第二个维度表示需要检索的元素，第三个维度表示元素的内容，可能经过了特殊的编码。

// NewMatcher creates a new pipeline for retrieving bloom bit streams and doing
// address and topic filtering on them. Setting a filter component to `nil` is
// allowed and will result in that filter rule being skipped (OR 0x11...1).
func NewMatcher(sectionSize uint64, filters [][][]byte) *Matcher {
	// Create the matcher instance
	m := &Matcher{
		sectionSize: sectionSize,
		schedulers:  make(map[uint]*scheduler),
		retrievers:  make(chan chan uint),
		counters:    make(chan chan uint),
		retrievals:  make(chan chan *Retrieval),
		deliveries:  make(chan *Retrieval),
	}
	// Calculate the bloom bit indexes for the groups we're interested in
	m.filters = nil

	for _, filter := range filters {
		// Gather the bit indexes of the filter rule, special casing the nil filter
		//如果这一位需要检索的元素不存在，则跳过。
		if len(filter) == 0 {
			continue
		}
		bloomBits := make([]bloomIndexes, len(filter))
		//对某一位需要检索的元素存在，则计算每一个元素对应的哈希后的三个位置
		for i, clause := range filter {
			if clause == nil {
				bloomBits = nil
				break
			}
			bloomBits[i] = calcBloomIndexes(clause)
		}

		//然后将对某一位需要检索的元素的位置用切片存储在匹配器中

		// Accumulate the filter rules if no nil rule was within
		if bloomBits != nil {
			m.filters = append(m.filters, bloomBits)
		}
	}
	//现在匹配器的 filters 装满了每一位对应的需要检索的元素的位置的矩阵。
	// 我们获取需要检索某一位的多个元素的位置的切片，再获这一位的每个元素，
	//再获取这个元素需要检索的位置，然后分发给 addScheduler 准备操作。

	// For every bit, create a scheduler to load/download the bit vectors
	for _, bloomIndexLists := range m.filters {
		for _, bloomIndexList := range bloomIndexLists {
			for _, bloomIndex := range bloomIndexList {
				m.addScheduler(bloomIndex)
			}
		}
	}
	return m
}

//对这bloom 的某一位新建检索任务，如果它之前还不存在。如果已经存在，那么使用已存在检索任务。

// addScheduler adds a bit stream retrieval scheduler for the given bit index if
// it has not existed before. If the bit is already selected for filtering, the
// existing scheduler can be used.
func (m *Matcher) addScheduler(idx uint) {
	if _, ok := m.schedulers[idx]; ok {
		return
	}
	m.schedulers[idx] = newScheduler(idx)
}

//对指定范围的区块执行匹配操作，返回结果（一串比特流），知道这个范围内的区块的匹配任务都完成
// begin, end 指区块范围
//

// Start starts the matching process and returns a stream of bloom matches in
// a given range of blocks. If there are no more matches in the range, the result
// channel is closed.
func (m *Matcher) Start(ctx context.Context, begin, end uint64, results chan uint64) (*MatcherSession, error) {
	// Make sure we're not creating concurrent sessions

	//在运行时的标志 running 中写入 1，并且返回 running 之前的值。
	//如果已经运行则不必再开始
	if atomic.SwapUint32(&m.running, 1) == 1 {
		return nil, errors.New("matcher already running")
	}
	//最后将 running 设置为 0，表示已经完成。
	defer atomic.StoreUint32(&m.running, 0)

	//新建会话，用于管理这次日志过滤的生命周期

	// Initiate a new matching round
	session := &MatcherSession{
		matcher: m,
		quit:    make(chan struct{}),
		ctx:     ctx,
	}
	for _, scheduler := range m.schedulers {
		scheduler.reset()
	}
	//部分匹配的结果
	sink := m.run(begin, end, cap(results), session)

	// Read the output from the result sink and deliver to the user
	session.pend.Add(1)
	go func() {
		defer session.pend.Done()
		defer close(results)

		for {
			select {
			case <-session.quit:
				return

			case res, ok := <-sink:
				// New match result found
				if !ok {
					return
				}
				// Calculate the first and last blocks of the section
				sectionStart := res.section * m.sectionSize

				first := sectionStart
				if begin > first {
					first = begin
				}
				last := sectionStart + m.sectionSize - 1
				if end < last {
					last = end
				}
				// Iterate over all the blocks in the section and return the matching ones
				for i := first; i <= last; i++ {
					// Skip the entire byte if no matches are found inside (and we're processing an entire byte!)
					next := res.bitset[(i-sectionStart)/8]
					if next == 0 {
						if i%8 == 0 {
							i += 7
						}
						continue
					}
					// Some bit it set, do the actual submatching
					if bit := 7 - i%8; next&(1<<bit) != 0 {
						select {
						case <-session.quit:
							return
						case results <- i:
						}
					}
				}
			}
		}
	}()
	return session, nil
}

//新建两条匹配的子流水线，一条匹配按地址检索，一条匹配按 topic 检索，然后对他们的结果按位取与。
//之所以称为流水线是因为它会一个一个地调用子匹配器，之前的子匹配器找到了匹配的区块后才会调用下一个子匹配器。

// run creates a daisy-chain of sub-matchers, one for the address set and one
// for each topic set, each sub-matcher receiving a section only if the previous
// ones have all found a potential match in one of the blocks of the section,
// then binary AND-ing its own matches and forwarding the result to the next one.
//
// The method starts feeding the section indexes into the first sub-matcher on a
// new goroutine and returns a sink channel receiving the results.
func (m *Matcher) run(begin, end uint64, buffer int, session *MatcherSession) chan *partialMatches {
	// Create the source channel and feed section indexes into
	source := make(chan *partialMatches, buffer)

	session.pend.Add(1)
	go func() {
		defer session.pend.Done()
		defer close(source)

		//将区块高度转化成 section 高度，因为它是匹配的最小单元。
		for i := begin / m.sectionSize; i <= end/m.sectionSize; i++ {
			//初始化部分匹配，二进制全 1，表示完全匹配
			select {
			case <-session.quit:
				return
			case source <- &partialMatches{i, bytes.Repeat([]byte{0xff}, int(m.sectionSize/8))}:
			}
		}
	}()
	// Assemble the daisy-chained filtering pipeline
	next := source
	dist := make(chan *request, buffer)

	//将部分匹配的结果接受者next，结果通信的 dist，需要检索某一位的多个元素的序列，检索任务的一次会话 传入子匹配函数中
	for _, bloom := range m.filters {
		next = m.subMatch(next, dist, bloom, session)
	}
	// Start the request distribution
	session.pend.Add(1)
	go m.distributor(dist, session)

	return next
}

//

// subMatch creates a sub-matcher that filters for a set of addresses or topics, binary OR-s those matches, then
// binary AND-s the result to the daisy-chain input (source) and forwards it to the daisy-chain output.
// The matches of each address/topic are calculated by fetching the given sections of the three bloom bit indexes belonging to
// that address/topic, and binary AND-ing those vectors together.
func (m *Matcher) subMatch(source chan *partialMatches, dist chan *request, bloom []bloomIndexes, session *MatcherSession) chan *partialMatches {
	// Start the concurrent schedulers for each bit required by the bloom filter
	sectionSources := make([][3]chan uint64, len(bloom))
	sectionSinks := make([][3]chan []byte, len(bloom))
	for i, bits := range bloom {
		for j, bit := range bits {
			sectionSources[i][j] = make(chan uint64, cap(source))
			sectionSinks[i][j] = make(chan []byte, cap(source))

			m.schedulers[bit].run(sectionSources[i][j], dist, sectionSinks[i][j], session.quit, &session.pend)
		}
	}

	process := make(chan *partialMatches, cap(source)) // entries from source are forwarded here after fetches have been initiated
	results := make(chan *partialMatches, cap(source))

	session.pend.Add(2)
	go func() {
		// Tear down the goroutine and terminate all source channels
		defer session.pend.Done()
		defer close(process)

		defer func() {
			for _, bloomSources := range sectionSources {
				for _, bitSource := range bloomSources {
					close(bitSource)
				}
			}
		}()
		// Read sections from the source channel and multiplex into all bit-schedulers
		for {
			select {
			case <-session.quit:
				return

			case subres, ok := <-source:
				// New subresult from previous link
				if !ok {
					return
				}
				// Multiplex the section index to all bit-schedulers
				for _, bloomSources := range sectionSources {
					for _, bitSource := range bloomSources {
						select {
						case <-session.quit:
							return
						case bitSource <- subres.section:
						}
					}
				}
				// Notify the processor that this section will become available
				select {
				case <-session.quit:
					return
				case process <- subres:
				}
			}
		}
	}()

	go func() {
		// Tear down the goroutine and terminate the final sink channel
		defer session.pend.Done()
		defer close(results)

		// Read the source notifications and collect the delivered results
		for {
			select {
			case <-session.quit:
				return

			case subres, ok := <-process:
				// Notified of a section being retrieved
				if !ok {
					return
				}
				// Gather all the sub-results and merge them together
				var orVector []byte
				for _, bloomSinks := range sectionSinks {
					var andVector []byte
					for _, bitSink := range bloomSinks {
						var data []byte
						select {
						case <-session.quit:
							return
						case data = <-bitSink:
						}
						if andVector == nil {
							andVector = make([]byte, int(m.sectionSize/8))
							copy(andVector, data)
						} else {
							bitutil.ANDBytes(andVector, andVector, data)
						}
					}
					if orVector == nil {
						orVector = andVector
					} else {
						bitutil.ORBytes(orVector, orVector, andVector)
					}
				}

				if orVector == nil {
					orVector = make([]byte, int(m.sectionSize/8))
				}
				if subres.bitset != nil {
					bitutil.ANDBytes(orVector, orVector, subres.bitset)
				}
				if bitutil.TestBytes(orVector) {
					select {
					case <-session.quit:
						return
					case results <- &partialMatches{subres.section, orVector}:
					}
				}
			}
		}
	}()
	return results
}

// distributor receives requests from the schedulers and queues them into a set
// of pending requests, which are assigned to retrievers wanting to fulfil them.
func (m *Matcher) distributor(dist chan *request, session *MatcherSession) {
	defer session.pend.Done()

	var (
		requests   = make(map[uint][]uint64) // Per-bit list of section requests, ordered by section number
		unallocs   = make(map[uint]struct{}) // Bits with pending requests but not allocated to any retriever
		retrievers chan chan uint            // Waiting retrievers (toggled to nil if unallocs is empty)
		allocs     int                       // Number of active allocations to handle graceful shutdown requests
		shutdown   = session.quit            // Shutdown request channel, will gracefully wait for pending requests
	)

	// assign is a helper method fo try to assign a pending bit an actively
	// listening servicer, or schedule it up for later when one arrives.
	assign := func(bit uint) {
		select {
		case fetcher := <-m.retrievers:
			allocs++
			fetcher <- bit
		default:
			// No retrievers active, start listening for new ones
			retrievers = m.retrievers
			unallocs[bit] = struct{}{}
		}
	}

	for {
		select {
		case <-shutdown:
			// Shutdown requested. No more retrievers can be allocated,
			// but we still need to wait until all pending requests have returned.
			shutdown = nil
			if allocs == 0 {
				return
			}

		case req := <-dist:
			// New retrieval request arrived to be distributed to some fetcher process
			queue := requests[req.bit]
			index := sort.Search(len(queue), func(i int) bool { return queue[i] >= req.section })
			requests[req.bit] = append(queue[:index], append([]uint64{req.section}, queue[index:]...)...)

			// If it's a new bit and we have waiting fetchers, allocate to them
			if len(queue) == 0 {
				assign(req.bit)
			}

		case fetcher := <-retrievers:
			// New retriever arrived, find the lowest section-ed bit to assign
			bit, best := uint(0), uint64(math.MaxUint64)
			for idx := range unallocs {
				if requests[idx][0] < best {
					bit, best = idx, requests[idx][0]
				}
			}
			// Stop tracking this bit (and alloc notifications if no more work is available)
			delete(unallocs, bit)
			if len(unallocs) == 0 {
				retrievers = nil
			}
			allocs++
			fetcher <- bit

		case fetcher := <-m.counters:
			// New task count request arrives, return number of items
			fetcher <- uint(len(requests[<-fetcher]))

		case fetcher := <-m.retrievals:
			// New fetcher waiting for tasks to retrieve, assign
			task := <-fetcher
			if want := len(task.Sections); want >= len(requests[task.Bit]) {
				task.Sections = requests[task.Bit]
				delete(requests, task.Bit)
			} else {
				task.Sections = append(task.Sections[:0], requests[task.Bit][:want]...)
				requests[task.Bit] = append(requests[task.Bit][:0], requests[task.Bit][want:]...)
			}
			fetcher <- task

			// If anything was left unallocated, try to assign to someone else
			if len(requests[task.Bit]) > 0 {
				assign(task.Bit)
			}

		case result := <-m.deliveries:
			// New retrieval task response from fetcher, split out missing sections and
			// deliver complete ones
			var (
				sections = make([]uint64, 0, len(result.Sections))
				bitsets  = make([][]byte, 0, len(result.Bitsets))
				missing  = make([]uint64, 0, len(result.Sections))
			)
			for i, bitset := range result.Bitsets {
				if len(bitset) == 0 {
					missing = append(missing, result.Sections[i])
					continue
				}
				sections = append(sections, result.Sections[i])
				bitsets = append(bitsets, bitset)
			}
			m.schedulers[result.Bit].deliver(sections, bitsets)
			allocs--

			// Reschedule missing sections and allocate bit if newly available
			if len(missing) > 0 {
				queue := requests[result.Bit]
				for _, section := range missing {
					index := sort.Search(len(queue), func(i int) bool { return queue[i] >= section })
					queue = append(queue[:index], append([]uint64{section}, queue[index:]...)...)
				}
				requests[result.Bit] = queue

				if len(queue) == len(missing) {
					assign(result.Bit)
				}
			}

			// End the session when all pending deliveries have arrived.
			if shutdown == nil && allocs == 0 {
				return
			}
		}
	}
}

// MatcherSession is returned by a started matcher to be used as a terminator
// for the actively running matching operation.
type MatcherSession struct {
	matcher *Matcher

	closer sync.Once     // Sync object to ensure we only ever close once
	quit   chan struct{} // Quit channel to request pipeline termination

	ctx     context.Context // Context used by the light client to abort filtering
	err     error           // Global error to track retrieval failures deep in the chain
	errLock sync.Mutex

	pend sync.WaitGroup
}

// Close stops the matching process and waits for all subprocesses to terminate
// before returning. The timeout may be used for graceful shutdown, allowing the
// currently running retrievals to complete before this time.
func (s *MatcherSession) Close() {
	s.closer.Do(func() {
		// Signal termination and wait for all goroutines to tear down
		close(s.quit)
		s.pend.Wait()
	})
}

// Error returns any failure encountered during the matching session.
func (s *MatcherSession) Error() error {
	s.errLock.Lock()
	defer s.errLock.Unlock()

	return s.err
}

// allocateRetrieval assigns a bloom bit index to a client process that can either
// immediately request and fetch the section contents assigned to this bit or wait
// a little while for more sections to be requested.
func (s *MatcherSession) allocateRetrieval() (uint, bool) {
	fetcher := make(chan uint)

	select {
	case <-s.quit:
		return 0, false
	case s.matcher.retrievers <- fetcher:
		bit, ok := <-fetcher
		return bit, ok
	}
}

// pendingSections returns the number of pending section retrievals belonging to
// the given bloom bit index.
func (s *MatcherSession) pendingSections(bit uint) int {
	fetcher := make(chan uint)

	select {
	case <-s.quit:
		return 0
	case s.matcher.counters <- fetcher:
		fetcher <- bit
		return int(<-fetcher)
	}
}

// allocateSections assigns all or part of an already allocated bit-task queue
// to the requesting process.
func (s *MatcherSession) allocateSections(bit uint, count int) []uint64 {
	fetcher := make(chan *Retrieval)

	select {
	case <-s.quit:
		return nil
	case s.matcher.retrievals <- fetcher:
		task := &Retrieval{
			Bit:      bit,
			Sections: make([]uint64, count),
		}
		fetcher <- task
		return (<-fetcher).Sections
	}
}

// deliverSections delivers a batch of section bit-vectors for a specific bloom
// bit index to be injected into the processing pipeline.
func (s *MatcherSession) deliverSections(bit uint, sections []uint64, bitsets [][]byte) {
	s.matcher.deliveries <- &Retrieval{Bit: bit, Sections: sections, Bitsets: bitsets}
}

// Multiplex polls the matcher session for retrieval tasks and multiplexes it into
// the requested retrieval queue to be serviced together with other sessions.
//
// This method will block for the lifetime of the session. Even after termination
// of the session, any request in-flight need to be responded to! Empty responses
// are fine though in that case.
func (s *MatcherSession) Multiplex(batch int, wait time.Duration, mux chan chan *Retrieval) {
	for {
		// Allocate a new bloom bit index to retrieve data for, stopping when done
		bit, ok := s.allocateRetrieval()
		if !ok {
			return
		}
		// Bit allocated, throttle a bit if we're below our batch limit
		if s.pendingSections(bit) < batch {
			select {
			case <-s.quit:
				// Session terminating, we can't meaningfully service, abort
				s.allocateSections(bit, 0)
				s.deliverSections(bit, []uint64{}, [][]byte{})
				return

			case <-time.After(wait):
				// Throttling up, fetch whatever's available
			}
		}
		// Allocate as much as we can handle and request servicing
		sections := s.allocateSections(bit, batch)
		request := make(chan *Retrieval)

		select {
		case <-s.quit:
			// Session terminating, we can't meaningfully service, abort
			s.deliverSections(bit, sections, make([][]byte, len(sections)))
			return

		case mux <- request:
			// Retrieval accepted, something must arrive before we're aborting
			request <- &Retrieval{Bit: bit, Sections: sections, Context: s.ctx}

			result := <-request
			if result.Error != nil {
				s.errLock.Lock()
				s.err = result.Error
				s.errLock.Unlock()
				s.Close()
			}
			s.deliverSections(result.Bit, result.Sections, result.Bitsets)
		}
	}
}
