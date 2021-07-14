# jsont(ype)

A JSON parser returns the native type of JSON following ECMA-404 standard.

Passed all test from <https://github.com/nst/JSONTestSuite>

## How to use

```bash
go get github.com/chikaku/jsont
```

### Types in jsont

- JSON type is one of Object, Array, String, Number, Boolean, Null
- Object is `map[string]JSON`
- Array is `[]JSON`
- String is `string`
- Number is `string` it can format to int64 or float64
- Boolean is `true` or `false`
- Null is `struct{}`

### Decode complete JSON

```go
var text = []byte(`"animals": ["🐱", "🐶", "🐭"]`)
result, _ := jsont.Decode(text)             // result is map[string]JSON
animals := result["animals"]                // animals is Array []JSON
cat, dog := animals[0], animals[1]          // cat, dog is string
```

### Decode part of JSON

Decode one JSON element and return position

```go
obj, pos, err = ReadObject([]byte(`  {"a": "b"}...`))    // {"a": "b"}
arr, pos, err = ReadArray([]byte(`  [1, {}, ""],...`))   // [1, {}, ""]
str, pos, err = ReadString([]byte(`  "   1""...`))       // "   1"
numstr, pos, err = ReadNumber([]byte(`  177178...`))     // 177178
num, err = num.Int64()                                   // 177178
_, _, err = ReadNull([]byte(`  null\...`))
_, _, err = ReadBool([]byte(`  false\...`))
```
