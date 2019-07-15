# 1.简介
Protobuf是Protocol Buffers的简称，它是Google公司开发的一种数据描述语言，用于序列化结构化数据 - 类似XML，但更小，更快，更简单。 您可以定义数据的结构化结构，然后使用特殊生成的源代码轻松地将结构化数据写入和读取各种数据流，并使用各种语言。可以作为RPC的基础工具
#2.定义一个Message类型
```
/* SearchRequest represents a search query, with pagination options to
 * indicate which results to include in the response. */
syntax = "proto3";
message SearchRequest {
    string query = 1;
    int32 page_number = 2; // Which page number do we want?
    int32 result_per_page = 3;// Number of results to return per page.
}
```
- 1.第一行指定语法版本，本示例采用的是proto3的语法，第三版的protobuf对语言进行了提炼简化，所有成员均采用类似Go语言中的零值初始化，因此，消息成员不再需要支持required特性。
- 2.定义三个字段，每个字段有一个名字和类型
- 3.每个字段都定义的有唯一数字标签，该数值定义好以后，不允许改变。这些字段用来标示字段在信息中的二进制形式，取值范围为1-15占用一个字节编码，16-2047，占用两个字节。
- 最小的标签值可以指定为 1, 最大的为 229 - 1, or 536,870,911. 不能使用 19000 到 19999 之间的值。
- 如果您通过完全删除一个字段或将其注释掉来更新消息类型，将来的用户可以在对该类型进行自己的更新时重用这个标记号。这可能会导致严重的问题，如果它们稍后加载相同的旧版本.proto，包括数据损坏、隐私bug等等。确保这不会发生的一种方法是指定字段标记(以及/或名称，也可能导致JSON序列化的问题)被保留。如果将来的用户尝试使用这些字段标识符，协议缓冲编译器将会抱怨。
```
message Foo {
   reserved 2, 15, 9 to 11;
   reserverd "foo", "bar";
}
```
注意，不能在同一保留语句中混合字段名称和标记号。

编译后产出的文件

对于go，产生一个.pb.go文件

#3.标量类型
```
.proto type ====> GO type
double            float64
float             float
int32/sint32/sfixed32      int32
int64/sint64/sfixed64      int64
uint32/fixed32    uint32
uint64/fixed64    uint64
bool              bool
string            string
bytes             []byte
```
# 4.默认值
- string 默认值为空的字符串
- bytes 默认值为空的bytes
- bools，默认值为false
- numeric，默认值为0
- enums,默认值是第一个定义的enum值，必需为0
- message字段未设置时，确切的值依赖语言环境。
- repeated 字段的默认值为空

#5.Enumerations
```
message SearchRequest {
  string query = 1;
  int32 page_number = 2;
  int32 result_per_page = 3;
  enum Corpus {
    UNIVERSAL = 0;
    WEB = 1;
    IMAGES = 2;
    LOCAL = 3;
    NEWS = 4;
    PRODUCTS = 5;
    VIDEO = 6;
  }
  Corpus corpus = 4;
}
```
Corpus字段是一个枚举类型，第一个常量影射为0.每个枚举类型必需包含一个常量映射为0，作为第一个元素。
可以定义一个别名通过分配同样的值给不同的枚举常量。这么做首先需要设置allow_alias选项为true.否则会报错。

```
enum EnumAllowingAlias {
  option allow_alias = true;
  UNKNOWN = 0;
  STARTED = 1;
  RUNNING = 1;
}
```
#6.使用其他message类型
可以使用其他message 类型作为字段类型。例如，Result 作为SearchResponse的字段。

```
message SearchResponse {
    repeated Result results = 1;
}
message Result {
    string url = 1;
    string title = 2;
    repeated string snippets = 3;
}
```
上面两个类型定义在同一个.proto文件中，如果Result定义在另一个文件中呢？
在需要导入其他文件定义的文件中头部

```
import "myproject/other_protos.proto";
```

#7.嵌套类型
```
message SearchResponse {
    message Result {
        string url = 1;
        string title = 2;
        repeated string snippets = 3;
    }
    repeated Result results = 1;
}
```
如果想重用在父类型之外的其他类型中，使用Parent.Type;

```
message SomeOtherMessage {
    SearchResponse.Result result = 1;
}
```
嵌套深度无限制。
# 8.更新一个Messsage 
如果一个类型不满足需求了，比如需要添加额外的字段。但是依然使用原来的形式，如下：

- 不要修改已存在字段的数字标签
- 如果您添加了新的字段，那么使用您的“旧”消息格式序列化的任何消息仍然可以通过您的新生成的代码进行解析。您应该记住这些元素的默认值，以便新代码能够正确地与旧代码生成的消息交互。类似地，新代码创建的消息可以通过旧代码解析:旧二进制代码在解析时忽略了新字段。有关详细信息，请参见未知字段部分。
- 字段可以移除，只要它的数字标签不在新的消息类型中使用。如果想要重命名字段。可添加前缀'OBSOLETE_',或者保留数字标签。
- int32、uint32、int64、uint64和bool都是兼容的，这意味着您可以在不破坏转发或向后兼容的情况下，将字段从这些类型转换为另一种。如果一个数字从与对应类型不匹配的线中解析，您将得到与在c++中(例如，如果一个64位数字被读取为int32，它将被截断为32位)的相同效果。
- sint32和sint64是兼容的，但是不兼容其他整数类型。
- string和bytes是兼容的，只要bytes是有效的UTF-8。
- 如果bytes包含消息的编码版本，嵌入式消息与bytes兼容。
- fixed32与sfixed32兼容，fixed64与sfixed64兼容。
- enum与int32、uint32、int64和uint64兼容，在有线格式方面(注意，如果不合适，值将被截断)。但是，请注意，当消息被反序列化时，客户端代码可能会以不同的方式对待它们:例如，在消息中会保留未被识别的proto3 enum类型，但是当消息被反序列化时，如何表示这个消息是依赖于语言的。Int字段总是保持它们的值。
# 9.oneof
oneof字段就像常规的字段，除了oneof中的字段共享内存，任何时候，只有一个字段的值被设置。设置任何一个字段，会自动清除其他成员。你可查询那个字段设置了，使用case()或者WhichOneof() 方法，根据你选择的语言。
```
message SampleMessage {
  oneof test_oneof {
    string name = 4;
    SubMessage sub_message = 9;
  }
}
```
- You then add your oneof fields to the oneof definition. You can add fields of any type, but cannot use repeated fields.
- Setting a oneof field will automatically clear all other members of the oneof. So if you set several oneof fields, only the last field you set will still have a value.
- 
