# Javascript AST scraper

[![GoDoc](https://godoc.org/github.com/jloup/scraper/js?status.svg)](https://godoc.org/github.com/jloup/scraper/js)

This package provides tools to easily scrap Javascript AST (using [github.com/robertkrimen/otto](https://github.com/robertkrimen/otto))

It can be used in combination of [github.com/jloup/scraper/html](https://github.com/jloup/scraper/tree/master/html) to scrap javascript contained in html script tags

A JSON config allows you to define what you want to scrap from javascript. Each node of your JSON should contain the following keys:
- [type](#node-types) : the js ast node type you want to scrap
- identifier (optional): the identifier of the node if any. An identifier is a name in js code (but not a literal like strings, number, ...). e.g. : var element = React.createElement("argument") -> 'element', 'React', 'createElement' are identifier for their respective node ; "argument" is a literal
- [extractors](#extractors) (optional): functions that scrap js node ast
- agg: aggregator type. See [aggregator explanation](https://github.com/jloup/scraper/tree/master/html#agg)

Example: 
this js code:
```js
(function() {

var element = myObject.play([
                              {"videoid": "Yfhuyr-T", "views": 10},
                              {"videoid": "dhy-Ggdz", "views": 356980}
                            ]);
})();
```
and the following scrap config:
```json
{
"config": [
          {
          "type": "var",
          "identifier": "element",
          "childs": [
                    {
                    "type": "*",
                    "identifier": "play",
                    "childs": [
                              {
                              "type": "[]",
                              "childs": [
                                        {
                                        "type": "{}",
                                        "agg": "1",
                                        "childs": [
                                                  {
                                                  "type": "*",
                                                  "identifier": "videoid",
                                                  "childs": [
                                                            {
                                                            "type": "literal",
                                                            "extractors": [{
                                                                          "name": "literal", 
                                                                          "key": "id"
                                                                          }]
                                                            }
                                                            ]
                                                  },
                                                  {
                                                  "type": "*",
                                                  "identifier": "views",
                                                  "childs": [
                                                            {
                                                            "type": "literal",
                                                            "extractors": [{
                                                                          "name": "literal",
                                                                          "key": "views"
                                                                          }]
                                                            }
                                                            ]
                                                  }
                                                  ]
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
will output this map[string]interface{}:
```go
map[videoid:"Yfhuyr-T" views: "10"]
map[videoid:"dhy-Ggdz" views:"356980"]
```

## Installation & use

Get the pkg
```
go get github.com/jloup/scraper/js
```

Use it in code
```
import "github.com/jloup/scraper/js"
```

## <a name="node-types"></a>Node types

The following values as "type" are handled:
- '{}': ObjectLiteral
- '[]': ArrayLiteral
- 'fn': FunctionLiteral
- 'var': VariableExpression
- 'string': StringLiteral
- 'bool': BooleanLiteral
- 'number': NumberLiteral
- 'literal': StringLiteral || BooleanLiteral || NumberLiteral || RegExpLiteral
- 'for': ForStatement
- 'new': NewExpression
- '*': will match any node type

If you are curious to know more about ast node type, you should see [github.com/robertkrimen/otto/ast](https://github.com/robertkrimen/otto/blob/master/ast/node.go)

## <a name="extractors"></a>Extractors

Extractors are implemented in [github.com/jloup/scraper/js/extractor](https://github.com/jloup/scraper/tree/master/js/extractor)

The following are currently available to use:
- extractor.Identifier : scrap identifier as string
```json
{"name": "identifier", "key": "<key store value>"}
```
- extractor.Literal : scrap literal value as string
```json
{"name": "literal", "key": "<key store value>"}
```

To make your own extractor accessible from JSON configuration files, you must register them with AddExtractorGenerator functions.
```go
js.AddExtractorGenerator("<name>" , <extractor>)
```

## Enhanced node types
As first example showed, the config can be really verbose. Some enhanced types allow the config to be lighter. 
The node must be passed an additional 'config' key containing a map[string]string:
- '{property:}': will extract config["property"] key from a js object
- '{}->map': will extract each config key from a js object (see example)
- to be continued...

Example using '{}->map' node type to make first example less verbose:

scrap config
```json
{
"config": [
          {
          "type": "var",
          "identifier": "element",
          "childs": [
                    {
                    "type": "*",
                    "identifier": "play",
                    "childs": [
                              {
                              "type": "[]",
                              "childs": [
                                        {
                                        "type": "{}->map",
                                        "config": {"videoid": "1", "views": "1"}
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
output:
```go
map[id:"Yfhuyr-T" views: "10"]
map[id:"dhy-Ggdz" views:"356980"]
```
