package errorcheck

import (
	"fmt"
)

func Check(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
