本文旨在快速介绍Go标准库中读取文件的许多选项。

在Go中（就此而言，大多数底层语言和某些动态语言（如Node））返回字节流。 不将所有内容自动转换为字符串的好处是，其中之一是避免昂贵的字符串分配，这会增加GC压力。

为了使本文更加简单，我将使用`string(arrayOfBytes)`将`bytes`数组转换为字符串。 但是，在发布生产代码时，不应将其作为一般建议。

## 1.读取整个文件到内存中

首先，标准库提供了多种功能和实用程序来读取文件数据。我们将从os软件包中提供的基本情况开始。这意味着两个先决条件：
- 该文件必须容纳在内存中
- 我们需要预先知道文件的大小，以便实例化一个足以容纳它的缓冲区。


有了`os.File`对象的句柄，我们可以查询大小并实例化一个字节列表。

```
package main


import (
	"os"
	"fmt"
)
func main() {
	file, err := os.Open("filetoread.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	bytesread, err := file.Read(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("bytes read: ", bytesread)
	fmt.Println("bytestream to string: ", string(buffer))
}
```

## 2.以块的形式读取文件
虽然大多数情况下可以一次读取文件，但有时我们还是想使用一种更加节省内存的方法。例如，以某种大小的块读取文件，处理它们，并重复直到结束。在下面的示例中，使用的缓冲区大小为100字节。
```
package main


import (
	"io"
	"os"
	"fmt"
)

const BufferSize = 100

func main() {
	
	file, err := os.Open("filetoread.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	buffer := make([]byte, BufferSize)

	for {
		bytesread, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		fmt.Println("bytes read: ", bytesread)
		fmt.Println("bytestream to string: ", string(buffer[:bytesread]))
	}
}

```
与完全读取文件相比，主要区别在于：
- 读取直到获得EOF标记，因此我们为`err == io.EOF`添加了特定检查
- 我们定义了缓冲区的大小，因此我们可以控制所需的“块”大小。 如果操作系统正确地将正在读取的文件缓存起来，则可以在正确使用时提高性能。
- 如果文件大小不是缓冲区大小的整数倍，则最后一次迭代将仅将剩余字节数添加到缓冲区中，因此调用`buffer [：bytesread]`。 在正常情况下，`bytesread`将与缓冲区大小相同。


对于循环的每次迭代，都会更新内部文件指针。 下次读取时，将返回从文件指针偏移开始直到缓冲区大小的数据。 该指针不是语言的构造，而是操作系统之一。 在Linux上，此指针是要创建的文件描述符的属性。 所有的read / Read调用（分别在Ruby / Go中）在内部都转换为系统调用并发送到内核，并且内核管理此指针。

## 3.并发读取文件块

如果我们想加快对上述块的处理，该怎么办？一种方法是使用多个go例程！与串行读取块相比，我们需要做的另一项工作是我们需要知道每个例程的偏移量。请注意，当目标缓冲区的大小大于剩余的字节数时，ReadAt的行为与Read的行为略有不同。

另请注意，我并没有限制`goroutine`的数量，它仅由缓冲区大小来定义。实际上，此数字可能会有上限。

