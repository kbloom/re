package re_test

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/ghemawat/re"
)

func TestFind(t *testing.T) {
	type mytype int

	type testcase struct {
		re       string
		input    string
		result   bool
		args     []interface{}
		expected []interface{}
	}

	// test returns a testcase that scans for the specified re in
	// input.  The even elements of argexpect are pointers that
	// will be passed as output arguments to re.Scan.  The odd
	// elements of argexpect are the expected values that re.Scan
	// should store into the preceding pointers.
	test := func(re, input string, result bool, argexpect ...interface{}) testcase {
		t := testcase{re, input, result, nil, nil}
		for i := 0; i < len(argexpect); i += 2 {
			t.args = append(t.args, argexpect[i])
			t.expected = append(t.expected, argexpect[i+1])
		}
		return t
	}

	for _, c := range []testcase{
		// Tests without any argument extraction.
		test(`(\w+):(\d+)`, "", false),
		test(`(\w+):(\d+)`, "host:1234x", true),
		test(`(\w+):(\d+)`, "-host:1234-", true),
		test(`(\w+):(\d+)`, "host:x1234", false),
		test(`^(\w+):(\d+)$`, "host:1234", true, nil, nil),
		test(`^(\w+):(\d+)$`, "host:1234x", false, nil, nil),

		// not enough matches
		test(`^\w+:\d+$`, "host:1234", false, new(string), nil),

		// extraction into nil
		test(`^(\w+):(\d+)$`, "host:1234", true, nil, nil, nil, nil),

		// missing sub-expression
		test(`^(\w+):((\d+))?`, "host:", true, nil, nil, nil, nil, nil, nil),
		test(`^(\w+):((\d+))?`, "host:", false, nil, nil, new(int), nil),

		// combination of multiple arguments
		test(`(\w+):(\d+)`, "h:80", true, new(string), "h", new(int), 80),

		// unsupported type
		test(`(.*)`, "1234", false, new(mytype), nil),

		// string
		test(`(.*):\d+`, "host:1234", true, new(string), "host"),

		// []byte
		test(`(.*):\d+`, "host:1234", true, new([]byte), []byte("host")),
		test(`(.*):\d+`, ":1234", true, new([]byte), []byte("")),

		// int
		test(`(\d+)`, "1234", true, new(int), 1234),
		test(`(.*)`, "-1234", true, new(int), -1234),
		test(`(.*)`, "123456789123456789123456789", false, new(int), nil),
		test(`(.*)`, "-123456789123456789123456789", false, new(int), nil),
		test(`(.*)`, "0x10", true, new(int), 0x10),
		test(`(.*)`, "010", true, new(int), 010),

		// uint
		test(`(\d+)`, "1234", true, new(uint), uint(1234)),
		test(`(\d+)`, "123456789123456789123456789", false, new(uint), nil),

		// uintptr
		test(`(\d+)`, "1234", true, new(uintptr), uintptr(1234)),
		test(`(\d+)`, "123456789123456789123456789", false, new(uintptr), nil),

		// uint8
		test(`(.*)`, "0", true, new(uint8), uint8(0)),
		test(`(.*)`, "17", true, new(uint8), uint8(17)),
		test(`(.*)`, "255", true, new(uint8), uint8(255)),
		test(`(.*)`, "256", false, new(uint8), nil),
		test(`(.*)`, "x", false, new(uint8), nil),

		// uint16
		test(`(.*)`, "0", true, new(uint16), uint16(0)),
		test(`(.*)`, "17", true, new(uint16), uint16(17)),
		test(`(.*)`, "65535", true, new(uint16), uint16(65535)),
		test(`(.*)`, "65536", false, new(uint16), nil),
		test(`(.*)`, "x", false, new(uint16), nil),

		// uint32
		test(`(.*)`, "0", true, new(uint32), uint32(0)),
		test(`(.*)`, "17", true, new(uint32), uint32(17)),
		test(`(.*)`, "4294967295", true, new(uint32), uint32(4294967295)),
		test(`(.*)`, "4294967296", false, new(uint32), nil),
		test(`(.*)`, "x", false, new(uint32), nil),

		// uint64
		test(`(.*)`, "0", true, new(uint64), uint64(0)),
		test(`(.*)`, "17", true, new(uint64), uint64(17)),
		test(`(.*)`, "18446744073709551615", true, new(uint64), uint64(18446744073709551615)),
		test(`(.*)`, "18446744073709551616", false, new(uint64), nil),
		test(`(.*)`, "x", false, new(uint64), nil),

		// int8
		test(`(.*)`, "0", true, new(int8), int8(0)),
		test(`(.*)`, "17", true, new(int8), int8(17)),
		test(`(.*)`, "127", true, new(int8), int8(127)),
		test(`(.*)`, "128", false, new(int8), nil),
		test(`(.*)`, "x", false, new(int8), nil),

		// int16
		test(`(.*)`, "0", true, new(int16), int16(0)),
		test(`(.*)`, "17", true, new(int16), int16(17)),
		test(`(.*)`, "32767", true, new(int16), int16(32767)),
		test(`(.*)`, "32768", false, new(int16), nil),
		test(`(.*)`, "x", false, new(int16), nil),

		// int32
		test(`(.*)`, "0", true, new(int32), int32(0)),
		test(`(.*)`, "17", true, new(int32), int32(17)),
		test(`(.*)`, "2147483647", true, new(int32), int32(2147483647)),
		test(`(.*)`, "2147483648", false, new(int32), nil),
		test(`(.*)`, "x", false, new(int32), nil),

		// int64
		test(`(.*)`, "0", true, new(int64), int64(0)),
		test(`(.*)`, "17", true, new(int64), int64(17)),
		test(`(.*)`, "9223372036854775807", true, new(int64), int64(9223372036854775807)),
		test(`(.*)`, "9223372036854775808", false, new(int64), nil),
		test(`(.*)`, "x", false, new(int64), nil),

		// float32
		test(`(.*)`, "0", true, new(float32), float32(0)),
		test(`(.*)`, "1.25e2", true, new(float32), float32(1.25e2)),
		test(`(.*)`, "1e40", false, new(float32), nil),
		test(`(.*)`, "x", false, new(float32), nil),

		// float64
		test(`(.*)`, "0", true, new(float64), float64(0)),
		test(`(.*)`, "1.25e2", true, new(float64), float64(1.25e2)),
		test(`(.*)`, "1e40", true, new(float64), float64(1e40)),
		test(`(.*)`, "1e400", false, new(float64), nil),
		test(`(.*)`, "x", false, new(float64), nil),
	} {
		err := re.Scan(regexp.MustCompile(c.re), []byte(c.input), c.args...)
		if !c.result {
			if err == nil {
				t.Errorf("Find(`%s`, `%s`, ...) succeeded unexpectedly", c.re, c.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("unexpected error: %s", err)
			continue
		}
		for i, a := range c.args {
			if a == nil && c.expected[i] == nil {
				continue
			}
			// c.args[i] wraps a *T and c.expected[i] wraps a T.
			// Dereference c.args[i] to get a T we can compare.
			av := reflect.Indirect(reflect.ValueOf(a)).Interface()
			if !reflect.DeepEqual(av, c.expected[i]) {
				t.Errorf("Find(`%s`, `%s`, ...): result[%d] is `%v`; expected `%v`\n",
					c.re, c.input, i, av, c.expected[i])
			}

		}
	}
}

func TestReFunc(t *testing.T) {
	var arg string
	savearg := func(a []byte) error {
		arg = string(a)
		return nil
	}
	hp := `^(\w+):(\d+)$`
	str := "host:1234"
	if err := re.Scan(regexp.MustCompile(hp), []byte(str), savearg); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if arg != "host" {
		t.Fatalf("Find(`%s`, `%s`, savearg): did not call function", hp, str)
	}

	fail := func(a []byte) error {
		arg = string(a)
		return fmt.Errorf("error")
	}
	if err := re.Scan(regexp.MustCompile(hp), []byte(str), fail); err == nil {
		t.Fatalf("Find(`%s`, `%s`, fail): succeeded unexpectedly", hp, str)
	}
}

func TestOptional(t *testing.T) {
	hp := regexp.MustCompile(`^(\w+)(?::(\d+))?$`)

	cases := []struct {
		haystack string
		host     string
		port     int
		hasPort  bool
	}{
		{
			haystack: "host:1234",
			host:     "host",
			port:     1234,
			hasPort:  true,
		},
		{
			haystack: "host",
			host:     "host",
			hasPort:  false,
		},
	}

	for _, c := range cases {
		var (
			host    string
			port    int
			hasPort bool
		)

		if err := re.Scan(hp, []byte(c.haystack), &host, re.Optional(&port, &hasPort)); err != nil {
			t.Errorf("re.Scan(%s, %s, ...) unexpected error: %s", hp, c.haystack, err)
			continue
		}
		if host != c.host {
			t.Errorf("re.Scan(%s, %s, ...): host =  %s, want %s", hp, c.haystack, host, c.host)
		}
		if port != c.port {
			t.Errorf("re.Scan(%s, %s, ...): port =  %d, want %d", hp, c.haystack, port, c.port)
		}
		if hasPort != c.hasPort {
			t.Errorf("re.Scan(%s, %s, ...): hasPort =  %T, want %T", hp, c.haystack, hasPort, c.hasPort)
		}
	}
}

func TestRePosition(t *testing.T) {
	hp := `(\w+):(\d+)`
	bytes := []byte("host:1234 host2:2345")
	var pos re.Position
	var host string
	var port int
	if err := re.Scan(regexp.MustCompile(hp), bytes, &pos, &host, &port); err != nil {
		t.Fatalf("First match: unexpected error: %s", err)
	} else {
		if pos.Start != 0 {
			t.Errorf("First match: pos.Start = %d, want 0", pos.Start)
		}
		if pos.End != 9 {
			t.Errorf("First match: pos.End = %d, want 9", pos.End)
		}
		if host != "host" {
			t.Errorf("First match: host = %s, want \"host\"", host)
		}
		if port != 1234 {
			t.Errorf("First match: port = %d, want 1234", port)
		}
	}

	bytes = bytes[pos.End:]

	if err := re.Scan(regexp.MustCompile(hp), bytes, &pos, &host, &port); err != nil {
		t.Fatalf("Second match: unexpected error: %s", err)
	} else {
		// Offsets in the second Scan call should count from the beginning of the tail,
		// not from the beginning of the original string.
		if pos.Start != 1 {
			t.Errorf("Second match: pos.Start = %d, want 1", pos.Start)
		}
		if pos.End != 11 {
			t.Errorf("Second match: pos.End = %d, want 11", pos.End)
		}
		if host != "host2" {
			t.Errorf("Second match: host = %s, want \"host2\"", host)
		}
		if port != 2345 {
			t.Errorf("Second match: port = %d, want 2345", port)
		}
	}

}

func TestReAliasing(t *testing.T) {
	b := []byte("hello")
	var m []byte
	if err := re.Scan(regexp.MustCompile(`(.*)`), b, &m); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if string(m) != "hello" {
		t.Fatalf("Find extracted wrong value")
	}
	b[0] = 'j'
	if string(m) != "jello" {
		t.Fatalf("extracted byte slice does not alias input")
	}
}
