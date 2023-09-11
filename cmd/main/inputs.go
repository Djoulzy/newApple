package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/mattn/go-tty"
)

var keyb *tty.TTY

func input() {
	for {
		r, _ := keyb.ReadRune()
		switch r {
		case 't':
			trace = !trace
		case ' ':
			trace = true
			stepper = true
			return
		case 'q':
			fmt.Printf("%s\n", cpu.FullDebug)
			os.Exit(0)
		}
	}
}

func InterractiveMode() bool {
	var addr string
	var endAddr bool

	for {
		r, _ := keyb.ReadRune()
		switch r {
		case ' ':
			return false
		case 'r':
			stepper = false
			return true
		case 'd':
			fmt.Printf("Dump> ")
			endAddr = false
			for !endAddr {
				rr, _ := keyb.ReadRune()
				switch rr {
				case 13:
					fmt.Println("")
					hx, _ := strconv.ParseInt(addr, 16, 64)
					MEM.Dump(uint16(hx))
					addr = ""
					endAddr = true
				default:
					fmt.Printf("%c", rr)
					addr += string(rr)
				}
			}
		case 'p':
			fmt.Printf("Disass> ")
			endAddr = false
			for !endAddr {
				rr, _ := keyb.ReadRune()
				switch rr {
				case 13:
					fmt.Println("")
					hx, _ := strconv.ParseInt(addr, 16, 64)
					cpu.Disassemble(uint16(hx), 40)
					addr = ""
					endAddr = true
					r = 'p'
				default:
					fmt.Printf("%c", rr)
					addr += string(rr)
				}
			}
		case 'q':
			fmt.Printf("%s\n", cpu.FullDebug)
			os.Exit(0)
		}
	}
}
