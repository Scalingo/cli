package term

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"syscall"
	"unicode/utf8"
	"unsafe"

	"gopkg.in/errgo.v1"
)

type wBool uint16

type inputRecord struct {
	EventType       uint16
	Unknown1        uint16
	KeyDown         wBool
	RepeatCount     uint16
	Unknown2        uint16
	VirtualKeyCode  uint16
	VirtualScanCode uint16
	UnicodeChar     [2]byte
	ControlKeyState uint16
}

const (
	inputRecordSize = unsafe.Sizeof(inputRecord{})
)

var (
	kernel32, _         = syscall.LoadLibrary("kernel32.dll")
	getStdHandle, _     = syscall.GetProcAddress(kernel32, "GetStdHandle")
	readConsoleInput, _ = syscall.GetProcAddress(kernel32, "ReadConsoleInputW")
	stdInputHandle      = int32(-10)
)

func Password(prompt string) (string, error) {
	defer syscall.FreeLibrary(kernel32)
	fmt.Print(prompt)

	var pass string
	var c rune
	for c != '\r' {
		var err error
		c, err = getChar()
		if err != nil {
			return "", err
		}
		fmt.Print("*")
		if c != '\r' {
			pass += string(c)
		}
	}
	fmt.Println()
	return pass, nil
}

func getChar() (rune, error) {
	handlePtr, _, errno := syscall.Syscall(uintptr(getStdHandle), uintptr(1), uintptr(stdInputHandle), uintptr(0), uintptr(0))
	if errno != 0 {
		return 0, fmt.Errorf("Fail to get terminal handle (%v)", errno)
	}

	var buffer [inputRecordSize]byte
	var nbEvents int32
	_, _, errno = syscall.Syscall6(uintptr(readConsoleInput), uintptr(4), handlePtr, uintptr(unsafe.Pointer(&buffer[0])), uintptr(1), uintptr(unsafe.Pointer(&nbEvents)), 0, 0)
	if errno != 0 {
		return 0, fmt.Errorf("Fail to read character with win#ReadConsoleInput (%v)", errno)
	}
	b := bytes.NewBuffer(buffer[:])
	input := &inputRecord{}
	err := binary.Read(b, binary.LittleEndian, input)
	if err != nil {
		return 0, errgo.Mask(err, errgo.Any)
	}
	// Remove keyUp events, take only keyboard event, and ctrl/alt/cap keys
	char, _ := utf8.DecodeRune(input.UnicodeChar[:])
	if input.KeyDown == 0 || input.EventType != 1 || char == 0 {
		return getChar()
	}
	return char, nil
}
