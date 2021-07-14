package jsont

import (
	"strconv"
	"strings"
)

func Decode(raw []byte) (JSON, error) {
	val, pos, err := ReadOffset(raw)
	if err != nil {
		return nil, errWithPosition(pos, err)
	}

	for index, ch := range raw[pos:] {
		switch ch {
		case ' ', '\r', '\n', '\t':
		default:
			return nil, errWithPosition(pos+index+1, ErrUnexpectedCharacter)
		}
	}

	return val, nil
}

func ReadOffset(raw []byte) (JSON, int, error) {
	for index, ch := range raw {
		switch ch {
		case ' ', '\r', '\n', '\t':
			continue
		case '{':
			val, count, err := ReadObject(raw[index:])
			return val, index + count, err
		case '[':
			val, count, err := ReadArray(raw[index:])
			return val, index + count, err
		case '"':
			val, count, err := ReadString(raw[index:])
			return val, index + count, err
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '+', '-':
			val, count, err := ReadNumber(raw[index:])
			return val, index + count, err
		case 'n':
			val, count, err := ReadNull(raw[index:])
			return val, index + count, err
		case 't', 'f':
			val, count, err := ReadBool(raw[index:])
			return val, index + count, err
		default:
			return nil, 0, ErrUnexpectedCharacter
		}
	}

	return nil, 0, ErrTooShort
}

const (
	stateLookupKey = iota
	stateReadingKey
	stateLookupColon
	stateLookupCommon
)

func ReadObject(raw []byte) (Object, int, error) {
	index := skipWhitespace(raw)
	if index > len(raw)-1 || raw[index] != '{' {
		return nil, 0, ErrUnexpectedCharacter
	}

	m := Object{}
	noElement := true
	key := new(strings.Builder)
	state := stateLookupKey

	var ch byte
	for index < len(raw)-1 {
		index++
		ch = raw[index]

		switch state {
		case stateLookupKey:
			switch ch {
			case ' ', '\r', '\n', '\t':
			case '}':
				if noElement {
					return m, index + 1, nil
				}
				return nil, index + 1, ErrUnexpectedCharacter
			case '"':
				state = stateReadingKey
				noElement = false
			default:
				return nil, index + 1, ErrUnexpectedCharacter
			}
		case stateReadingKey:
			switch ch {
			case '"':
				state = stateLookupColon
			case '\n', '\r':
				return nil, index + 1, ErrUnexpectedNewline
			default:
				key.WriteByte(ch)
			}
		case stateLookupColon:
			switch ch {
			case ' ':
			case ':':
				val, count, err := ReadOffset(raw[index+1:])
				if err != nil {
					return nil, index + 1 + count, err
				}

				m[key.String()] = val
				key.Reset()
				index += count
				state = stateLookupCommon
			default:
				return nil, index + 1, ErrUnexpectedCharacter
			}
		case stateLookupCommon:
			switch ch {
			case ' ', '\r', '\n', '\t':
			case '}':
				return m, index + 1, nil
			case ',':
				state = stateLookupKey
			default:
				return nil, index + 1, ErrUnexpectedCharacter
			}
		}
	}

	return nil, index + 1, ErrUnclosed
}

func ReadArray(raw []byte) (Array, int, error) {
	array := make([]JSON, 0, 8)
	lookupCommon := false
	noElement := true

	index := skipWhitespace(raw)
	if index > len(raw)-1 || raw[index] != '[' {
		return nil, 0, ErrUnexpectedCharacter
	}

	var ch byte
	for index < len(raw)-1 {
		index++
		ch = raw[index]
		if lookupCommon {
			switch ch {
			case ' ', '\n', '\r', '\t':
			case ',':
				lookupCommon = false
			case ']':
				return array, index + 1, nil
			default:
				return nil, index + 1, ErrUnexpectedCharacter
			}
			continue
		}

		if ch == ']' && noElement {
			return array, index + 1, nil
		}

		val, count, err := ReadOffset(raw[index:])
		if err != nil {
			return nil, index + 1 + count, err
		}

		array = append(array, val)
		index += count - 1
		lookupCommon = true
		noElement = false
	}

	return nil, index + 1, ErrUnclosed
}

func ReadString(raw []byte) (String, int, error) {
	index := skipWhitespace(raw)
	if index > len(raw)-1 || raw[index] != '"' {
		return "", 0, ErrUnexpectedCharacter
	}

	index++
	var ch byte
	var escape bool
	bs := new(strings.Builder)
	for ; index < len(raw); index++ {
		ch = raw[index]
		if escape {
			switch ch {
			case 'n':
				bs.WriteByte('\n')
			case 'r':
				bs.WriteByte('\r')
			case 't':
				bs.WriteByte('\t')
			case '\\':
				bs.WriteByte('\\')
			case 'b':
				bs.WriteByte('\b')
			case 'f':
				bs.WriteByte('\f')
			case '"':
				bs.WriteByte('"')
			case '/':
				bs.WriteByte('/')
			case 'u':
				if index+4 >= len(raw) {
					return "", index + 1, ErrUnexpectedCharacter
				}
				num, err := strconv.ParseInt(string(raw[index+1:index+5]), 16, 32)
				if err != nil {
					return "", index + 5, ErrIllegalUnicode
				}

				bs.WriteRune(rune(num))
				index += 4
			default:
				return "", index + 1, ErrIllegalEscape
			}
			escape = false
			continue
		}

		switch ch {
		case '\\':
			escape = true
		case '\n', '\t', 0:
			return "", index + 1, ErrUnexpectedCharacter
		case '"':
			return String(bs.String()), index + 1, nil
		default:
			bs.WriteByte(ch)
		}
	}

	return "", index + 1, ErrUnclosed
}

