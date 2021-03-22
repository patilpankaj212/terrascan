package utils

import (
	"bytes"
	"runtime"
	"strings"
)

func IsWindowsPlatform() bool {
	return runtime.GOOS == "windows"
}

func ReplaceCarriageReturnBytes(input []byte) []byte {
	return bytes.ReplaceAll(input, []byte("\r\n"), []byte("\n"))
}

func ReplaceCarriageReturnString(input string) string {
	return strings.ReplaceAll(input, "\r\n", "\n")
}
