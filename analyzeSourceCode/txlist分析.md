# TxList-Analysis

> 本文旨在分析清楚tx_list.go中这个工具包里面的重要代码

## 堆排序

以下为`tx_list.go`中的`heap.Interface`的全部实现代码，非常标准，和默认的一样；

```go
//heap的整个实现过程
type nonceHeap []uint64

func (h nonceHeap) Len() int           { return len(h) }
func (h nonceHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h nonceHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *nonceHeap) Push(x interface{}) {
	*h = append(*h, x.(uint64))
}

func (h *nonceHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
```

以下展示`heap.Interface`的结构：

```go
type Interface interface {
    sort.Interface
    Push(x interface{}) // add x as element Len()
    Pop() interface{}   // remove and return element Len() - 1.
}
```

其中`sort.Interface`这个接口里包含`Len() Less() Swap()`这三个方法，也就是对应上面的前三个方法；

加上后面的`Push() Pop()`两个方法，也就是我们实现了`heap.Interface`这个接口；

然后我们就可以使用`heap`包里面的相关功能性函数（因为他们的参数要求基本上都包含`heap.Interface`这个接口)，`heap`包代码量非常少，算上注释才一百多行，很容易也推荐看完；

在这里附上一段试验的源码以供参考：

```go
package main

import (
	"container/heap"
	"fmt"
	"sort"
)

type intHeap []int

func (a intHeap) Len() int           { return len(a) }
func (a intHeap) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a intHeap) Less(i, j int) bool { return a[i] < a[j] }

func (s *intHeap) Push(x interface{}) {
	*s = append(*s, x.(int))
}

func (s *intHeap) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}
func main() {
	//必须要是初始化为指针类型的 这样才算是实现了该接口 因为上面的参数有指针类型的
	res := &intHeap{2,1,4,6,3}
	//s:=heap.Interface(res)
	heap.Init(res)
	fmt.Println(res)
	heap.Push(res,3)
	fmt.Println(res)
	fmt.Println((*res)[4])
	sort.Sort(res)
	fmt.Println(res)
}
```

运行结果如下：

<img src="C:\Users\Administrator\AppData\Roaming\Typora\typora-user-images\image-20220311150102351.png" alt="image-20220311150102351" style="zoom:50%;" />

---

## 函数功能解析

### txSortedMap

具体的结构如下：

```go
type txSortedMap struct {
	//建立一个nonce->transaction的map
	items map[uint64]*types.Transaction // Hash map storing the transaction data
	index *nonceHeap                    // Heap of nonces of all the stored transactions (non-strict mode)
	cache types.Transactions            // Cache of the transactions already sorted
}
```

以下为其对应的所有方法：

1. `newTxSortedMap()`  进行初始化并返回初始化后的`*txSortedMap`

2. `Get(nonce uint64)`获取指定`nonce`的交易并返回该笔交易

3. `Put(tx *types.transaction)`将该笔交易该笔交易添加到`txSortedMap`中，无论之前是否存在

4. `Foward(threshold uint64)`将低于这个门槛的`nonce`的交易全部剔除

5. `reheap()`根据当前的`map`重新进行`nonceheap`的排序

6. `filter(filter func(*types.Transaction) bool)`其中的参数`filter func(*types.Transaction) bool)`

   它的源码是`func(tx *types.Transaction) bool { return tx.Nonce() > lowest }`

   我们可以发现该函数的作用其实为了过滤`nonce`小于最低要求的交易，而`Filter(filter func(*types.Transaction) bool)`调用了以上的函数，所以功能差不多

7. `Cap(threshold int)`如果该`map`中的交易数量超过了限制，就删除最高`nonce`的交易直至数量达到要求，并返回删除掉的`drops`

8. `Remove(nonce uint64)`删除成功则返回`true`,没有找到就返回`false`

9. `Ready(start uint64)`准备好`nonce`高于`start`并且`连续`的交易

10. `Len()`返回`map`的大小

11. `Flatten()`获取全部的交易，`flatten()`将全部按照`nonce`排序好的交易进行缓存

12. `LastElement()`返回`cashe`中`nonce值`最高的交易

---

### txList

具体结构如下：

```go
type txList struct {
    //nonce是否严格连续
	strict bool         // Whether nonces are strictly continuous or not
	//nonce
    txs    *txSortedMap // Heap indexed sorted hash map of the transactions

	costcap *big.Int // Price of the highest costing transaction (reset only if exceeds balance)
	gascap  uint64   // Gas limit of the highest spending transaction (reset only if exceeds block limit)
}
```

重要函数如下：

1. `Overlaps(tx *types.Transaction)`若是已有这笔交易就返回`true`，否则返回`false`
2. `Add(tx *types.Transaction,priceBump uint64)`若是已有这笔交易就尝试加入，不存在就直接加入，返回交易`true old`
3. `Filter(costLimit *big.Int,gasLimit uint64)`过滤掉拥有过高的`cost`或者`gas`的交易，同时过滤掉后面`nonce`不连续的交易
4. `Remove(tx *types.Transaction)`尝试移除掉指定的交易并移除后面`nonce`值不连续的交易

---

### priceHeap

具体结构如下：

```go
type priceHeap struct {
	baseFee *big.Int // heap should always be re-sorted after baseFee is changed
	list    []*types.Transaction
}
```

实现了`heap.Interface`这个接口，排序方式是`先比较gasFee,再之tipFee`

---

### txPricedList

具体结构如下：

```go
type txPricedList struct {
    //一个过时的计数器
	stales int64

	all              *txLookup  // Pointer to the map of all transactions
	urgent, floating priceHeap  // Heaps of prices of all the stored **remote** transactions
	reheapMu         sync.Mutex // Mutex asserts that only one routine is reheaping the list
}
```

首先提一下，本部分源码中大量使用了原子操作，`Go语言`中提供的原子操作都是非侵入式的，在标准库代码包`sync/atomic`中提供了相关的原子函数,具体功能如下：

原子操作即是进行过程中不能被中断的操作，针对某个值的原子操作在被进行的过程中，CPU 绝不会再去进行其他的针对该值的操作。为了实现这样的严谨性，原子操作仅会由一个独立的CPU指令代表和完成。原子操作是无锁的，常常直接通过CPU指令直接实现。 事实上，其它同步技术的实现常常依赖于原子操作。





现在理解不来为何 `priceHeap`里要有`urgent 和 floating`这两个 ，（注意：应该是新增的，网上找不到资料；）

