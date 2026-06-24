package fusion

import (
	"admin/fusion/base"
	. "admin/fusion/buildin"
	"strconv"
	"strings"
)

type TextUnpacker struct {
	text string
	pos  int
}

func NewTextUnpacker(text string) *TextUnpacker {
	return &TextUnpacker{text, 0}
}

func (self *TextUnpacker) IsEmpty() bool {
	return self.pos >= len(self.text)
}

func (self *TextUnpacker) IsDelimiter(delimiter rune) bool {
	return self.pos > 0 && rune(self.text[self.pos-1]) == delimiter
}

func (self *TextUnpacker) IsAnchor(anchor rune) bool {
	if self.pos < len(self.text) && rune(self.text[self.pos]) == anchor {
		self.pos += 1
		return true
	}
	return false
}

func (self *TextUnpacker) UnpackBool() bool {
	stop := findNumberStop(self.text, self.pos)
	value, _ := strconv.ParseBool(self.text[self.pos:stop])
	self.pos = base.MinInt(len(self.text), stop+1)
	return value
}

func (self *TextUnpacker) UnpackInt() int64 {
	stop := findIntegerStop(self.text, self.pos)
	value, _ := strconv.ParseInt(self.text[self.pos:stop], 10, 64)
	self.pos = base.MinInt(len(self.text), stop+1)
	return value
}

func (self *TextUnpacker) UnpackUint() uint64 {
	stop := findNumberStop(self.text, self.pos)
	value, _ := strconv.ParseUint(self.text[self.pos:stop], 10, 64)
	self.pos = base.MinInt(len(self.text), stop+1)
	return value
}

func (self *TextUnpacker) UnpackFloat() float64 {
	stop := findFloatStop(self.text, self.pos)
	value, _ := strconv.ParseFloat(self.text[self.pos:stop], 64)
	self.pos = base.MinInt(len(self.text), stop+1)
	return value
}

func (self *TextUnpacker) UnpackString() string {
	var value string
	if bytes := int(self.UnpackUint()); bytes > 0 {
		if self.pos+bytes+1 <= len(self.text) {
			value = self.text[self.pos : self.pos+bytes]
			self.pos += bytes + 1
		} else {
			value = self.text[self.pos:]
			self.pos = len(self.text)
		}
	}
	return value
}

func (self *TextUnpacker) UnpackXString(delimiter rune) string {
	var value string
	if bytes := strings.IndexRune(self.text[self.pos:], delimiter); bytes != -1 {
		value = self.text[self.pos : self.pos+bytes]
		self.pos += bytes + 1
	} else {
		value = self.text[self.pos:]
		self.pos = len(self.text)
	}
	return value
}

func findNumberStop(text string, start int) int {
	var pos = start
	for len := len(text); pos < len; pos++ {
		ch := text[pos]
		if ch < '0' || ch > '9' {
			break
		}
	}
	return pos
}

func findIntegerStop(text string, start int) int {
	return findNumberStop(text, start+IfInt(isSign(text, start), 1, 0))
}

func findFloatStop(text string, start int) int {
	pos := findIntegerStop(text, start)
	if isDot(text, pos) {
		pos = findNumberStop(text, pos+1)
	}
	if isExp(text, pos) {
		pos = findIntegerStop(text, pos+1)
	}
	return pos
}

func isSign(text string, pos int) bool {
	if pos < len(text) {
		return text[pos] == '-'
	}
	return false
}
func isExp(text string, pos int) bool {
	if pos < len(text) {
		ch := text[pos]
		return ch == 'e' || ch == 'E'
	}
	return false
}
func isDot(text string, pos int) bool {
	if pos < len(text) {
		return text[pos] == '.'
	}
	return false
}
