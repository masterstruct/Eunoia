package board

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// https://chessprogramming.org/Extended_Position_Description
// like above, but without any fancy shit like "opcode mnemonics"
// an EPD string MUST start with a FEN, followed by `;` separated
// key-value combos (value is optional, defaults to empty string)
// 4k3/8/8/8/8/8/8/4K2R w K - 0 1 ;D1 15 ;D2 66 ;D3 1197 ;D4 7059

// Example use:
// r := io.Reader...
// err := ParseEPD(r, func(epd EPD) error {
// 	fmt.Println(epd)
// 	return nil
// })

type EPD struct {
	FEN        string
	Operations []Operation
}

type Operation struct {
	Key   string
	Value string
}

func ParseEPD(r io.Reader, callback func(EPD) error) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	for lineNumber := 1; scanner.Scan(); lineNumber++ {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		epd, err := ParseEPDLine(line)
		if err != nil {
			return fmt.Errorf("epd: line %d: %w", lineNumber, err)
		}
		if err := callback(epd); err != nil {
			return fmt.Errorf("epd: line %d: %w", lineNumber, err)
		}
	}
	return scanner.Err()
}

func ParseEPDLine(line string) (EPD, error) {
	var epd EPD

	fen, operations, hasAny := strings.Cut(line, ";")
	epd.FEN = strings.TrimSpace(fen)
	if epd.FEN == "" {
		return epd, fmt.Errorf("empty FEN")
	}
	if !hasAny {
		return epd, nil
	}

	for raw := range strings.SplitSeq(operations, ";") {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}

		fields := strings.Fields(raw)
		key := fields[0]
		value := strings.TrimSpace(raw[len(key):])

		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) && len(value) >= 2 {
			value = value[1 : len(value)-1]
		}

		epd.Operations = append(epd.Operations, Operation{Key: key, Value: value})
	}

	return epd, nil
}
