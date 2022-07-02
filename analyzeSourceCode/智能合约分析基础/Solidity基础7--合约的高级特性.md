## 合约的高级特性

### 继承

继承的机制和python的非常相似，但是存在差异。一般而言使用过 C++, 基本已经掌握。

当合约继承其他的合约时，只会在区块链上生成一个合约，**所有相关的合约都会编译进这个合约，调用机制和写在一个合约上一致。**

**继承时，全局变量无法覆盖，如果出现可见的同名变量会编译错误**。通过例子来体会细节，重点理解语法，而不是程序逻辑。

```js
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;


contract Owned {
    constructor() { owner = payable(msg.sender); }//构造函数中的msg.sender 是部署者
    address payable owner;
}

	// `is` 是继承的关键词. 子合约可以接受父合约所有非 private 的东西. 
contract Destructible is Owned {
    //  `virtual` 表示函数可以被重写
    function destroy() virtual public {
        if (msg.sender == owner) selfdestruct(owner);//只有调用函数的人是部署者，才能执行自毁操作
    }
}

// abstract用于提取合约的 接口，重写后实现更多的功能
abstract contract Config {
    function lookup(uint id) public virtual returns (address adr);
}

abstract contract NameReg {
    function register(bytes32 name) public virtual;
    function unregister() public virtual;
}

// 允许从多个合约继承. 
contract Named is Owned, Destructible {
    constructor(bytes32 name) {
        Config config = Config(0xD5f9D8D94886E70b06E474c3fB14Fd43E2f23970);//从地址创建 满足接口的 Condig 合约实例，用于调用
        NameReg(config.lookup(1)).register(name);//这里并未重写 lookup函数，因此返回值都是默认零值，这里创建0地址上的NameReg合约实例，然后注册管理者
    }

    // 将重写的函数需要使用overridden的标识，并且被重写的函数之前有virtual标识。
    //注意重写函数的名字，参数以及返回值类型都不能变。
    function destroy() public virtual override {
        if (msg.sender == owner) {
            Config config = Config(0xD5f9D8D94886E70b06E474c3fB14Fd43E2f23970);
            NameReg(config.lookup(1)).unregister();
            Destructible.destroy();
        }
    }
}

// 如果父合约有构造函数，则需要填上参数。
contract PriceFeed is Owned, Destructible, Named("GoldFeed") {
    function updateInfo(uint newInfo) public {
        if (msg.sender == owner) info = newInfo;
    }

    // 如果从多个合约继承了同名的可重写函数，需要在override后面指明所有同名函数所在的合约。
    function destroy() public override(Destructible, Named) { Named.destroy(); }
    function get() public view returns(uint r) { return info; }

    uint info;
}
```



但是，继承是从右到左深度优先搜索来寻找同名函数（搜索的顺序是按 ”辈分“ 从小到大，而且继承多个合约时也要按着从右到左的顺序填上，如下图继承链是 D, C, B, A），一旦找到同名函数就停止，不会执行后面重复出现的重名函数。所以如果继承了多个合约，希望把上一级父合约的同名函数都执行一遍，就需要 `super` 关键词。

```js
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;

/* Inheritance tree
   A
 /  \
B   C
 \ /
  D
*/

contract A {
    event Log(string message);

    function foo() public virtual {
        emit Log("A.foo called");
    }

    function bar() public virtual {
        emit Log("A.bar called");
    }
}

contract B is A {
    function foo() public virtual override {
        emit Log("B.foo called");
        A.foo();
    }

    function bar() public virtual override {
        emit Log("B.bar called");
        super.bar();
    }
}

contract C is A {
    function foo() public virtual override {
        emit Log("C.foo called");
        A.foo();
    }

    function bar() public virtual override {
        emit Log("C.bar called");
        super.bar();
    }
}

contract D is B, C {
    // Try:
    // - Call D.foo and check the transaction logs.
    //   Although D inherits A, B and C, it only called C and then A.
    // - Call D.bar and check the transaction logs
    //   D called C, then B, and finally A.
    //   Although super was called twice (by B and C) it only called A once.

    function foo() public override(B, C) {
        super.foo();
    }

    function bar() public override(B, C) {
        super.bar();
    }
}

```

