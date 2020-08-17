### 环境
- 服务器环境：ubuntu 
- elasticsearch 版本为7.6.0
- 使用`go-elasticsearch`来操作elasticsearch

## 1.坑1--在安装完成中文分词后，未重启elasticsearch

针对于elasticsearch的安装就不多说了，大家可以在网上搜出很多种，当然，也可以直接使用docker来运行elasticsearch。当然了，这种方式的话，在需要定制化elasticsearch时，需要把本地的配置文件通过docker进行映射。

回过头来，先讲一下关于中文分词的插件安装，我使用的是`ik`.如果大家有更好的推荐，欢迎在评论区留言。关于`ik`的安装，我简单介绍一下。
首先，找到安装elastic的目录，进入可执行文件目录，可以使用插件安装命令，安装要安装的插件。
```
./bin/elasticsearch-plugin install https://github.com/medcl/elasticsearch-analysis-ik/releases/download/v7.6.0/elasticsearch-analysis-ik-7.6.0.zip
```

其中`v7.6.0`是插件的版本，大家可以根据自己的elastic版本，进行安装自己版本相关的插件。

其次，最重要的是，**在第一步，安装插件命令执行完成后，需要我们手动重启一下elasticsearch。因为只有重启后，在建立索引时，才能使用我们的插件分词器。**

在ubuntu环境中，直接`systemctl restart elasticsearch`

说明一下，为什么要特别注意这一点呢，原因如下：
- elasticsearch的索引，在使用索引创建语言创建失败后，在首次执行插入数据的时候，elasticsearch会根据插入数据的特征，自动创建索引。哈哈，是不是感觉很受伤，我都将数据成功写入到elastic了，分词也安装了，为啥查询效果不佳，特别留意哦


举个例子来看看，我们的分词器是否正常可用了呢。

```
➜  go-elasticsearch git:(master) curl -XPUT 'http://127.0.0.1:9200/bbs?pretty'
{
  "acknowledged" : true,
  "shards_acknowledged" : true,
  "index" : "bbs"
}

```

如果在安装了分词器，但是为重启的情况下，直接创建带有中文分词的映射时，就会出现创建失败的情况：

```
➜  go-elasticsearch git:(master) curl -XPOST -H 'Content-Type:application/json' '127.0.0.1:9200/bbs/_mapping?pretty' -d '{"properties":{"contetn":{"type":"text","analyzer":"ik_max_word","search_analyzer":"ik_smart"}}}'             
{
  "error" : {
    "root_cause" : [
      {
        "type" : "mapper_parsing_exception",
        "reason" : "analyzer [ik_smart] not found for field [contetn]"
      }
    ],
    "type" : "mapper_parsing_exception",
    "reason" : "analyzer [ik_smart] not found for field [contetn]"
  },
  "status" : 400
}
```

重启elastic后，再次创建映射，就可以成功创建了

```
➜  curl  '127.0.0.1:9200/bbs/_mappings?pretty' -H 'Content-Type:application/json' 
{
  "bbs" : {
    "mappings" : {
      "properties" : {
        "content" : {
          "type" : "text",
          "analyzer" : "ik_max_word",
          "search_analyzer" : "ik_smart"
        }
      }
    }
  }
}

```

## 2.坑2--go-elasticsearch的api在操作数组类型的值时

关于这个问题，首先一点，需要说明的是，这个其实跟自身的golang语言的熟练程度是有很大关系的。这里特别提出来，倒不是说这人家`go-elasticsearch`的工具库的问题，希望不要跟大家带来困扰。


具体问题描述，比如，我现在在`bbs`这个映射上加了一个标签功能，这样可以在对内容进行搜索的时候，可以指定标签进行搜索。

我们通过例子来看看，到底是个什么情况：

```
➜  curl -XPOST -H 'Content-Type:application/json' '127.0.0.1:9200/bbs/create?pretty' -d '{"tags":["科技","财经"]}'
{
  "_index" : "bbs",
  "_type" : "create",
  "_id" : "AXPwEfYmtmwwbP36AZPn",
  "_version" : 1,
  "result" : "created",
  "_shards" : {
    "total" : 2,
    "successful" : 1,
    "failed" : 0
  },
  "created" : true
}
```
我们通过插入一条数据到elastic中，其中，tags字段带有一个列表值，其元素代表，这条bbs记录的打上了科技和财经的标签。

这种数据的`_mappings`形式，我们来看一下：

```
➜  curl '127.0.0.1:9200/bbs/_mappings?pretty'                                                         
{
  "bbs" : {
    "mappings" : {
      "create" : {
        "properties" : {
          "tags" : {
            "type" : "text",
            "fields" : {
              "keyword" : {
                "type" : "keyword",
                "ignore_above" : 256
              }
            }
          }
        }
      }
    }
  }
}
```
它实际上是用到了子字段的语法形式。所以，我们在搜索特定标签的时候，需要特别注意这一点，下面的时候，我还继续会特别提这一点。

继续看，我们如果使用curl在命令行上进行查询操作，你发现，其实并没有任何问题。比如：

