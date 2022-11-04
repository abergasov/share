package utils

import (
	"fmt"
	"io"
	"os"
	"strings"

	"rsc.io/qr"
)

const (
	blackWhite = "▄"
	blackBlack = " "
	whiteBlack = "▀"
	whiteWhite = "█"
	// size of border in full blocks
	borderSize = 2
)

func PrintQRCode(data string) error {
	code, err := qr.Encode(data, qr.Q)
	if err != nil {
		return fmt.Errorf("failed to encode qr code: %w", err)
	}
	return writeHalfBlocks(os.Stdout, code, borderSize)
}

func writeHalfBlocks(w io.Writer, code *qr.Code, borderSize int) error {
	var res strings.Builder
	// top border
	if borderSize%2 != 0 {
		res.WriteString(stringRepeat(blackWhite, code.Size+borderSize*2) + "\n")
		res.WriteString(stringRepeat(stringRepeat(whiteWhite, code.Size+borderSize*2)+"\n", borderSize/2))
	} else {
		res.WriteString(stringRepeat(stringRepeat(whiteWhite, code.Size+borderSize*2)+"\n", borderSize/2))
	}
	for i := 0; i <= code.Size; i += 2 {
		res.WriteString(stringRepeat(whiteWhite, borderSize))
		for j := 0; j <= code.Size; j++ {
			nextBlack := false
			if i+1 < code.Size {
				nextBlack = code.Black(j, i+1)
			}
			currBlack := code.Black(j, i)
			if currBlack && nextBlack {
				res.WriteString(blackBlack)
			} else if currBlack && !nextBlack {
				res.WriteString(blackWhite)
			} else if !currBlack && !nextBlack {
				res.WriteString(whiteWhite)
			} else {
				res.WriteString(whiteBlack)
			}
		}
		res.WriteString(stringRepeat(whiteWhite, borderSize-1) + "\n")
	}
	// bottom border
	if borderSize%2 == 0 {
		res.WriteString(stringRepeat(stringRepeat(whiteWhite, code.Size+borderSize*2)+"\n", borderSize/2-1))
		res.WriteString(stringRepeat(whiteBlack, code.Size+borderSize*2) + "\n")
	} else {
		res.WriteString(stringRepeat(stringRepeat(whiteWhite, code.Size+borderSize*2)+"\n", borderSize/2))
	}
	_, err := w.Write([]byte(res.String()))
	return err
}

func stringRepeat(s string, count int) string {
	if count <= 0 {
		return ""
	}
	return strings.Repeat(s, count)
}
