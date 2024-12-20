# ValueBox

## Overview
ValueBox is a Golang module designed to provide a simple and efficient way to handle and manipulate JSON structures.

## Installation
To install ValueBox, use the following command:

```sh
go get github.com/duhnnie/valuebox
```

## Example of Use
Here is a basic example of how to use ValueBox in your project:

```go
package main

import (
	"fmt"
	"strings"

	"github.com/duhnnie/godash"
	"github.com/duhnnie/valuebox"
)

func main() {
	b := valuebox.New()

	j := []byte(`{
		"name": "Pepito",
		"age": 20
	}`)

	b.Set("myJSON", j)
	b.Set("myJSON.name", []byte("Pepito Watson"))
	b.Set("myJSON.isStudent", []byte("true"))
	b.Set("myJSON.hobbies", []byte(`["music", "programming", "outdoors"]`))
	b.Set("myJSON.grades", []byte(`{"maths": 10, "english": 9, "science": 8}`))

	name, _ := b.GetString("myJSON.name")
	age, _ := b.GetFloat64("myJSON.age")
	hobbies, _ := b.GetStringSlice("myJSON.hobbies")
	grades, _ := b.GetFloat64Map("myJSON.grades")

	gradesAvg := godash.ReduceMap(grades, func(acc float64, _ string, grade float64, _ map[string]float64) float64 {
		return acc + grade
	}, 0.0) / float64(len(grades))

	summary := fmt.Sprintf(
		"\"Hi, my name is %s, I'm %d, my hobbies are %s, my grades average is %0.2f\"",
		name,
		int(age)+1,
		strings.Join(hobbies, ", "),
		gradesAvg,
	)

	_ = b.Set("myJSON.summary", []byte(summary)) 
	j, _ = b.ToJSON() 
	// You can also use b.ValueToJSON("myJSON") for getting only that part of the JSON

	fmt.Println(string(j))
}
```

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Credits
Developed by:
- Daniel Canedo Ramos (@duhnnie)

Special thanks to all contributors and the open-source community.