# JSON Pointer

Parse struct according to jsonpointer

## Use

> go get -u github.com/zc2638/jsonpointer

## Example

```go
import (
    "github.com/zc2638/jsonpointer"
)

type Test struct {
    Name     string  `json:"name"`
    Children []Child `json:"children"`
}

type Child struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
    Data *Data  `json:"data"`
}

type Data struct {
    Message string            `json:"message"`
    Labels  map[string]string `json:"labels"`
    Line    *int              `json:"line"`
}

func main() {
    line := 5
    data := Test{
        Name: "张三",
        Children: []Child{
        	{
                Name: "李四",
                Age:  20,
                Data: &Data{
    	            Message: "a12312",
                    Labels: map[string]string{
                        "a/b": "data",
                    },
                    Line: &line,
                },
            },
        },
    }
	
    parser, err := jsonpointer.NewParser(data)	
    if err != nil {
        log.Fatal(err)	
    }
    
    exists := parser.Check("/children/0/age")
    fmt.Printf("path exists: %s\n", exists)e
    
    value, err := parser.Get("/children/0/age")
    if err != nil {
        log.Fatal(err)	
    }
    fmt.Println(value == 20)
}
```