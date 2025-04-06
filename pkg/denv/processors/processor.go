package processors

import "os"

type Processor interface {
	Run(options any, stdin *os.File, stdout *os.File, stderr *os.File) error
}