```
➜  curl '127.0.0.1:9200/bbs/_search?pretty' -H 'Content-Type:application/json' -d '{"query":{"bool":{"must":[{"match":{"content":"拼多多"}}],"filter":[{"terms":{"tags.keyword":["科技"]}}]}},"size":2}' 
{
  "took" : 73,
  "timed_out" : false,
  "_shards" : {
    "total" : 5,
    "successful" : 5,
    "skipped" : 0,
    "failed" : 0
  },
  "hits" : {
    "total" : 1,
    "max_score" : 1.0707065,
    "hits" : [
      {
        "_index" : "bbs",
        "_type" : "create",
        "_id" : "AXPwGUpltmwwbP36AZPp",
        "_score" : 1.0707065,
        "_source" : {
          "content" : "拼多多支持消费者维权",
          "tags" : [
            "科技",
            "财经"
          ]
        }
      }
    ]
  }
}
```

在使用`go-elasticsearch`进行查询的时候，我当初使用的方式是，使用`buff.Buffer`类型的变量，进行字符串写入，在写入列表形式的数据时，我使用的是`fmt.Sprintf("%v", list)`.我自认为应该是没问题的。而实际上现实给我了一记响亮的耳光。我们来看一下，我当初的代码细节

```
    var buf strings.Builder
	query := fmt.Sprintf(`{
		"query": {
			"bool":{
				"must":[{
					"match":{ "content":   "%s"       }
				}],
				"filter":[
					{
						"terms":{
							"tags.keyword":"%v"
						}
					}
	`, "拼多多", []string{"科技"，"财经"})

	buf.WriteString(query)

	buf.WriteString(`]}},`)
	buf.WriteString(fmt.Sprintf(`"size":%d}`, 10))
	var r map[string]interface{}

	req := esapi.SearchRequest{
		Index:  []string{"bbs"},
		Body:   strings.NewReader(buf.String()),
		Pretty: true,
	}
	// Perform the search request.
	res, err := req.Do(context.Background(), esC)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error getting response: %s", err))
	}
	defer res.Body.Close()

```

但是，执行查询的过程中，一直报错，具体错误细节为：

```
[400 Bad Request] parsing_exception: [terms] query does not support [tags.keyword]
```

这真让人头大，难道是语法不支持？

一顿操作，测试后发现，根本不是语法不支持。只是提供的字符串形式的列表，在通过json处理之后，提供给elastic进行查询的时候，语法不对，elastic无法进行识别。

```
 curl '127.0.0.1:9200/bbs/_search?pretty' -H 'Content-Type:application/json' -d '{"query":{"bool":{"must":[{"match":{"content":"拼多多"}}],"filter":[{"terms":{"tags.keyword":"["科技"]"}}]}},"size":2}'
{
  "error" : {
    "root_cause" : [
      {
        "type" : "parsing_exception",
        "reason" : "[terms] query does not support [tags.keyword]",
        "line" : 1,
        "col" : 97
      }
    ],
    "type" : "parsing_exception",
    "reason" : "[terms] query does not support [tags.keyword]",
    "line" : 1,
    "col" : 97
  },
  "status" : 400
}
```
看到么，就这个问题，搞了我一夜，而这仅仅是golang语法的细节掌握不到位。而通过调整后，一且都豁然开朗

```
    var buf strings.Builder
	query := fmt.Sprintf(`{
		"query": {
			"bool":{
				"must":[{
					"match":{ "content":   "%s"       }
				}],
				"filter":[
	`, "拼多多")

	buf.WriteString(query)
    topicTops := []string{"科技"，"财经"}
	if len(topicTops) != 0 {
		topicbuf := bytes.NewBuffer([]byte{})
		topicse := map[string]interface{}{
			"terms": map[string][]string{
				"tags.keyword": topicTops,
			},
		}
		if err := json.NewEncoder(topicbuf).Encode(topicse); err != nil {
			fmt.Println(fmt.Sprintf("Error encoding query: %s", err))
		}
		buf.WriteString(fmt.Sprintf(
			",%s", topicbuf.String()))
	}
	buf.WriteString(`]}},`)
	buf.WriteString(fmt.Sprintf(`"size":%d}`, 10))
	var r map[string]interface{}

	req := esapi.SearchRequest{
		Index:  []string{"bbs"},
		Body:   strings.NewReader(buf.String()),
		Pretty: true,
	}
	// Perform the search request.
	res, err := req.Do(context.Background(), esC)
	if err != nil {
		fmt.Println(fmt.Sprintf("Error getting response: %s", err))
	}
	defer res.Body.Close()
```

上面的查询语法中，特别注意一点，**因为tags字段包含有子字段**，所以在查询的时候特别注意。

```
curl '127.0.0.1:9200/bbs/_search?pretty' -H 'Content-Type:application/json' -d '{"query":{"bool":{"must":[{"match":{"content":"拼多多"}}],"filter":[{"terms":{"tags":["科技"]}}]}},"size":2}'  
{
  "took" : 7,
  "timed_out" : false,
  "_shards" : {
    "total" : 5,
    "successful" : 5,
    "skipped" : 0,
    "failed" : 0
  },
  "hits" : {
    "total" : 0,
    "max_score" : null,
    "hits" : [ ]
  }
}
```

通过这次的搜索功能的开发，可以总结的地方，其实还有很多，因为时间的问题，后续有新的内容，在进行详细的总结，通过本次总结，可以发现：

- 在使用包含这种通过插入数据即可成功创建表结构的工具，一定要留意是否真的是我们想要的结果，最好封装出自己的专门做特地功能的工具，比如创建映射就专门用来创建映射，如果创建失败就报错，不要继续往下走，让工具来帮我们自动创建
- 好好学习语法，再细小的语法知识，都能在实战中给我们一痛击。
- 再接再厉
