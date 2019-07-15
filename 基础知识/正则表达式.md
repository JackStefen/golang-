```
package main


import (
    "fmt"
    "regexp"
)


func main() {
    text := `<p>更多分析师观点详见&nbsp;
<a href="https://m.xxx.cn/article/160267" target="_blank" rel="noopener">
一周策略前瞻：周期之火破灭了？</a></p>
<p><strong>来看主题：</strong></p>
<p><strong>
1、<a class="" href="https://api.xxx.cn/web/subjects/186">PPP</a>&nbsp;：
</strong>
<a href="https://m.xxx.cn/article/160320" target="_blank" rel="noopener">
国办发文力促民资参与PPP，经济回落下行业有望再成>稳增长抓手&nbsp;</a></p>
<p>参见上文逻辑，此处不多说了。地产板块也是类似。</p>`

    //var SubjectRegexp = regexp.MustCompile(`<a href="https://[[:ascii:]]*">(?P<ct>.*)</a>`)
    var ArticleRegexp = regexp.MustCompile(`<a.*href="https://(m|api).xxx.cn(.*)/(article|subjects)/[\d]+"(.*)>(.+)</a>`)
    fmt.Println(ArticleRegexp.FindAllString(text, -1))
    //fmt.Println(SubjectRegexp.FindAllString(text, -1))
    //fmt.Println(SubjectRegexp.ReplaceAllString(text, `${ct}`))
    text2 := `I'm singing while you're dancing.`
    RegExpIng := regexp.MustCompile(`((\')\w{1,2})`)
    fmt.Println(RegExpIng.FindAllString(text2, -1))
}
```
Output:
```
[<a href="https://m.xxx.cn/article/160267" target="_blank" rel="noopener">一周策略前瞻：周期之火破灭了？</a> <a class="" href="https://api.xxx.cn/web/subjects/186">PPP</a>&nbsp;：</strong><a href="https://m.xxx.cn/article/160320" target="_blank" rel="noopener">国办发文力促民资参与PPP，经济回落下行业有望再成稳增长抓手&nbsp;</a>]
['m 're]
```
# 1.`MustCompile(...)` VS `Compile(...)`
```
func Compile(expr string) (*Regexp, error) {
    return compile(expr, syntax.Perl, false)
}
```
MustComile实际上调用的是Compile。加了错误检测。
```
func MustCompile(str string) *Regexp {
    regexp, error := Compile(str)
    if error != nil {
        panic(`regexp: Compile(` + quote(str) + `): ` + error.Error())
    }
    return regexp
}
```

# 2. MatchString检测是否匹配正则，参数为被检测的字符串，返回布尔值
```
// MatchString reports whether the string s
// contains any match of the regular expression re.
func (re *Regexp) MatchString(s string) bool {
	return re.doMatch(nil, nil, s)
}
```

# 3. `FindAllString(...)`
有两个参数，第一个参数为要处理的字符串，第二个参数获取匹配的结果数量，如果为负数，则取出所有满足条件的匹配结果
```
// FindAllString is the 'All' version of FindString; it returns a slice of all
// successive matches of the expression, as defined by the 'All' description
// in the package comment.
// A return value of nil indicates no match.
func (re *Regexp) FindAllString(s string, n int) []string {
	if n < 0 {
		n = len(s) + 1
	}
	result := make([]string, 0, startSize)
	re.allMatches(s, nil, n, func(match []int) {
		result = append(result, s[match[0]:match[1]])
	})
	if len(result) == 0 {
		return nil
	}
	return result
}
```
其中核心是调用了 allMatches的私有方法获取的结果。该方法的第一个参数为要处理的文本字符串，第二个参数为字节数字切片，在FindAllString中使用的空指针。第三个参数为FindAllString的第二个参数n,第四个参数为一个函数，它负责把所有的收集。
```
// Find matches in slice b if b is non-nil, otherwise find matches in string s.
func (re *Regexp) allMatches(s string, b []byte, n int, deliver func([]int)) {
	var end int
	if b == nil {
		end = len(s)
	} else {
		end = len(b)
	}

	for pos, i, prevMatchEnd := 0, 0, -1; i < n && pos <= end; {
		matches := re.doExecute(nil, b, s, pos, re.prog.NumCap, nil)
		if len(matches) == 0 {
			break
		}

		accept := true
		if matches[1] == pos {
			// We've found an empty match.
			if matches[0] == prevMatchEnd {
				// We don't allow an empty match right
				// after a previous match, so ignore it.
				accept = false
			}
			var width int
			// TODO: use step()
			if b == nil {
				_, width = utf8.DecodeRuneInString(s[pos:end])
			} else {
				_, width = utf8.DecodeRune(b[pos:end])
			}
			if width > 0 {
				pos += width
			} else {
				pos = end + 1
			}
		} else {
			pos = matches[1]
		}
		prevMatchEnd = matches[1]

		if accept {
			deliver(re.pad(matches))
			i++
		}
	}
}
```

```
re := regexp.MustCompile("a.")
fmt.Println(re.FindAllString("paranormal", -1))
fmt.Println(re.FindAllString("paranormal", 2))
fmt.Println(re.FindAllString("graal", -1))
fmt.Println(re.FindAllString("none", -1))
```
Output:
```
[ar an al]
[ar an]
[aa]
[]
```
# 4.`ReplaceAllString(...)`
替换所有匹配到的结果为指定的字符串。第二个参数给出了要替换的值
```
func (re *Regexp) ReplaceAllString(src, repl string) string {
    n := 2
    if strings.Contains(repl, "$") {
        n = 2 * (re.numSubexp + 1)
    }
    b := re.replaceAll(nil, src, n, func(dst []byte, match []int) []byte {
        return re.expand(dst, repl, nil, src, match)
    })
    return string(b)
}
```
replaceAll
```
func (re *Regexp) replaceAll(bsrc []byte, src string, nmatch int, repl func(dst []byte, m []int) []byte) []byte {
	lastMatchEnd := 0 // end position of the most recent match
	searchPos := 0    // position where we next look for a match
	var buf []byte
	var endPos int
	if bsrc != nil {
		endPos = len(bsrc)
	} else {
		endPos = len(src)
	}
	if nmatch > re.prog.NumCap {
		nmatch = re.prog.NumCap
	}

	var dstCap [2]int
	for searchPos <= endPos {
		a := re.doExecute(nil, bsrc, src, searchPos, nmatch, dstCap[:0])
		if len(a) == 0 {
			break // no more matches
		}

		// Copy the unmatched characters before this match.
		if bsrc != nil {
			buf = append(buf, bsrc[lastMatchEnd:a[0]]...)
		} else {
			buf = append(buf, src[lastMatchEnd:a[0]]...)
		}

		// Now insert a copy of the replacement string, but not for a
		// match of the empty string immediately after another match.
		// (Otherwise, we get double replacement for patterns that
		// match both empty and nonempty strings.)
		if a[1] > lastMatchEnd || a[0] == 0 {
			buf = repl(buf, a)
		}
		lastMatchEnd = a[1]

		// Advance past this match; always advance at least one character.
		var width int
		if bsrc != nil {
			_, width = utf8.DecodeRune(bsrc[searchPos:])
		} else {
			_, width = utf8.DecodeRuneInString(src[searchPos:])
		}
		if searchPos+width > a[1] {
			searchPos += width
		} else if searchPos+1 > a[1] {
			// This clause is only needed at the end of the input
			// string. In that case, DecodeRuneInString returns width=0.
			searchPos++
		} else {
			searchPos = a[1]
		}
	}

	// Copy the unmatched characters after the last match.
	if bsrc != nil {
		buf = append(buf, bsrc[lastMatchEnd:]...)
	} else {
		buf = append(buf, src[lastMatchEnd:]...)
	}

	return buf
}
```
# 5. ReplaceAllStringFunc
```
func ConvertTabToEmptyString(text string) string {
	return TabRegExp.ReplaceAllStringFunc(text, func(matchedStr string) string {
		return strings.Replace(matchedStr, "   ", " ", -1)
	})
}
```

# 6.FindAllStringSubmatch
找出有匹配到的字符串子组列表，第二个参数小于0，表示全部匹配
```
// FindAllStringSubmatch is the 'All' version of FindStringSubmatch; it
// returns a slice of all successive matches of the expression, as defined by
// the 'All' description in the package comment.
// A return value of nil indicates no match.
func (re *Regexp) FindAllStringSubmatch(s string, n int) [][]string {
	if n < 0 {
		n = len(s) + 1
	}
	var result [][]string
	re.allMatches(s, nil, n, func(match []int) {
		if result == nil {
			result = make([][]string, 0, startSize)
		}
		slice := make([]string, len(match)/2)
		for j := range slice {
			if match[2*j] >= 0 {
				slice[j] = s[match[2*j]:match[2*j+1]]
			}
		}
		result = append(result, slice)
	})
	return result
}
```
示例：
```
func AlliRemLinkUrls(articleArr []*Article) {
	for _, article := range articleArr {
		if strArrArr := LinkUrlRegExp.FindAllStringSubmatch(article.Content, -1); strArrArr != nil {
			for _, strArr := range strArrArr {
				article.Content = strings.Replace(article.Content, strArr[0], strArr[3], 1)
			}
		}
	}
}
```