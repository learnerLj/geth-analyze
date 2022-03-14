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
	"sync"
)

//以太坊的布隆过滤器总共有 2048 位，以太坊会把若干个区块分成段，段作为检索的基本单位，4096 个区块为一段。
//checkpoint 和 时间检索也是这样。
//下面是一个检索请求，表示在特定的一段 section 中希望匹配 2048 位的过滤器中的哪一位。

// request represents a bloom retrieval task to prioritize and pull from the local
// database or remotely from the network.
type request struct {
	section uint64 // Section index to retrieve the a bit-vector from
	bit     uint   // Bit index within the section to retrieve the vector of
}

//和上面的 request 对应，表示被检索的 bit 向量的状态（即上面提到的 2048 个比特的集合）
//cached 缓存检索结果，用于去重

// response represents the state of a requested bit-vector through a scheduler.
type response struct {
	cached []byte        // Cached bits to dedup multiple requests
	done   chan struct{} // Channel to allow waiting for completion
}

//scheduler 是布隆过滤器的一次任务，在指定的区块中返回结果。
//主要包含了查找哪一位，以及每个区块高度的结果。

// scheduler handles the scheduling of bloom-filter retrieval operations for
// entire section-batches belonging to a single bloom bit. Beside scheduling the
// retrieval operations, this struct also deduplicates the requests and caches
// the results to minimize network/database overhead even in complex filtering
// scenarios.
type scheduler struct {
	bit       uint                 // Index of the bit in the bloom filter this scheduler is responsible for
	responses map[uint64]*response // Currently pending retrieval requests or already cached responses
	lock      sync.Mutex           // Lock protecting the responses from concurrent access
}

// newScheduler creates a new bloom-filter retrieval scheduler for a specific
// bit index.
func newScheduler(idx uint) *scheduler {
	return &scheduler{
		bit:       idx,
		responses: make(map[uint64]*response),
	}
}

//并发地分发检索任务，sections chan uint64 获取需要检索的某一段，
//dist 表示需要分发到 schedule 的请求。
//done 表示这次任务的结果，最后是退出的标识
//
//最终的效果是

// run creates a retrieval pipeline, receiving section indexes from sections and
// returning the results in the same order through the done channel. Concurrent
// runs of the same scheduler are allowed, leading to retrieval task deduplication.
func (s *scheduler) run(sections chan uint64, dist chan *request, done chan []byte, quit chan struct{}, wg *sync.WaitGroup) {
	//请求和回应之间的缓冲通道，用于阻塞和控制。

	// Create a forwarder channel between requests and responses of the same size as
	// the distribution channel (since that will block the pipeline anyway).
	pend := make(chan uint64, cap(dist))

	// Start the pipeline schedulers to forward between user -> distributor -> user
	wg.Add(2)
	go s.scheduleRequests(sections, dist, pend, quit, wg)
	go s.scheduleDeliveries(pend, done, quit, wg)
}

// reset cleans up any leftovers from previous runs. This is required before a
// restart to ensure the no previously requested but never delivered state will
// cause a lockup.
func (s *scheduler) reset() {
	s.lock.Lock()
	defer s.lock.Unlock()

	for section, res := range s.responses {
		if res.cached == nil {
			delete(s.responses, section)
		}
	}
}

// scheduleRequests reads section retrieval requests from the input channel,
// deduplicates the stream and pushes unique retrieval tasks into the distribution
// channel for a database or network layer to honour.
func (s *scheduler) scheduleRequests(reqs chan uint64, dist chan *request, pend chan uint64, quit chan struct{}, wg *sync.WaitGroup) {
	// Clean up the goroutine and pipeline when done
	defer wg.Done()
	defer close(pend)

	// Keep reading and scheduling section requests
	for {
		select {
		case <-quit: //收到退出信号就返回
			return

		case section, ok := <-reqs: //需要检索的段
			// New section retrieval requested
			if !ok {
				return
			}
			// Deduplicate retrieval requests
			unique := false //因为并发执行时可能一已经进入了协程，这里用于去重

			//阻塞其他协程，直到这一部分完成，Unlock
			s.lock.Lock()
			//如果请求为空，那么设置为已完成，避免重复执行
			if s.responses[section] == nil {
				s.responses[section] = &response{
					done: make(chan struct{}),
				}
				unique = true
			}
			s.lock.Unlock()

			// Schedule the section for retrieval and notify the deliverer to expect this section
			if unique {
				//如果对应的请求为空，但是还没有在其他协程结束，那么把结果分发出去
				select {
				case <-quit:
					return
				case dist <- &request{bit: s.bit, section: section}:
				}
			}
			//如果还没有结束，那么给 pend 赋值，pend 进入阻塞状态。
			select {
			case <-quit:
				return
			case pend <- section:
			}
		}
	}
}

//用于分发上文执行但是未完成的任务，写入 schedule 中每个区块对应的 response

// scheduleDeliveries reads section acceptance notifications and waits for them
// to be delivered, pushing them into the output data buffer.
func (s *scheduler) scheduleDeliveries(pend chan uint64, done chan []byte, quit chan struct{}, wg *sync.WaitGroup) {
	// Clean up the goroutine and pipeline when done
	defer wg.Done()
	defer close(done)

	// Keep reading notifications and scheduling deliveries
	for {
		select {
		case <-quit:
			return
		//结束 pend 的阻塞
		case idx, ok := <-pend:
			// New section retrieval pending
			if !ok {
				return
			}
			// Wait until the request is honoured
			s.lock.Lock()
			res := s.responses[idx] //写入 response，非常关键的一步。
			s.lock.Unlock()

			select {
			case <-quit:
				return
			case <-res.done:
			}
			// Deliver the result
			select {
			case <-quit:
				return
			case done <- res.cached: //如果没有结束，将缓存的写入 done
			}
		}
	}
}

// deliver is called by the request distributor when a reply to a request arrives.
func (s *scheduler) deliver(sections []uint64, data [][]byte) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for i, section := range sections {
		if res := s.responses[section]; res != nil && res.cached == nil { // Avoid non-requests and double deliveries
			res.cached = data[i]
			close(res.done)
		}
	}
}