```
package main

import (
	"fmt"
	"os"
	"sync"
)

const BufferSize = 100

type chunk struct {
	bufsize int
	offset  int64
}

func main() {
	
	file, err := os.Open("filetoread.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	filesize := int(fileinfo.Size())
	// Number of go routines we need to spawn.
	concurrency := filesize / BufferSize
	// buffer sizes that each of the go routine below should use. ReadAt
	// returns an error if the buffer size is larger than the bytes returned
	// from the file.
	chunksizes := make([]chunk, concurrency)

	// All buffer sizes are the same in the normal case. Offsets depend on the
	// index. Second go routine should start at 100, for example, given our
	// buffer size of 100.
	for i := 0; i < concurrency; i++ {
		chunksizes[i].bufsize = BufferSize
		chunksizes[i].offset = int64(BufferSize * i)
	}

	// check for any left over bytes. Add the residual number of bytes as the
	// the last chunk size.
	if remainder := filesize % BufferSize; remainder != 0 {
		c := chunk{bufsize: remainder, offset: int64(concurrency * BufferSize)}
		concurrency++
		chunksizes = append(chunksizes, c)
	}

	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 0; i < concurrency; i++ {
		go func(chunksizes []chunk, i int) {
			defer wg.Done()

			chunk := chunksizes[i]
			buffer := make([]byte, chunk.bufsize)
			bytesread, err := file.ReadAt(buffer, chunk.offset)

			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("bytes read, string(bytestream): ", bytesread)
			fmt.Println("bytestream to string: ", string(buffer))
		}(chunksizes, i)
	}

	wg.Wait()
}
```
与以前的任何方法相比，这种方法要多得多：
- 我正在尝试创建特定数量的Go例程，具体取决于文件大小和缓冲区大小（在本例中为100）。
- 我们需要一种方法来确保我们正在“等待”所有执行例程。 在此示例中，我使用的是wait group。
- 在每个例程结束的时候，从内部发出信号，而不是`break for`循环。因为我们延时调用了`wg.Done()`,所以在每个例程返回的时候才调用它。

**注意：始终检查返回的字节数，并重新分配输出缓冲区。**



使用`Read()`读取文件可以走很长一段路，但是有时您需要更多的便利。`Ruby`中经常使用的是`IO`函数，例如`each_line`,`each_char`, `each_codepoint` 等等.通过使用`Scanner`类型以及`bufio`软件包中的关联函数，我们可以实现类似的目的。 

`bufio.Scanner`类型实现带有“ split”功能的函数，并基于该功能前进指针。例如，对于每个迭代，内置的`bufio.ScanLines`拆分函数都会使指针前进，直到下一个换行符为止.
在每个步骤中，该类型还公开用于获取开始位置和结束位置之间的字节数组/字符串的方法。

```
package main

import (
	"fmt"
	"os"
	"bufio"
)

const BufferSize = 100

type chunk struct {
	bufsize int
	offset  int64
}

func main() {
	file, err := os.Open("filetoread.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// Returns a boolean based on whether there's a next instance of `\n`
	// character in the IO stream. This step also advances the internal pointer
	// to the next position (after '\n') if it did find that token.
	for {
		read := scanner.Scan()
		if !read {
			break
			
		}
		fmt.Println("read byte array: ", scanner.Bytes())
		fmt.Println("read string: ", scanner.Text())
	}
	
}
```

因此，要以这种方式逐行读取整个文件，可以使用如下所示的内容：
```
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("filetoread.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	// This is our buffer now
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	fmt.Println("read lines:")
	for _, line := range lines {
		fmt.Println(line)
	}
}
```

## 4.逐字扫描
bufio软件包包含基本的预定义拆分功能：
- ScanLines (默认)
- ScanWords
- ScanRunes(对于遍历UTF-8代码点（而不是字节）非常有用)
- ScanBytes

因此，要读取文件并在文件中创建单词列表，可以使用如下所示的内容：
```
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("filetoread.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var words []string

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	fmt.Println("word list:")
	for _, word := range words {
		fmt.Println(word)
	}
}
```
`ScanBytes`拆分函数将提供与早期`Read()`示例相同的输出。 两者之间的主要区别是在扫描程序中，每次需要附加到字节/字符串数组时，动态分配问题。 可以通过诸如将缓冲区预初始化为特定长度的技术来避免这种情况，并且只有在达到前一个限制时才增加大小。 使用与上述相同的示例：
```
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	file, err := os.Open("filetoread.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	// initial size of our wordlist
	bufferSize := 50
	words := make([]string, bufferSize)
	pos := 0

	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			// This error is a non-EOF error. End the iteration if we encounter
			// an error
			fmt.Println(err)
			break
		}

		words[pos] = scanner.Text()
		pos++

		if pos >= len(words) {
			// expand the buffer by 100 again
			newbuf := make([]string, bufferSize)
			words = append(words, newbuf...)
		}
	}

	fmt.Println("word list:")
	// we are iterating only until the value of "pos" because our buffer size
	// might be more than the number of words because we increase the length by
	// a constant value. Or the scanner loop might've terminated due to an
	// error prematurely. In this case the "pos" contains the index of the last
	// successful update.
	for _, word := range words[:pos] {
		fmt.Println(word)
	}
}
```
因此，我们最终要进行的切片“增长”操作要少得多，但最终可能要根据缓冲区大小和文件中的单词数在结尾处留出一些空插槽，这是一个折衷方案。

