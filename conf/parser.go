package conf

import (
	"bufio"
	"io"
	"reflect"
	"strings"

	"gopkg.in/logex.v1"
)

var (
	kindMap = map[string]*KindParser{}
)

func init() {
	var mapSlice map[string][]string
	var mapString map[string]string
	var slice []string
	kindMap[reflect.TypeOf(mapSlice).String()] = &KindParser{
		ParseMapSlice,
		func() interface{} {
			return map[string][]string{}
		},
	}
	kindMap[reflect.TypeOf(mapString).String()] = &KindParser{
		ParseMapString,
		func() interface{} {
			return map[string]string{}
		},
	}
	kindMap[reflect.TypeOf(slice).String()] = &KindParser{
		ParseSlice,
		func() interface{} {
			return []string{}
		},
	}
}

func Parse(ret interface{}, reader io.Reader) error {
	val := reflect.ValueOf(ret)
	if val.Type().Kind() == reflect.Ptr {
		val = val.Elem()
	}

	r := bufio.NewScanner(reader)
	for r.Scan() {
		line := r.Text()
		title, ok := IsGroupTitle(line)
		if !ok {
			continue
		}

		v := reflect.TypeOf(ret)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		k, ok := v.FieldByName(title)
		if !ok {
			return logex.NewTraceError("not such field name '", title, "'")
		}

		f := kindMap[k.Type.String()]
		if f == nil {
			return logex.NewTraceError("unknown type to parse '", k.Type.String(), "'")
		}

		tmpRet := f.Value()

		err := f.Func(&tmpRet, r)
		if err != nil {
			return logex.Trace(err)
		}

		vSet := reflect.ValueOf(tmpRet)
		val.FieldByName(title).Set(vSet)
	}
	return nil
}

type KindParser struct {
	Func  func(interface{}, *bufio.Scanner) error
	Value func() interface{}
}

func IsGroupTitle(l string) (string, bool) {
	if len(l) <= 2 {
		return "", false
	}
	if l[0] == '[' && l[len(l)-1] == ']' {
		return l[1 : len(l)-1], true
	}
	return "", false
}

func Set(to, from interface{}) {
	v := reflect.ValueOf(from)
	reflect.ValueOf(to).Elem().Set(v)
}

func ParseMapString(ret interface{}, s *bufio.Scanner) error {
	data := make(map[string]string)
	for s.Scan() {
		line := s.Text()
		if len(line) == 0 {
			break
		}

		sp := strings.Split(line, ":")
		if len(sp) != 2 {
			return logex.NewTraceError(`parse error in "`, line, `"`)
		}
		data[strings.TrimSpace(sp[0])] = strings.TrimSpace(sp[1])
	}

	Set(ret, data)
	return nil
}

func ParseSlice(ret interface{}, s *bufio.Scanner) error {
	data := []string{}
	for s.Scan() {
		line := s.Text()
		if len(line) == 0 {
			break
		}

		data = append(data, strings.TrimSpace(line))
	}

	Set(ret, data)
	return nil
}

func ParseMapSlice(ret interface{}, s *bufio.Scanner) error {
	data := make(map[string][]string)
	key := ""
	var item []string
	for s.Scan() {
		line := s.Text()
		if len(line) == 0 {
			break
		}

		if strings.HasSuffix(line, ":") {
			key = line[:len(line)-1]
			item = []string{}
			data[key] = item
			continue
		}

		if key != "" {
			item = append(item, strings.TrimSpace(line))
			data[key] = item
		}
	}

	Set(ret, data)
	return nil
}
