# YAML 1.2 marshaling and unmarshaling support for Go

[![coverage](https://codecov.io/gh/stellirin/go-yaml/branch/main/graph/badge.svg?token=16jEi8Qbct)](https://codecov.io/gh/stellirin/go-yaml)
[![tests](https://github.com/stellirin/go-yaml/workflows/Go/badge.svg)](https://github.com/stellirin/go-yaml/actions?query=workflow%3AGo)

A simple package to marshal and unmarshal YAML 1.2 in Go.

## ⚙️ Installation

```sh
go get -u czechia.dev/yaml
```

## 📝 Introduction

This package is a rewrite of [github.com/ghodss/yaml](https://github.com/ghodss/yaml) using the new `yaml.Node` API in [yaml.v3](https://github.com/go-yaml/yaml/tree/v3). The API of *this* package is the same as ghodss/yaml so you can use it as a drop-in replacement wherever ghodss/yaml is used.

The original library was designed to enable a better way of handling YAML when marshaling to and from structs.

This library first converts YAML to JSON using yaml.v3 and then uses `json.Unmarshal` to convert to the desired struct. This means that it uses the same JSON struct tags as well as the custom JSON methods `MarshalJSON` and `UnmarshalJSON`.

For a detailed overview of the reasoning behind this method, [see this archived blog post](http://web.archive.org/web/20190603050330/http://ghodss.com/2014/the-right-way-to-handle-yaml-in-golang/).

## 👍 Compatibility

This package uses [yaml.v3](https://github.com/go-yaml/yaml/tree/v3) and therefore supports [everything yaml.v3 supports](https://github.com/go-yaml/yaml/tree/v3#compatibility).

The most noticable difference between yaml.v2 and yaml.v3 (or rather the YAML 1.1 and YAML 1.2 specifications) is around booleans. Specifically only variants of `true` and `false` are booleans, so `y` is now treated as a string.



## ⛔️ Caveats

**Caveat #1:**  When using `yaml.Marshal` and `yaml.Unmarshal`, binary data should **NOT** be preceded with the `!!binary` YAML tag. If you do, yaml.v3 will convert the binary data from base64 to native binary data, which is not compatible with JSON. You can still use binary in your YAML files - just store them without the `!!binary` tag and decode the base64 in your code (e.g. in the custom JSON methods `MarshalJSON` and `UnmarshalJSON`). This has the benefit that binary data will be decoded exactly the same way for both YAML and JSON.

**Caveat #2:** When using `YAMLToJSON` directly, maps with keys that are maps will result in an error since this is not supported by JSON. This error will occur in `Unmarshal` as well since you can't unmarshal map keys anyways since struct fields can't be keys.

## 💻 Usage

Import using:

```go
import "czechia.dev/yaml"
```

Usage is very similar to the JSON library:

```go
package main

import (
	"fmt"

	"czechia.dev/yaml"
)

type Person struct {
	Name string `json:"name"` // Affects YAML field names too.
	Age  int    `json:"age"`
}

func main() {
	// Marshal a Person struct to YAML.
	p := Person{"John", 30}
	y, err := yaml.Marshal(p)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(string(y))
	/* Output:
	name: John
	age: 30
	*/

	// Unmarshal the YAML back into a Person struct.
	var p2 Person
	err = yaml.Unmarshal(y, &p2)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(p2)
	/* Output:
	{John 30}
	*/
}
```

The intermediate `yaml.YAMLToJSON` and `yaml.JSONToYAML` methods are also available:

```go
package main

import (
	"fmt"

	"czechia.dev/yaml"
)

func main() {
	j := []byte(`{"name": "John", "age": 30}`)
	y, err := yaml.JSONToYAML(j)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(string(y))
	/* Output:
	age: 30
	name: John
	*/

	j2, err := yaml.YAMLToJSON(y)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Println(string(j2))
	/* Output:
	{"age":30,"name":"John"}
	*/
}
```
