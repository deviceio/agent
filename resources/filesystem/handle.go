package filesystem

import "os"

type handle struct {
	file *os.File
}
