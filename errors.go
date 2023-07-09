package clear

import "fmt"

type insufficientArg struct {
	min     int
	current int
}

func (i *insufficientArg) Error() string {
	return fmt.Sprintf("not enough arguments need %d for %d", i.min, i.current)
}
