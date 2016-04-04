# HTML scraper
[![GoDoc](https://godoc.org/github.com/jloup/scraper/html?status.svg)](https://godoc.org/github.com/jloup/scraper/html)

Package scraper provides tools to easily scrap HTML pages. 

A JSON config allows you to define what you want to scrap from HTML. There are three main directives the parser understands:
- tag : which HTML tag you are interested in ('div', 'a',...)
- validators : functions that return if a HTML tag is something of interest
- extractors: functions that scrap a HTML tag attributes and content

Config example:
```json
{
    "config": [
    {
     "tag": "ul",
     "validators": [
                {
                    "name": "exists",
                    "attr": "id"
                },
                {
                    "name": "attrEquals",
                    "attr": "id",
                    "value": "item-list"
                }
            ],
            "childs": [
                {
                    "tag": "li",
                    "validators": [
                        {
                            "name": "exists",
                            "attr": "class"
                        },
                        {
                            "name": "attrEquals",
                            "attr": "class",
                            "value": "item-name"
                        }
                    ],
                    "extractors": [
                        {
                            "name": "extractAttr",
                            "attr": "data-item-category"
                        }
                    ],
                    "childs": [
                        {
                            "tag": "text",
                            "extractors": [
                                {
                                    "name": "extractText",
                                    "key": "item-name"
                                }
                            ]
                        }
                    ]
                }
            ]
        }
    ]
}
```

that could output a map:
```go
 []map[string]interface{}{
                           {"data-item-category": "food", "item-name": "tomato"},
                           {"data-item-category": "cosmetic", "item-name": "shampoo"},
}
```

#### Real world example : scraping Soundcloud iframes

with this config soundcloud.json:
```json
{"config": [
  { "tag": "iframe",
    "validators": [
              {"name": "exists", "attr": "src"},
              {"name": "scIFast"},
              {"name": "scI"}
             ],
    "extractors": [
               {"name": "sc", "attr": "src"}
             ]
  },
  {
    "tag": "object",
    "childs": [
               {
                "tag": "param",
                "validators": [
                          {"name": "scOFast"},
                          {"name": "scO"}
                         ],
                "extractors": [
                           {"name": "sc", "attr": "value"}
                         ]
               }
             ]
  }
]}
```

and the Go code:
```go
    // errors are not checked for reading ease
	s, _ := html.JSONFileToScraper("soundcloud.json")
	
	// http://www.lagasta.com (music blog)
	f, _ := os.Open("testdata/exampledata/website.html")

	items, _ := html.ScrapHTML(s, f)

	for _, item := range items {
		fmt.Printf("Soundcloud type '%s' id %s\n", item["type"], item["id"])
	}
```

Ouput:
```
Soundcloud type 'tracks' id 196618442
Soundcloud type 'tracks' id 195052762
Soundcloud type 'tracks' id 195529989
Soundcloud type 'tracks' id 196295430
Soundcloud type 'tracks' id 187956230
Soundcloud type 'tracks' id 196308250
Soundcloud type 'tracks' id 195962464
Soundcloud type 'tracks' id 196130447
```

## Installation & Use

Get the pkg
```
go get github.com/jloup/scraper/html
```

Use it in code
```
import "github.com/jloup/scraper/html"
```

documentation on [godoc](https://godoc.org/github.com/jloup/scraper/html)

## Validators and Extratcors
Validators and Extractors are implemented in their respective packages: 
- [http://github.com/jloup/scraper/html/validator](https://godoc.org/github.com/jloup/scraper/html/validator)
- [http://github.com/jloup/scraper/html/extractor](https://godoc.org/github.com/jloup/scraper/html/extractor)

Current list of validators and how to use them in JSON: 

- validator.Exists : check that a tag attribute exists 
```json
{"name": "exists", "attr": "<attribute-name>"}
```
- validator.AttrEquals : check that a tag attribute matches a value.
```json
{"name": "attrEquals", "attr": "<attribute name>", "value": "<value to match>"}
```
- validator.AttrContains : check that a tag attribute value passes strings.Contains.
```json
{"name": "attrContains", "attr": "<attribute name>", "value": "<value to test>"}
```
Current list of extractors and how to use it in JSON: 

- extractor.Attribute : scrap attribute value 
```json
{"name": "extractAttr", "attr": "<attribute-name>"}
```
- extractor.TextContent : extract text content of a text node ("tag": "text")
```json
{"name": "extractText", "key": "<store key>"}
```
- [extractor.Js](#javascript) : scrap javascript ast
```json
{"name": "js", "config": "<js scraper name>"}
```

For further documentation, see validator/extractor packages doc. They serve as examples for your own implementation. To make your own validator/extractor accessible from JSON configuration files, you must register them with AddValidatorGenerator and AddExtractorGenerator functions.
```go
html.AddValidatorGenerator("<name>" , <validator>)
html.AddExtractorGenerator("<name>" , <extractor>)
```

## <a name="javascript"></a>Javascript
Javascript code in HTML (e.g. script tag) can be scraped as well using js AST parsing. Note that javascript code is not interpreted, the js AST is just walked to extract objects, variables, arrays,...

You must define your js scraper using [github.com/jloup/scraper/js](https://github.com/jloup/scraper/tree/master/js) and register it in the module registry. You scraper can then be used with extractor.Js.

Refer to package documentation for example.

##<a name="agg"></a> Aggregator

There is a fourth directive you can pass to the parser to reflect how the data is structured in HTML. In JSON, you can add a 'type' attribute at the same level the 'tag' attribute is. It can be one of the following:

#### "type": "N" (default)
"N" node duplicates its own data for each of its children

Example:
```html
<div id="dad" name="Jean"> <!-- parent type 'N' -->
   <li class="child">Pierre</li> <!-- child 1 -->
   <li class="child">Romain</li> <!-- child  2-->
 </div>`
```
Output:
```go
[]map[string]interface{}{
  {"parent": "Jean", "child": "Pierre"},
  {"parent": "Jean", "child": "Romain"},
}
```
#### "type": "list"
"list" node indexes data that share the same 'key'
Example:
```html
<div id="dad" name="Jean"> <!-- parent type 'list' -->
   <li class="child">Pierre</li> <!-- child 1 -->
   <li class="child">Romain</li> <!-- child  2-->
 </div>`
```
Output:
```go
[]map[string]interface{}{
  {"parent": "Jean", "child0": "Pierre", "child1": "Romain"},
}
```
#### "type": "array"
"array" gathers node data in an array (array key set by user)
Example:
```html
<div id="dad" name="Jean"> <!-- parent type 'list' -->
   <li class="child">Pierre</li> <!-- child 1 -->
   <li class="child">Romain</li> <!-- child 2-->
 </div>`
```
Output:
```go
[]map[string]interface{}{
  {"parent": "Jean", "array": ["Pierre", "Romain"]},
}
```

#### "type": "1"
"1" node aggregates data from its children in an unique hash.

Example:
```html
<div id="SELL455225" class="house"> <!-- parent "1" -->
  <div class="price">100.000$</div> <!-- child -->
  <div class="location">
    <span class="state">TX</span> <!-- child -->
    <span class="city">Austin</span> <!-- child -->
  </div>
</div>
```
Output:
```go
[]map[string]interface{}{
  {"id": "SELL455225", "price": "100,000$", "state": "TX", "city": "Austin"},
}
```
