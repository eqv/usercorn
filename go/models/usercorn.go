package models

import (
	"encoding/binary"
	"github.com/lunixbochs/ghostrace/ghost/memio"
	uc "github.com/unicorn-engine/unicorn/bindings/go/unicorn"
)

type Usercorn interface {
	uc.Unicorn
	Arch() *Arch
	Bits() uint
	ByteOrder() binary.ByteOrder
	Disas(addr, size uint64) (string, error)
	Symbolicate(addr uint64) (string, error)

	Brk(addr uint64) (uint64, error)
	Mmap(addr, size uint64) (uint64, error)
	MmapWrite(addr uint64, p []byte) (uint64, error)
	Mem() memio.MemIO

	PackAddr(buf []byte, n uint64) error
	UnpackAddr(buf []byte) uint64
	PopBytes(p []byte) error
	PushBytes(p []byte) (uint64, error)
	Pop() (uint64, error)
	Push(n uint64) (uint64, error)
	ReadRegs(reg []int) ([]uint64, error)
	RegDump() ([]RegVal, error)

	Loader() Loader
	InterpBase() uint64
	Base() uint64
	Entry() uint64
	BinEntry() uint64

	PrefixPath(s string, force bool) string
	PosixInit(args, env []string, auxv []byte) error
	Syscall(num int, name string, getArgs func(n int) ([]uint64, error)) (uint64, error)
	Exit(status int)
}