更多的介绍请见[官方文档](https://docs.soliditylang.org/en/latest/contracts.html#inheritance)。



### 函数重写

父合约中被标记为`virtual`的非private函数可以在子合约中用`override`重写。

重写可以改变函数的标识符，规则如下：

- 可见性只能单向从 `external` 更改为 `public。`
- `nonpayable` 可以被 `view` 和 `pure` 覆盖。
- `view` 可以被 `pure` 覆盖。 
- `payable` 不可被覆盖。

如果有多个父合约有相同定义的函数， `override` 关键字后必须指定所有父合约的名字，且这些父合约没有被继承链上的其他合约重写。

接口会自动作为 `virtual` 。

注意：特殊的，如果 `external` 函数的参数和返回值和 `public` 全局变量一致的话，可以把函数重写全局变量。

```js
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.6.0 <0.9.0;

contract A
{
    function f() external view virtual returns(uint) { return 5; }
}

contract B is A
{
    uint public override f;
}
```

**注意：**函数修饰器也支持重写，且和函数重写规则一致。

```js
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.6.0 <0.9.0;

contract Base
{
    modifier foo() virtual {_;}
}

contract Inherited is Base
{
    modifier foo() override {_;}
}
```

### 抽象合约

如果合约至少有一个函数没有完成 (例如：`function foo(address) external returns (address);`)，则该合约会被视为抽象合约，需要用 `abstract` 标明。

```js
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.6.0 <0.9.0;

abstract contract Feline {
    function utterance() public pure virtual returns (bytes32);
}

contract Cat is Feline {
    function utterance() public pure override returns (bytes32) { return "miaow"; }
}
```

如果子合约没有重写父合约中所有未完成的函数，那么子合约也需要标注`abstract`

注意：声明函数类型的变量和未实现的函数的不同：

```js
function(address) external returns (address) foo;//函数类型变量
function foo(address) external returns (address);//抽象合约的函数
```

抽象合约可以将定义合约和实现合约的过程分离开，具有更佳的可拓展性。

### 接口

接口和抽象合约的作用很类似，但是它的每一个函数都没有实现，而且不可以作为其他合约的子合约，只能作为父合约被继承。

接口中所有的函数必须是`external`，且**不包含构造函数和全局变量**。接口的所有函数都会隐式标记为`external`，可以重写。多次重写的规则和多继承的规则和一般函数重写规则一致。

```js
pragma solidity >=0.6.2 <0.9.0;

interface Token {
    enum TokenType { Fungible, NonFungible }
    struct Coin { string obverse; string reverse; }
    function transfer(address recipient, uint amount) external;
}
```

### 库

库与合约类似，但是它们只在某个合约地址部署一次，并且通过 EVM 的`DELEGATECALL` （为了实现上下文更改）来实现复用。

当库中的函数被调用时，它的代码在当前合约的上下文中执行，并且只可以访问调用时显式提供的调用合约的状态变量。库本身没有状态记录（如 全局变量）。

如果库被继承的话，库函数在子合约是可见的，也可以直接使用，和普通的继承相同（属于库的内部调用方式）。为了改变状态，**内部的库（即不是通过地址引入的库）所有`data a rea` 的传参需要都是传递一个引用**（库函数使用`storage`标识），在EVM中，编译也是直接把库包含进调用合约。

```js
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.6.0 <0.9.0;

struct Data {
    mapping(uint => bool) flags;
}

library Set {
    // 注意到这里使用的是storage引用类型  
    function insert(Data storage self, uint value)
        public
        returns (bool)
    {
        if (self.flags[value])
            return false; // 如果已经存在停止插入
        self.flags[value] = true;
        return true;
    }

    function remove(Data storage self, uint value)
        public
        returns (bool)
    {
        if (!self.flags[value]) 
            return false; // 如果不存在就比用移除
        self.flags[value] = false;
        return true;
    }

    function contains(Data storage self, uint value)
        public
        view
        returns (bool)
    {
        return self.flags[value];
    }
}


contract C {
    Data knownValues;

    function register(uint value) public {
        require(Set.insert(knownValues, value));
    }
    // In this contract, we can also directly access knownValues.flags, if we want.
}
```

库具有以下特性：

- 没有状态变量
- 不能够继承或被继承
- 不能接收以太币
- 不可以被销毁

### Using For

 `using A for B;` 可用于附加库函数（从库 `A`）到任何类型（`B`）

**`using A for *;` 的效果是，库 `A` 中的函数被附加在任意的类型上，这个类型可以使用A内的函数。**

```js
pragma solidity >=0.6.0 <0.9.0;

// 这是和之前一样的代码，只是没有注释。
struct Data { mapping(uint => bool) flags; }

library Set {

  function insert(Data storage self, uint value)
      public
      returns (bool)
  {
      if (self.flags[value])
        return false; // 已经存在
      self.flags[value] = true;
      return true;
  }

  function remove(Data storage self, uint value)
      public
      returns (bool)
  {
      if (!self.flags[value])
          return false; // 不存在
      self.flags[value] = false;
      return true;
  }

  function contains(Data storage self, uint value)
      public
      view
      returns (bool)
  {
      return self.flags[value];
  }
}

contract C {
    using Set for Data; // 这里是关键的修改
    Data knownValues;

    function register(uint value) public {
        // Here, all variables of type Data have
        // corresponding member functions.
        // The following function call is identical to
        // `Set.insert(knownValues, value)`
        // 这里， Data 类型的所有变量都有与之相对应的成员函数。
        // 下面的函数调用和 `Set.insert(knownValues, value)` 的效果完全相同。
        require(knownValues.insert(value));
    }
}
```

引用存储变量或者 internal 库调用 是唯一不会发生拷贝的情况。