const (
	stateLookupSign = iota
	stateLookupFirstDigit
	stateDigitBeforeDot
	stateFirstDigitAfterDot
	stateDigitAfterDot
	stateESign
	stateEFirstDigit
	stateEDigit
)

func ReadNumber(raw []byte) (Number, int, error) {
	var ch byte
	var state = stateLookupSign
	var numstr = new(strings.Builder)
	var index = skipWhitespace(raw) - 1

L:
	for index < len(raw)-1 {
		index++
		ch = raw[index]
		switch state {
		case stateLookupSign:
			if ch == '-' {
				numstr.WriteByte(ch)
			} else {
				index--
			}
			state = stateLookupFirstDigit
		case stateLookupFirstDigit:
			switch ch {
			case '0':
				numstr.WriteByte('0')
				index++
				if index == len(raw) { // finish
					break L
				}
				ch = raw[index]
				switch ch {
				case '.':
					numstr.WriteByte(ch)
					state = stateFirstDigitAfterDot
				case 'e', 'E':
					numstr.WriteByte(ch)
					state = stateESign
				default:
					index--
					break L
				}
				if index >= len(raw)-1 {
					return "", index + 1, ErrUnexpectedCharacter
				}
			case '1', '2', '3', '4', '5', '6', '7', '8', '9':
				numstr.WriteByte(ch)
				state = stateDigitBeforeDot
			default:
				return "", index + 1, ErrUnexpectedCharacter
			}
		case stateDigitBeforeDot:
			switch ch {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				numstr.WriteByte(ch)
			case '.':
				numstr.WriteByte(ch)
				state = stateFirstDigitAfterDot
			case 'e', 'E':
				numstr.WriteByte(ch)
				state = stateESign
			default:
				index--
				break L
			}
		case stateFirstDigitAfterDot:
			if '0' <= ch && ch <= '9' {
				numstr.WriteByte(ch)
				state = stateDigitAfterDot
			} else {
				index--
				break L
			}
		case stateDigitAfterDot:
			if '0' <= ch && ch <= '9' {
				numstr.WriteByte(ch)
			} else if ch == 'e' || ch == 'E' {
				numstr.WriteByte(ch)
				state = stateESign
			} else {
				index--
				break L
			}
		case stateESign:
			if '0' <= ch && ch <= '9' {
				numstr.WriteByte(ch)
				state = stateEDigit
			} else if ch == '+' || ch == '-' {
				numstr.WriteByte(ch)
				state = stateEFirstDigit
			} else {
				return "", index + 1, ErrUnexpectedCharacter
			}
		case stateEFirstDigit:
			if '0' <= ch && ch <= '9' {
				numstr.WriteByte(ch)
				state = stateEDigit
			} else {
				return "", index + 1, ErrUnexpectedCharacter
			}
		case stateEDigit:
			if '0' <= ch && ch <= '9' {
				numstr.WriteByte(ch)
			} else {
				index--
				break L
			}
		}
	}

	if state == stateFirstDigitAfterDot || state == stateESign || state == stateEFirstDigit {
		return "", index, ErrUnexpectedFinished
	}

	return Number(numstr.String()), index + 1, nil
}

func ReadNull(raw []byte) (Null, int, error) {
	index := skipWhitespace(raw)
	raw = raw[index:]

	if len(raw) < 4 {
		return Null{}, 0, ErrTooShort
	}

	if string(raw[:4]) != "null" {
		return Null{}, 0, ErrUnexpectedCharacter
	}

	return Null{}, 4, nil
}

func ReadBool(raw []byte) (Bollean, int, error) {
	index := skipWhitespace(raw)
	raw = raw[index:]

	if len(raw) < 4 {
		return false, len(raw), ErrTooShort
	}

	if len(raw) >= 4 && string(raw[:4]) == "true" {
		return true, 4, nil
	}

	if len(raw) >= 5 && string(raw[:5]) == "false" {
		return false, 5, nil
	}

	return false, 5, ErrUnexpectedCharacter
}

func skipWhitespace(b []byte) int {
	index := 0
	for index <= len(b)-1 {
		switch b[index] {
		case ' ', '\t', '\r', '\n':
			index++
		default:
			return index
		}
	}
	return index
}
