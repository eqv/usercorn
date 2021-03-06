package syscalls

import (
	"fmt"
	"os"
	"strings"

	"../models"
)

const (
	INT = iota
	FD
	STR
	BUF
	OBUF
	LEN
	OFF
	PTR
)

func traceBasicArg(u models.Usercorn, name string, arg uint64, t int) string {
	switch t {
	case INT, FD:
		return fmt.Sprintf("%d", int32(arg))
	case STR:
		s, _ := u.Mem().ReadStrAt(arg)
		return fmt.Sprintf("%#v", s)
	default:
		return fmt.Sprintf("0x%x", arg)
	}
}

func traceArg(u models.Usercorn, name string, args []uint64, t int) string {
	switch t {
	case BUF:
		mem, _ := u.MemRead(args[0], args[1])
		return fmt.Sprintf("%#v", string(mem))
	default:
		return traceBasicArg(u, name, args[0], t)
	}
}

func traceArgs(u models.Usercorn, name string, args []uint64) string {
	types := syscalls[name].Args
	ret := make([]string, 0, len(types))
	for i, t := range types {
		s := traceArg(u, name, args[i:], t)
		ret = append(ret, s)
	}
	return strings.Join(ret, ", ")
}

func Trace(u models.Usercorn, name string, args []uint64) {
	fmt.Fprintf(os.Stderr, "%s(%s)", name, traceArgs(u, name, args))
}

func TraceRet(u models.Usercorn, name string, args []uint64, ret uint64) {
	types := syscalls[name].Args
	var out []string
	for i, t := range types {
		if t == OBUF {
			r := int(ret)
			if uint64(r) <= args[i+1] && r >= 0 {
				mem, _ := u.MemRead(args[i], uint64(r))
				out = append(out, fmt.Sprintf("%#v", string(mem)))
			}
		}
	}
	out = append(out, traceBasicArg(u, name, ret, syscalls[name].Ret))
	fmt.Fprintf(os.Stderr, " = %s\n", strings.Join(out, ", "))
}
