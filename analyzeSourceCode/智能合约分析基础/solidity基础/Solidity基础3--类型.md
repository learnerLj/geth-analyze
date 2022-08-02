## 类型

Solidity的值传递和引用传递有自己的规则，通过不同的存储域决定，后面详述。

**默认值**：Solidity中不存在`undefined`或`null`，每种变量都有自己的默认值，一般是“零状态”。

 **[运算符优先级](https://learnblockchain.cn/docs/solidity/cheatsheet.html#order)** 

### 布尔类型

`bool`:常量值为`true`和`false`

### 整型

`int` / `uint` ：分别表示有符号和无符号的不同位数的整型变量。 支持关键字 `uint8` 到 `uint256` （无符号，从 8 位到 256 位）以及 `int8` 到 `int256`，以 `8` 位为步长递增。 `uint` 和 `int` 分别是 `uint256` 和 `int256` 的别名。

可以用`type(x).min` `type(x).max`来获取这个类型地最小值和最大值。

### 位运算

在二进制地补码上操作，特别的 `~int256（0）== int256（-1）`

**移位：**

左移则会截断最高位；右移操作数必须是无符号地整型，否则会编译错误。

- `x<<y`相当于`x*2**y`,(其实这里体现了`**`的优先级比较高)
- 如果`x>0`:`x>>y`相当于`x/2**y`
- 如果`x<0`:`x>>y`相当于`(x+1)/2**y - 1`(如果不是整数，则向下取整)（注意：0.5.0之前是向上取整）

### 加减乘除

在 `0.8.0 `之后加入了溢出检查，值超过上限或者下限则会回滚，我们可以使用`unchecked{}`来取消检查。在此之前需要使用` OpenZepplin SafeMath`库。

**注意：**`unchecked{}` 不能替代码块的花括号，而且不支持嵌套，只对花括号内的语句有效，且对其中调用的函数无效，并且花括号内不能出现 `_`。

除 0 或者模 0 会报错。`type(int).min / (-1)`是唯一的整除向上溢出的情况。

**注意移位操作符造成的溢出并不会报错**，需要额外注意溢出问题。

**幂运算只适用于无符号的整型**，有时为了减少gas消耗，编译器会建议用`x*x*x`来代替`x**3`。

定义`0**0=1`

### 定长浮点型

由于 EVM 只支持整数运算并且需要严格控制计算资源，因此浮点数的计算的实现有一定的挑战，采用了严格限制整数位数和小数位数的方式。

`fixed` / `ufixed`：表示各种大小的有符号和无符号的定长浮点型。 在关键字 `ufixedMxN` 和 `fixedMxN` 中，`M` 表示该类型占用的位数，`N` 表示可用的小数位数。 `M` 必须能整除 8，即 8 到 256 位。 `N` 则可以是从 0 到 80 之间的任意数。 `ufixed` 和 `fixed` 分别是 `ufixed128x19` 和 `fixed128x19` 的别名。

注意：solidity还没有完全的支持定长浮点型，**只能声明，但是不可以给他赋值，也不能用它给其他变量赋值**，只可以下面那样。用的很少。

```solidity
fixed8x4 a;
```

### 地址类型

这是比较特殊的类似，其他语言没有，实际上是储存字节。

地址类型有两种，

- `address`：保存一个20字节的值（以太坊地址的大小），**不支持作为转账地址**。
- `address payable` ：**可参与转账的地址**，与 `address` 相同，不过有成员函数 `transfer` 和 `send` 。

注意：`address` 和 `address payable` 的区别是在 0.5.0 版本引入的***

**地址成员：**

地址类型有默认的成员，方便查看它的属性。

- `<address>.balance` 返回 `uint256`

  以 Wei 为单位的余额。

- `<address>.code` 返回 `bytes memory`

  地址上的字节码(可以为空)

- `<address>.codehash` (`bytes32`)

  地址上的字节码哈希

- `<address payable>.transfer(uint256 amount)`

  向该地址发送数量为 amount 的 Wei，失败时抛出异常，并且会回滚。使用固定（不可调节）的 2300 gas 的矿工费。

- `<address payable>.send(uint256 amount) returns (bool)`

  向该地址发送数量为 amount 的 Wei，失败时返回 `false`，发送 2300 gas 的矿工费用，不可调节。

  **注意**：`send`安全等级比较低，他失败时（比如因为堆栈在1024或者gas不足）不会发生异常，因此往往要检查它的返回值，或者直接用`transfer`

```solidity
// SPDX-License-Identifier: MIT
// compiler version must be greater than or equal to 0.8.3 and less than 0.9.0
pragma solidity ^0.8.3;
contract HelloWorld {
    string public greet = "Hello World!";
    address public myAddress=address(this);
    uint public myBalance = myAddress.balance;
    bytes public myCode = myAddress.code;
    bytes32 public myCodehash = myAddress.codehash;
    function getstr() public view returns (string memory){
        return greet;
    }
}
```



### 合约类型

每一个合约都有自己的类型，也可以用合约名定义其他变量，相当于创建了一个接口。

合约可以通过`address(x)`转换成`address`类型；只有可支付的合约（具有receive函数或者payable fallback函数），才可以使用`payable(address(x))`转换成`address payable`类型（0.5.0版本之后才区分`address和payable address`）

下面是示例用法：

新建两个文件放在同个文件夹下：

<img src="https://gitee.com/learnerLj/typora/raw/master/image-20220119192540950.png" alt="image-20220119192540950" style="zoom: 200%;" />

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.3;
contract HelloWorld {
    string public greet = "Hello World!";
    function getstr() public view returns (string memory){
        return greet;
    }
}
```

```solidity
// SPDX-License-Identifier: GPL-3.0
import "./Hello.sol";
pragma solidity >=0.5.0 <0.9.0;
contract CallHello {
    HelloWorld public hello;
    constructor(address _addr){
        hello = HelloWorld(_addr);
    }
    function f()public view returns(string memory){
        return hello.getstr();
    }
    function g()public view returns(address){
        return address(hello);
    }
}
```





### 枚举类型

枚举类型至少需要一个成员，且不能多于256个成员。整数类型和枚举类型只能显式转化，不能隐式转化。整数转枚举时需要在枚举值的范围内，否则会引发`Panic error`。

可以使用 `type(NameOfEnum).min` 和 `type(NameOfEnum).max` 获取这个枚举类型的最小值和最大值

```solidity
    enum ActionChoices { GoLeft, GoRight, GoStraight, SitStill }
    ActionChoices choice;
 	function setGoStraight() public {
        choice = ActionChoices.GoStraight;
    }
```



### 函数类型

函数可以作为类型，可以被赋值，而且也可以作为其他函数的参数或者返回值，这一点和 Go 语言是一致的。

```solidity
pragma solidity ^0.8.3;
contract A{

    function foo() external pure returns(uint){
    uint a =5;
    return a;
    }

    function () external returns(uint) f=this.foo;//注意，访问函数类型，一定要从 this访问
   // f=this.foo;注意无法这样赋值，只能初始化时赋值
```

函数类型实际上包括一个 20 个字节的地址和 4 个字节的函数选择器，等价于 `byte24` 类型

**函数类型的调用限制**

有两种：

- *内部（internal）* 函数类型，只能在当前合约内被调用（包括内部库函数和继承的函数），不能在合约的上下文外执行。调用内部函数时通过跳转到函数的标签片段实现。
- *外部（external）* 函数类型，由一个地址和函数签名组成，在调用时会被视作`function`类型，函数的地址后面后紧跟函数标识符一起编码成`bytes24`类型。

下面是函数的类型表示：

```
function (<parameter types>) {internal|external} [pure|constant|view|payable] [returns (<return types>)]
```

函数类型默认是内部函数，但是**在合约内定义的函数可见性必须明确声明**。**在合约内定义函数的位置时任意的，可以调用后面才定义的函数。**

**函数类型的成员**

public（或 external）函数都有下面的成员：

- `.address` 返回函数的合约地址。
- `.selector` 返回 ABI 函数选择器

注意在过去还有两个成员：`.gas(uint)` 和 `.value(uint)` 在0.6.2中弃用了，在 0.7.0 中移除了。 用 `{gas: ...}` 和 `{value: ...}` 代替。

```solidity
pragma solidity ^0.8.3;
contract A{
    
    function foo() public pure returns(uint){
    uint a =5;
    return a;
    }

    function getAddr() public view returns(address){

    return this.foo.address;
    }

    function getSekector() public pure returns(bytes4){
    return this.foo.selector;
    }
}
```

内部函数的例子：(这里采用了库函数)

```solidity
library ArrayUtils {
  // 内部函数可以在内部库函数中使用，
  // 因为它们会成为同一代码上下文的一部分
  function map(uint[] memory self, function (uint) pure returns (uint) f)
    internal
    pure
    returns (uint[] memory r)
  {
    r = new uint[](self.length);
    for (uint i = 0; i < self.length; i++) {
      r[i] = f(self[i]);
    }
  }
  function reduce(
    uint[] memory self,
    function (uint, uint) pure returns (uint) f
  )
    internal
    pure
    returns (uint r)
  {
    r = self[0];
    for (uint i = 1; i < self.length; i++) {
      r = f(r, self[i]);
    }
  }
  function range(uint length) internal pure returns (uint[] memory r) {
    r = new uint[](length);
    for (uint i = 0; i < r.length; i++) {
      r[i] = i;
    }
  }
}

contract Pyramid {
  using ArrayUtils for *;
  function pyramid(uint l) public pure returns (uint) {
    return ArrayUtils.range(l).map(square).reduce(sum);//前一个的返回值作为后一个的参数
  }
  function square(uint x) internal pure returns (uint) {
    return x * x;
  }
  function sum(uint x, uint y) internal pure returns (uint) {
    return x + y;
  }
}
```

使用外部函数的例子：(对于不习惯将函数当作类型的读者，可能会比较陌生)

```solidity
pragma solidity >=0.4.22  <0.9.0;

contract Oracle {
  struct Request {
    bytes data;
    function(uint) external callback;
  }
  Request[] private requests;
  event NewRequest(uint);
  function query(bytes memory data, function(uint) external callback) public {
    requests.push(Request(data, callback));
    emit NewRequest(requests.length - 1);
  }
  function reply(uint requestID, uint response) public {
    // 这里检查回复来自可信来源
    requests[requestID].callback(response);
  }
}

contract OracleUser {
  Oracle constant private ORACLE_CONST = Oracle(address(0x00000000219ab540356cBB839Cbe05303d7705Fa)); // known contract
  uint private exchangeRate;
  function buySomething() public {
    ORACLE_CONST.query("USD", this.oracleResponse);
  }
  function oracleResponse(uint response) public {
    require(
        msg.sender == address(ORACLE_CONST),
        "Only oracle can call this."
    );
    exchangeRate = response;
  }
}
```

### 引用类型

引用类型可以通过不同变量名来修改指向的同一个值。目前的引用类型包括：结构体、数组和映射。

在使用引用类型时，需要指明这个类型存储在哪个数据域（data area）

- memory:存储在内存里，只在函数内部使用，**函数内变量不做特殊说明为`memory`类型**。
- storage:相当于全局变量。**函数外合约内的都是`storage`类型**。
- calldata:保存函数的参数的特殊储存位置，只读，大多数时候和`memory`相似。

如果可以的话，尽可能使用`calldata` 临时存储传入函数的参数，因为它既不会复制，也不能修改，而且还可以作为函数的返回值。

### 数据的赋值

**更改位置或者类型转化是拷贝；同一位置赋值一般是引用**

- `storage`和`memory`之间的赋值或者用`calldata`对它们赋值，都是产生独立的拷贝，不修改原来的值。
- `memory`之间的赋值，是引用。
- `storage`给合约的全局变量赋值总是引用。
- 其他向`storage` 赋值是拷贝。
- 结构体里面的赋值是一个拷贝。

```solidity
pragma solidity >=0.5.0 <0.9.0;

contract C {
    uint[] x; //函数外变量都默认 storage

    // 函数内变量都是 memory.
    function f(uint[] memory memoryArray) public {
        x = memoryArray; // memory 给函数外的storage变量赋值，拷贝
        uint[] storage y = x; // storage 之间 指针传递，节省内存
        y.pop(); // 同时修改X
        delete x; // 重置X,同时修改Y
        g(x); // 函数传参时,也遵守规则，这里是传引用
        h(x); //这里传复制
    }

    function g(uint[] storage) internal pure {}
    function h(uint[] memory) public pure {}
}
```



### 数组

- 创建多维数组时，下标的用法有些不一样，`a[2][4]`表示4个子数列，每个子数列里2个元素，所以`a.length`等于4。但是访问数组时下标的顺序和大多数语言相同。

- `a[3]`，其中`a`也可以是数组，即 `a[3]` 是多维数组。

- 多维数组的声明不要求写明长度，初始化如下`    uint[][] a=[[1,2,3],[4,5,6]];`，当然也可以`    uint[5][7] a=[[1,2,3],[4,5,6]];`，不够的位置用0来补上。

- 动态数组也支持切片访问，`x[start:end]` 其中的`start` 和`end` 会隐式的转化成`uint256`类型，`start` 默认是0，`end` 默认到最后，因此可以省略其中一个。切片不能够使用数组的操作成员，但是会隐式地转换成新的数组，支持进一步地按索引访问。目前，**只有 `calldata` 的数组才支持切片。**

- 数组可以是任何类型，包括映射和结构体。但是**数组中的映射只能是storage类型**。

- `bytes.concat`函数可以把`bytes` 或者 `bytes1 ... bytes32` 拼接起来，返回一个`bytes memory`的数组。

- `.push`在数列末尾添加元素，返回值是对这个元素的引用。

- 使用 `new` 创建的 `memory` 类型的数组，内存一旦分配就是固定的，因此，不能够使用`.push`改变数组的大小。

  ```solidity
  pragma solidity >=0.4.16 <0.9.0;
  
  contract TX {
      function f(uint len) public pure {
          uint[] memory a = new uint[](7);
          bytes memory b = new bytes(len);
  
          assert(a.length == 7);
          assert(b.length == len);
  
          a[6] = 8;
      }
  }
  ```

  

#### 定长字节数组

`bytes1`， `bytes2`， `bytes3`， …， `bytes32` 是存放1，2，3，直到 32 个字节的字节序列。它们看成是数组。比较特别的是，它们也可以比较大小，移位，但是不能够进行四则运算。

对于多个字节序列，可以使用`bytes32[k]` 之类的数组存储，但是这样使用很浪费空间，往往还是当成整体来使用，太长时就用下面将介绍的 `bytes` 类型.

- `byte`作为`bytes1`的别名（在0.8.0之前）

- 可以使用`.length`获取字节数，即字节数组长度。

#### 变长字节数组

`bytes`和`string`，当然还有一般的数组类型如`uint[]`

Solidity 没有字符串操作函数但是可以使用第三方的字符串库，不过可以使用keccak256-hash来比较两个字符串`keccak256(abi.encodePacked(s1)) == keccak256(abi.encodePacked(s2))`，或者使用`bytes.concat(bytes(s1), bytes(s2))`来连接两个字符串。

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.3;
contract Hello {
    string public greet = "Hello, ";
    function getstr(string calldata a) public view  returns (string memory){
        return string(bytes.concat(bytes(greet),bytes(a)));
    }
}
```

`bytes` 和 `string`是特殊的数组，一方面元素的内存是紧密连续存放的，不是按照32个字节一单元的方式存放。其中`bytes`可以通过下标访问（`bytes(Name)[index]`或者`Name[index]`)，返回的是底层的`bytes`类型的 UTF-8 码； `string`不能够通过下标访问。我们一般是用固定的`bytes`类型(如`bytes1`,`bytes2` ,...., `bytes32` )，因为 `byte[]` 类型的可变长数组每个元素是占32个字节，一个单元用不完会自动填充0，消耗更多的gas。

#### 数组的赋值和字面常量

数组字面常量是在方括号中（ `[...]` ） 包含一个或多个逗号分隔的表达式。 例如 `[1, a, f(3)]` 。

它是静态（固定大小）的memory类型的数组，长度由元素个数确定，数组的基本类型这样确定：

**通过字面量创建的数组以列表中第一个元素的类型为准，其他的元素会隐式转化**，但是这种转换需要合法。

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.4.16 <0.9.0;

contract C {
    function f() public pure {
        g([uint(1), 2, 3]);
    }
    function g(uint[3] memory) public pure {
        // ...
    }
}
```

上面就是`uint`类型的数组字面常量。`[1,-1]`就是不合法的，因为正整数默认是`uint8`类型，而第二个元素是`-1`，是`int8`类型，数组字面常量的元素的类型就不一致了。`[int8(1),-1]`就是合法的。

更进一步，在2多维数组中，每个子列表的第一个元素都要是同样的类型:

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.4.16 <0.9.0;

contract C {
    function f() public pure returns (uint24[2][4] memory) {
        uint24[2][4] memory x = [[uint24(0x1), 1], [0xffffff, 2], [uint24(0xff), 3], [uint24(0xffff), 4]];
        // The following does not work, because some of the inner arrays are not of the right type.
        // uint[2][4] memory x = [[0x1, 1], [0xffffff, 2], [0xff, 3], [0xffff, 4]];
        return x;
    }
}
```

**通过数组的字面常量创建数组，不支持动态分配内存，必须预设数组大小**。`uint[] memory x = [uint(1), 3, 4];`报错，必须写成`uint[3] memory x = [uint(1), 3, 4];`。这个考虑移除这个特性，但是会造成ABI中数组传参的一些麻烦。

**如果是先创建 memory 的数组，再传参，也不能通过数组的字面常量赋值，必须单独给元素赋值**：

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.4.16 <0.9.0;

contract C {
    function f() public pure {
        uint[] memory x = new uint[](3);
        x[0] = 1;
        x[1] = 3;
        x[2] = 4;
    }
}
```

#### 数组的成员

- `.length`: 返回当前数组的长度。
- `.push()`: 除了`string`类型，其他的动态storage数组和`bytes`都可以使用这个函数在数组的末尾添加一个元素，这个元素默认是0，**返回对这个元素的引用**，`x.push() = b`，修改`b`即可实现对数组元素的修改。
- `.push(x)`: 将 `x` 添加到数组末尾，没有返回值。
- `.pop()`:除了`string`类型，其他的动态数组和`bytes`都可以使用这个函数删除数组的最后一个元素，相当于隐式地`delete`这个元素。（注意 delete 的效果，并不是删除）

可以看出，`push`增加一个元素的gas是固定的，因为储存单元的大小是确定的，但是使用`pop()`等同执行`delete`操作，擦除大量的空间可能会消耗很多gas。

注意：如果需要在外部（external）函数中使用多维数组，这需要启用 ABI coder v2 (在合约最开头加上 `pragma experimental ABIEncoderV2;`，这是为了方便 ABI 编码）。 公有（public）函数中默认支持的使用多维数组。

注意：在Byzantium（在2017-10-16日4370000区块上进行硬分叉升级）之前的EVM版本中，无法访问从函数调用返回动态数组。 如果要调用返回动态数组的函数，请确保 EVM 在拜占庭模式或者之后的模式上运行。



### 结构体

**结构体的辖域**

- 定义在合约之外的结构体类型，可以被所有合约引用。

- 合约内定义的结构体，只能在合约内或者是继承后的合约内可见。

- 结构体的使用和C语言类似，但是注意，结构体不能使用自身。

注意：在 Solidity 0.70 以前`memory`结构体只有`storage`的成员。

**结构体赋值办法：**

`structName(para1, para2, para3, para4)` 或者 `structName(paraName1:para1, paraName2:para2, paraName3:para3)` 



### 映射

映射类型在声明时的形式为 `mapping(_KeyType => _ValueType)`。声明映射类型的变量的形式为 `mapping(_KeyType => _ValueType) _VariableName`.

其中`_KeyType`可以是任何内置的类型，包括`bytes`、`string`以及合约类型和枚举类型，但是不能是自定义的复杂类型，映射、结构体以及数组。`_ValueType`可以是任何类型。但是，映射实际上是哈希表，`key`存储的是`keccak256`的哈希值而不是真实的`key`。因此，底层存储方式上并不是键值对的集合。

**映射只能被声明为 `storage` 类型，不可以作为`public`函数的参数或返回值**。如果结构体或者数组含有映射类型，也需要满足这个规则。

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.4.0 <0.9.0;

contract MappingExample {
    mapping(address => uint) public balances;

    function update(uint newBalance) public {
        balances[msg.sender] = newBalance;
    }
}

contract MappingUser {
    function f() public returns (uint) {
        MappingExample m = new MappingExample();
        m.update(100);
        return m.balances(address(this));
    }
}
```

#### 可迭代的映射

我们使用嵌套映射和结构体，来实现复杂的数据结构，比如链表。以下例子有些难懂，这是通过位置（索引）和关键词构成的链式结构。理解思想即可，需要用到再深入学习。

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.6.0 <0.9.0;

//在合约外定义的变量，是全局变量。

struct IndexValue { uint keyIndex; uint value; } //关键词对应的索引和对应的值
struct KeyFlag { uint key; bool deleted; } //标记关键词是否删除

//这类似于链表。data 用于从当前位置传递到下一个位置，每一个位置都有关键词的索引和值，构成链式结构。
//而KeyFlag 用于记录每个节点（关键词+值）是否删除
//size 标记链表长度
struct itmap {
    mapping(uint => IndexValue) data;
    KeyFlag[] keys;
    uint size;
}

//这是库，里面很多函数可用
library IterableMapping {
    //插入
    function insert(itmap storage self, uint key, uint value) internal returns (bool replaced) {
        uint keyIndex = self.data[key].keyIndex;
        self.data[key].value = value;
        if (keyIndex > 0)
            return true;//已经存在
        else {
            keyIndex = self.keys.length;

            self.keys.push();
            self.data[key].keyIndex = keyIndex + 1;
            self.keys[keyIndex].key = key;
            self.size++;
            return false;
        }
    }
    
    //删除
    function remove(itmap storage self, uint key) internal returns (bool success) {
        uint keyIndex = self.data[key].keyIndex;
        if (keyIndex == 0)
            return false;
        delete self.data[key];
        self.keys[keyIndex - 1].deleted = true;
        self.size --;
    }
    
    //是否包含某个元素
    function contains(itmap storage self, uint key) internal view returns (bool) {
        return self.data[key].keyIndex > 0;
    }
    
    function iterate_start(itmap storage self) internal view returns (uint keyIndex) {
        return iterate_next(self, type(uint).max);
    }

    function iterate_valid(itmap storage self, uint keyIndex) internal view returns (bool) {
        return keyIndex < self.keys.length;
    }

    function iterate_next(itmap storage self, uint keyIndex) internal view returns (uint r_keyIndex) {
        keyIndex++;
        while (keyIndex < self.keys.length && self.keys[keyIndex].deleted)
            keyIndex++;
        return keyIndex;
    }

    function iterate_get(itmap storage self, uint keyIndex) internal view returns (uint key, uint value) {
        key = self.keys[keyIndex].key;
        value = self.data[key].value;
    }
}

// 如何使用
contract User {
    // Just a struct holding our data.
    itmap data;
    // Apply library functions to the data type.
    using IterableMapping for itmap;

    // Insert something
    function insert(uint k, uint v) public returns (uint size) {
        // This calls IterableMapping.insert(data, k, v)
        data.insert(k, v);
        // We can still access members of the struct,
        // but we should take care not to mess with them.
        return data.size;
    }

    // Computes the sum of all stored data.
    function sum() public view returns (uint s) {
        for (
            uint i = data.iterate_start();
            data.iterate_valid(i);
            i = data.iterate_next(i)
        ) {
            (, uint value) = data.iterate_get(i);
            s += value;
        }
    }
}
```

## 类型转换

### 自定义类型

注意 `type UFixed256x18 is uint256 ` 的定义方式

`UFixed256x18.unwrap(a)` 从自定义类型，解封装成内置类型

`UFixed256x18.wrap(a * multiplier)` 从内置类型封装成自定义类型。

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.8;

// Represent a 18 decimal, 256 bit wide fixed point type using a user defined value type.
type UFixed256x18 is uint256;

/// A minimal library to do fixed point operations on UFixed256x18.
library FixedMath {
    uint constant multiplier = 10**18;

    /// Adds two UFixed256x18 numbers. Reverts on overflow, relying on checked
    /// arithmetic on uint256.
    function add(UFixed256x18 a, UFixed256x18 b) internal pure returns (UFixed256x18) {
        return UFixed256x18.wrap(UFixed256x18.unwrap(a) + UFixed256x18.unwrap(b));
    }
    /// Multiplies UFixed256x18 and uint256. Reverts on overflow, relying on checked
    /// arithmetic on uint256.
    function mul(UFixed256x18 a, uint256 b) internal pure returns (UFixed256x18) {
        return UFixed256x18.wrap(UFixed256x18.unwrap(a) * b);
    }
    /// Take the floor of a UFixed256x18 number.
    /// @return the largest integer that does not exceed `a`.
    function floor(UFixed256x18 a) internal pure returns (uint256) {
        return UFixed256x18.unwrap(a) / multiplier;
    }
    /// Turns a uint256 into a UFixed256x18 of the same value.
    /// Reverts if the integer is too large.
    function toUFixed256x18(uint256 a) internal pure returns (UFixed256x18) {
        return UFixed256x18.wrap(a * multiplier);
    }
}
```

### 基本类型转换

 **隐式转换**：隐式转换发生在编译时期，如果不出现信息丢失，其实都可以进行隐式转换，比如`uint8`可以转成`uint16`。隐式转换常发生在不同的操作数一起用操作符操作时发生。

 **显式转换**：如果编译器不允许隐式转换，而你足够自信没问题，那么就去尝试显示转换，但是这很容易造成安全问题。

高版本的Solidity不支持常量的不符合编译器的显式转换，但是允许变量之间进行显式转换。对于`int`转 `uint`就是找补码，负数可以理解为下溢。如果是`uint`或者`int`同类型强制转换，就是从最低位截断(十六进制下，或者从最高位补0。

```solidity
uint32 a = 0x12345678;
uint16 b = uint16(a); // b will be 0x5678 now
```

```solidity
uint16 a = 0x1234;
uint32 b = uint32(a); // b will be 0x00001234 now
assert(a == b);
```

对于`bytes`类型就是从最低位补0或者从最高位开始保留，这样就没有改变原来的下标。

```solidity
bytes2 a = 0x1234;
bytes4 b = bytes4(a); // b will be 0x12340000
assert(a[0] == b[0]);
assert(a[1] == b[1]);
```

只有具有相同字节数的整数和`bytes`类型才允许之间的强制转换，不同长度的需要中间过渡。注意：`bytes32`,表示32个字节，一个字节是8位；`int256`这样指的是二进制位。

```solidity
bytes2 a = 0x1234;
uint32 b = uint16(a); // b will be 0x00001234
uint32 c = uint32(bytes4(a)); // c will be 0x12340000
uint8 d = uint8(uint16(a)); // d will be 0x34
uint8 e = uint8(bytes1(a)); // e will be 0x12
```

`bytes`数组和`calldate`的 `bytes`的切片转换成`bytes32`这样的定长字节类型，截断和填充和定长`bytes`一致。

```solidity
// SPDX-License-Identifier: GPL-3.0
pragma solidity ^0.8.5;

contract C {
    bytes s = "abcdefgh";
    function f(bytes calldata c, bytes memory m) public view returns (bytes16, bytes3) {
        require(c.length == 16, "");
        bytes16 b = bytes16(m);  // if length of m is greater than 16, truncation will happen
        b = bytes16(s);  // padded on the right, so result is "abcdefgh\0\0\0\0\0\0\0\0"
        bytes3 b1 = bytes3(s); // truncated, b1 equals to "abc"
        b = bytes16(c[:8]);  // also padded with zeros
        return (b, b1);
    }
}
```

### 地址类型转换

 `address payable` 可以完成到 `address` 的隐式转换，但是从 `address` 到 `address payable` 必须显式的转换, 通过 `payable(<address>)` 进行转换。 某些函数会严格限制采用哪一种类型。

实际上，合约类型、`uint160`、整数字面常量、`bytes20`都可以与`address`类型互相转换。

- 如果有需要截断的情况，byte类型需要转换成 `uint` 之后才能转换成地址类型。`bytes32`就会被截断，且在`0.4.24`之后需要做显式处理`address(uint(bytes20(b)))`）。
- 合约类型如果已经绑定到已部署的合约，可以显式转换成已部署合约的地址 。
- 字面常量是 "0xabc...."的字符串，可以当作地址类型直接使用。



### 字面常量类型转换

- 0.8.0以后整型的字面产常量的强转必须在满足隐式转化的条件之上，而且整数的隐式转换非常严格，不存在截断。
- 字节型的字面常量只支持同等大小的十六进制数转化，不能由十进制转化。但是如果字面常量是十进制的0或者十六进制的0，那么就允许转换成任何的定长字节类型。

```solidity
bytes2 a = 54321; // not allowed
bytes2 b = 0x12; // not allowed
bytes2 c = 0x123; // not allowed
bytes2 d = 0x1234; // fine
bytes2 e = 0x0012; // fine
bytes4 f = 0; // fine
bytes4 g = 0x0; // fine
```

- 字符串字面常量转换成定长字节类型也需要大小相同。

```solidity
bytes2 a = hex"1234"; // fine
bytes2 b = "xy"; // fine
bytes2 c = hex"12"; // not allowed
bytes2 d = hex"123"; // not allowed
bytes2 e = "x"; // not allowed
bytes2 f = "xyz"; // not allowed
```

- 只有大小正确（40位十六进制，160个字节）的满足检验和的十六进制常量才能转换成地址类型。

### 函数可见性类型转化

函数的可见性类型可以发生隐式的转化，**规则是**：只能变得比以前更严格，不能改变原来的限制条件，只能增加更多的限制条件。

**有且仅有以下三种转化：**

- `pure` 函数可以转换为 `view` 和 `non-payable` 函数
- `view` 函数可以转换为 `non-payable` 函数
- `payable` 函数可以转换为 `non-payable` 函数

如果在{internal,external}的位置是`public`，那么函数既可以当作内部函数，也可以当作外部函数使用，如果只想当内部函数使用，就用`f`（函数名）调用，如果想当作外部函数调用，使用`this.f`（地址+函数名，合约对象.函数名）

