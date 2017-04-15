package d71

import "fmt"

var (
	ErrDiskFull = fmt.Errorf("disk full")
	ErrDirFull  = fmt.Errorf("directory full")
)
