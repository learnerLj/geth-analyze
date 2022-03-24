## EIP-2124

因为 `core/forkId` 包是 EIP-2124 的实现，因此我们先了解 [EIP-2124](https://eips.ethereum.org/EIPS/eip-2124) 的内容。

### 目的

> **记住以太坊是向后兼容！**

以太坊节点之前寻找其他节点的方式是：随机的选择可以连接的节点，再去判断节点是否对自己有用。但是其他的节点可能是主网节点、测试网节点、私网节点、共识机制不一致的节点，这样“盲目”的寻找会浪费资源。

于是这个提案希望在连接前，其他节点可以发送自己的链配置，实现精确连接到有用节点。这个传递的信息在提案中被称作 `fork identifier`，它实现如下功能：

- 如果两个节点在不同的网络上，他们不应该考虑连接。
- 如果硬分叉通过，升级的节点应该拒绝未升级的节点。
- 如果两条链有相同的创世区块配置和链配置，但没有分叉（ETH / ETC），它们应该拒绝连接。

这个提案并没有完整地处理分叉的问题，例如对于分叉 A、B、C，它没有单独处理每个分叉。

### 具体实现方式

- **`FORK_HASH`**: 创世配置哈希和分叉区块高度的校验和，4 字节。
  - 分叉的区块高度按照顺序，计算 CRC32 校验和。
  - 分叉的区块高度视作 `uint64`，计算校验和时按照大端序。
  - 如果创世配置中不包括 Frontier 版本，就不考虑分叉带来的影响。
- **`FORK_NEXT`**: 下一个将分叉的区块高度，如果未知则为 0。

例如主网的 `ForK_HASH` 的计算：

- forkhash₀ = `0xfc64ec04` (Genesis) = `CRC32(<genesis-hash>)`
- forkhash₁ = `0x97c2c34c` (Homestead) = `CRC32(<genesis-hash> || uint64(1150000))`
- forkhash₂ = `0x91d1f948` (DAO fork) = `CRC32(<genesis-hash> || uint64(1150000) || uint64(1920000))`

然后 `fork indentifier` 定义为 `RLP([FORK_HASH, FORK_NEXT])`。提案后面给出了详细的例子，这里不再讨论。

### 验证规则

- 本地和远程节点的 `FORK_HASH` 匹配，比较当前的区块高度和 `FORK_NEXT`。这说明当前的两个节点是兼容的，未来的某次分叉后可能不兼容。
  - 如果存在远程节点宣布的分叉区块还没有传递给本地节点，但本地节点已经接收了分叉的区块，则断开的当前连接。（因为本地升级了，连接着的远程节点还没有，这就不兼容）
  - 没有远程节点宣布分叉区块，而且本地节点也没有接收到分叉区块，则继续连接。
- 如果远程节点的 `FORK_HASH` 是当前节点的子集，并且远程节点的 `FORK_NEXT` 与本地节点将接收的分叉高度相同，则连接。
  - 远程节点没有实现某些分叉，但是它可以获取信息，虽然后面可能会“渐行渐远”。
- 如果远程节点的 `FORK_HASH` 是当前节点的超集，并且本地节点可以在将来的分叉中逐渐与远程节点的 `FORK_HASH` 相同。连接。
- 如果不是前面提到的情况，则拒绝所有连接。



## 源码实现

### 兼容性错误和生成 ID 的依据

首先定义远程节点和本地节点之间不兼容的情况。一种是远程节点的检验和是本地的子集，说明远程节点还没升级，不兼容。另一种是远程节点的校验和不匹配，说明它的分叉顺序、链配置、创世区块配置与本地不兼容。

```go
var (
	// ErrRemoteStale is returned by the validator if a remote fork checksum is a
	// subset of our already applied forks, but the announced next fork block is
	// not on our already passed chain.
	ErrRemoteStale = errors.New("remote needs update")

	// ErrLocalIncompatibleOrStale is returned by the validator if a remote fork
	// checksum does not match any local checksum variation, signalling that the
	// two chains have diverged in the past at some point (possibly at genesis).
	ErrLocalIncompatibleOrStale = errors.New("local incompatible or needs update")
)
```

为了传递信息，定义了链配置、创世区块配置和当前区块头的接口。然后定义了 `ID` 包括了前面提到的 `FORK_HASH` 和 `FORK_NEXT`，最后定义了判断是否兼容的函数 `Filter`。

```go
// Blockchain defines all necessary method to build a forkID.
type Blockchain interface {
	// Config retrieves the chain's fork configuration.
	Config() *params.ChainConfig

	// Genesis retrieves the chain's genesis block.
	Genesis() *types.Block

	// CurrentHeader retrieves the current head header of the canonical chain.
	CurrentHeader() *types.Header
}

// ID is a fork identifier as defined by EIP-2124.
type ID struct {
	Hash [4]byte // CRC32 checksum of the genesis block and passed fork block numbers
	Next uint64  // Block number of the next upcoming fork, or 0 if no forks are known
}

// Filter is a fork id filter to validate a remotely advertised ID.
type Filter func(id ID) error
```

### 创建 ID 标识

节点根据对方的 ID 的内容，选择是否连接。从链配置、创世区块哈希、当前分叉区块高度获取节点的标识

```go
// NewID calculates the Ethereum fork ID from the chain config, genesis hash, and head.
func NewID(config *params.ChainConfig, genesis common.Hash, head uint64) ID {
	// Calculate the starting checksum from the genesis hash
	hash := crc32.ChecksumIEEE(genesis[:]) //计算创世区块校验和

	// Calculate the current fork checksum and the next fork block
	var next uint64
	for _, fork := range gatherForks(config) {
		if fork <= head {
			// Fork already passed, checksum the previous hash and the fork number
			hash = checksumUpdate(hash, fork) //将之前的分叉逐个添加进 fork_hash
			continue
		}
		next = fork //如果超过了当前的分叉的区块高度，那么这是即将迎来的分叉。否则 next 为 0
		break
	}
	return ID{Hash: checksumToBytes(hash), Next: next}
}
```

### 创建过滤器

首先根据本地节点的配置，创建过滤器。然后输入远程节点的 ID，通过检查过滤器是否抛出不兼容的错误，就可以知道是否连接远程节点。

过滤规则详见 EIP-2124，本文前面提到过。

```go
// newFilter is the internal version of NewFilter, taking closures as its arguments
// instead of a chain. The reason is to allow testing it without having to simulate
// an entire blockchain.
func newFilter(config *params.ChainConfig, genesis common.Hash, headfn func() uint64) Filter {
	// Calculate the all the valid fork hash and fork next combos
	var (
		//各种分叉的区块高度
		forks = gatherForks(config)
		//每个分叉的对应的累积校验和
		sums = make([][4]byte, len(forks)+1) // 0th is the genesis
	)
	hash := crc32.ChecksumIEEE(genesis[:])
	sums[0] = checksumToBytes(hash)

	//整合校验和
	for i, fork := range forks {
		hash = checksumUpdate(hash, fork)
		sums[i+1] = checksumToBytes(hash)
	}

	//最后一个位置作为 "哨兵"，用于方便处理

	// Add two sentries to simplify the fork checks and don't require special
	// casing the last one.
	forks = append(forks, math.MaxUint64) // Last fork will never be passed

	// Create a validator that will filter out incompatible chains
	return func(id ID) error {
		// Run the fork checksum validation ruleset:
		//   1. If local and remote FORK_CSUM matches, compare local head to FORK_NEXT.
		//        The two nodes are in the same fork state currently. They might know
		//        of differing future forks, but that's not relevant until the fork
		//        triggers (might be postponed, nodes might be updated to match).
		//      1a. A remotely announced but remotely not passed block is already passed
		//          locally, disconnect, since the chains are incompatible.
		//      1b. No remotely announced fork; or not yet passed locally, connect.
		//   2. If the remote FORK_CSUM is a subset of the local past forks and the
		//      remote FORK_NEXT matches with the locally following fork block number,
		//      connect.
		//        Remote node is currently syncing. It might eventually diverge from
		//        us, but at this current point in time we don't have enough information.
		//   3. If the remote FORK_CSUM is a superset of the local past forks and can
		//      be completed with locally known future forks, connect.
		//        Local node is currently syncing. It might eventually diverge from
		//        the remote, but at this current point in time we don't have enough
		//        information.
		//   4. Reject in all other cases.
		head := headfn()
		for i, fork := range forks {
			//如果当前区块高度超过了某个分叉，就继续往后检查。前面设置的哨兵的作用就是在这里跳过这个检查

			// If our head is beyond this fork, continue to the next (we have a dummy
			// fork of maxuint64 as the last item to always fail this check eventually).
			if head >= fork {
				continue
			}
			//达到了本地节点的下一个分叉区块高度

			// Found the first unpassed fork block, check if our current state matches
			// the remote checksum (rule #1).
			if sums[i] == id.Hash { //如果匹配
				// Fork checksum matched, check if a remote future fork block already passed
				// locally without the local node being aware of it (rule #1a).
				if id.Next > 0 && head >= id.Next {
					//但是当前区块高度高于远程节点的下一个分叉高度，那么不兼容
					return ErrLocalIncompatibleOrStale
				}
				// Haven't passed locally a remote-only fork, accept the connection (rule #1b).
				return nil
			}

			//开始部分匹配。远程节点与本地节点处于不同的分叉状态，这是要求远程节点在本地节点“之前”,
			//表现出来是远程节点的校验和是本地节点的子集

			// The local and remote nodes are in different forks currently, check if the
			// remote checksum is a subset of our local forks (rule #2).
			for j := 0; j < i; j++ {
				//找子集的匹配部分，并且要求远程节点接着的分叉与本地节点对应
				if sums[j] == id.Hash {
					// Remote checksum is a subset, validate based on the announced next fork
					if forks[j] != id.Next {
						return ErrRemoteStale
					}
					return nil
				}
			}

			//如果远程节点包括了当前节点的所有分叉，那么可以连接，告诉远程节点，当前节点没有完成同步

			// Remote chain is not a subset of our local one, check if it's a superset by
			// any chance, signalling that we're simply out of sync (rule #3).
			for j := i + 1; j < len(sums); j++ {
				if sums[j] == id.Hash {
					// Yay, remote checksum is a superset, ignore upcoming forks
					return nil
				}
			}
			// No exact, subset or superset match. We are on differing chains, reject.
			return ErrLocalIncompatibleOrStale
		}
		log.Error("Impossible fork ID validation", "id", id)
		return nil // Something's very wrong, accept rather than reject
	}
}
```

### 提取分叉高度

从上面过滤器的规则可以看出，它主要根据分叉的高度来判断是否兼容。以下的函数用于从链配置中提取分叉高度。

```go
// gatherForks gathers all the known forks and creates a sorted list out of them.
func gatherForks(config *params.ChainConfig) []uint64 {
   // Gather all the fork block numbers via reflection
   kind := reflect.TypeOf(params.ChainConfig{})
   conf := reflect.ValueOf(config).Elem()

   var forks []uint64
   for i := 0; i < kind.NumField(); i++ {
      // Fetch the next field and skip non-fork rules
      field := kind.Field(i)

      //处理链配置中的分叉区块，因为它们末尾都是 Block，而且都是 bigInt 类型
      if !strings.HasSuffix(field.Name, "Block") {
         continue
      }
      if field.Type != reflect.TypeOf(new(big.Int)) {
         continue
      }
      // Extract the fork rule block number and aggregate it
      rule := conf.Field(i).Interface().(*big.Int)
      if rule != nil {
         forks = append(forks, rule.Uint64())
      }
   }

   //数据量很小，冒泡排序也不差，区块高度需要按升序，表达分叉顺序

   // Sort the fork block numbers to permit chronological XOR
   for i := 0; i < len(forks); i++ {
      for j := i + 1; j < len(forks); j++ {
         if forks[i] > forks[j] {
            forks[i], forks[j] = forks[j], forks[i]
         }
      }
   }

   //处理同一个区块高度多个分叉的情况，删除前面的重复的分叉。这种情况几乎不会发生。
   //如A、B、C 的分叉高度都是 1000，那么删除 A、B对应的区块高度。

   // Deduplicate block numbers applying multiple forks
   for i := 1; i < len(forks); i++ {
      if forks[i] == forks[i-1] {
         forks = append(forks[:i], forks[i+1:]...)
         i--
      }
   }

   //跳过高度为 0 的分叉，因为它是写在创世区块的配置里。

   // Skip any forks in block 0, that's the genesis ruleset
   if len(forks) > 0 && forks[0] == 0 {
      forks = forks[1:]
   }
   return forks
}
```