package fusion

import (
	"io"
	"reflect"
	"time"
)

func Must(args ...interface{}) {
	if err := args[len(args)-1]; err != nil {
		panic(err)
	}
}

func Check(args ...interface{}) []interface{} {
	if err := args[len(args)-1]; err != nil {
		panic(err)
	}
	return args[:len(args)-1]
}

func Ingore(args ...interface{}) []interface{} {
	return args[:len(args)-1]
}

func FromReader7BitEncodedInt(r io.Reader) (uint, error) {
	var value uint
	var count int
	for {
		var bytes [1]byte
		if _, err := io.ReadFull(r, bytes[:]); err != nil {
			return 0, err
		}
		value |= uint(bytes[0]&0x7f) << (count * 7)
		if bytes[0] > 0x7f {
			count += 1
		} else {
			break
		}
	}
	return value, nil
}

func ToWriter7BitEncodedInt(value uint, w io.Writer) error {
	var buffer [10]byte
	var count int
	for {
		var flag byte
		if value > 0x7f {
			flag = 0x80
		}
		buffer[count] = flag | byte(value)
		count += 1
		if flag > 0 {
			value >>= 7
		} else {
			break
		}
	}
	if _, err := w.Write(buffer[:count]); err != nil {
		return err
	}
	return nil
}

func NotifyNonBlock(c chan bool) {
	if len(c) == 0 {
		select {
		case c <- true:
		default:
		}
	}
}

func WaitBlock(c chan bool) bool {
	select {
	case _, ok := <-c:
		return ok
	}
}

func WaitTimeout(c chan bool, d time.Duration) bool {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-c:
		return false
	case <-timer.C:
		return true
	}
}

func Contain(slice interface{}, value interface{}) bool {
	sliceVal := reflect.ValueOf(slice)
	for i, n := 0, sliceVal.Len(); i < n; i++ {
		if sliceVal.Index(i).Interface() == value {
			return true
		}
	}
	return false
}

func RemoveIf(slice interface{}, test func(interface{}) bool) interface{} {
	sliceVal, numVal := reflect.ValueOf(slice), 0
	for i, n := 0, sliceVal.Len(); i < n; i++ {
		if !test(sliceVal.Index(i).Interface()) {
			if i != numVal {
				sliceVal.Index(numVal).Set(sliceVal.Index(i))
			}
			numVal++
		}
	}
	if numVal != sliceVal.Len() {
		return sliceVal.Slice(0, numVal).Interface()
	}
	return slice
}