## 5.将长字符串拆分为单词
`bufio.NewScanner`使用满足`io.Reader`接口的类型作为参数，这意味着它将与定义了`Read`方法的任何类型一起使用。
标准库中返回`reader`类型的`string`实用程序方法之一是`strings.NewReader`函数。当从字符串中读取单词时，我们可以将两者结合起来：
```
package main

import (
	"bufio"
	"fmt"
	"strings"
)

func main() {
	longstring := "This is a very long string. Not."
	var words []string
	scanner := bufio.NewScanner(strings.NewReader(longstring))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	fmt.Println("word list:")
	for _, word := range words {
		fmt.Println(word)
	}
}
```
## 6.扫描以逗号分隔的字符串
手动解析CSV文件/字符串通过基本的`file.Read()`或者`Scanner`类型是复杂的。因为根据拆分功能`bufio.ScanWords`，“单词”被定义为一串由unicode空间界定的符文。读取各个符文并跟踪缓冲区的大小和位置（例如在词法分析中所做的工作）是太多的工作和操作。

但这可以避免。 我们可以定义一个新的拆分函数，该函数读取字符直到读者遇到逗号，然后在调用`Text（）`或`Bytes（）`时返回该块。`bufio.SplitFunc`函数的函数签名如下所示：
```
type SplitFunc func(data []byte, atEOF bool) (advance int, token []byte, err error)
```

为简单起见，我展示了一个读取字符串而不是文件的示例。 使用上述签名的CSV字符串的简单阅读器可以是：

```
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

func main() {
	csvstring := "name, age, occupation"

	// An anonymous function declaration to avoid repeating main()
	ScanCSV := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		commaidx := bytes.IndexByte(data, ',')
		if commaidx > 0 {
			// we need to return the next position
			buffer := data[:commaidx]
			return commaidx + 1, bytes.TrimSpace(buffer), nil
		}

		// if we are at the end of the string, just return the entire buffer
		if atEOF {
			// but only do that when there is some data. If not, this might mean
			// that we've reached the end of our input CSV string
			if len(data) > 0 {
				return len(data), bytes.TrimSpace(data), nil
			}
		}

		// when 0, nil, nil is returned, this is a signal to the interface to read
		// more data in from the input reader. In this case, this input is our
		// string reader and this pretty much will never occur.
		return 0, nil, nil
	}

	scanner := bufio.NewScanner(strings.NewReader(csvstring))
	scanner.Split(ScanCSV)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}
```

## 7.`ioutil`
我们已经看到了多种读取文件的方式.但是，如果您只想将文件读入缓冲区怎么办？


`ioutil`是标准库中的软件包，其中包含一些使它成为单行的功能。

### 读取整个文件

```
package main

import (
	"io/ioutil"
	"log"
	"fmt"
)

func main() {
	bytes, err := ioutil.ReadFile("filetoread.txt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bytes read: ", len(bytes))
	fmt.Println("String read: ", string(bytes))
}
```

这更接近我们在高级脚本语言中看到的内容。

### 读取文件的整个目录

不用说，如果您有大文件，请不要运行此脚本
```
package main

import (
	"io/ioutil"
	"log"
	"fmt"
)

func main() {
	filelist, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}
	for _, fileinfo := range filelist {
		if fileinfo.Mode().IsRegular() {
			bytes, err := ioutil.ReadFile(fileinfo.Name())
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Bytes read: ", len(bytes))
			fmt.Println("String read: ", string(bytes))
		}
	}
}
```



## 参考文献

- [go语言读取文件概述](https://kgrz.io/reading-files-in-go-an-overview.html#reading-byte-wise)
