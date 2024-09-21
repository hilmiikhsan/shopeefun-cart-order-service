package helper

import (
	"fmt"
	"strings"
)

func QueryLog(query string, args ...interface{}) {
	for i, v := range args {
		query = strings.ReplaceAll(query, fmt.Sprintf("$%d", (i+1)), fmt.Sprintf("'%v'", v))
	}
	fmt.Println(query)
}
