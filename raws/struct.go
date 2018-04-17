package raws // import "github.com/BenLubar/dfide/raws"

import (
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type stringOrIndex struct {
	String string
	Index  int
}

func makeIndexString(s string) []stringOrIndex {
	parts := strings.Split(s, ".")
	converted := make([]stringOrIndex, len(parts))
	for i, p := range parts {
		if n, err := strconv.Atoi(p); err == nil {
			converted[i].Index = n
		} else {
			converted[i].String = p
		}
	}
	return converted
}

type structTag struct {
	Name  []stringOrIndex
	Char  bool
	Union bool
}

var structTagElements = func() map[string]reflect.StructField {
	m := make(map[string]reflect.StructField)

	t := reflect.TypeOf(structTag{})
	for i, n := 0, t.NumField(); i < n; i++ {
		f := t.Field(i)
		if f.Name == "Name" {
			continue
		}

		m[strings.ToLower(f.Name)] = f
	}

	return m
}()

func parseStructTag(t reflect.StructTag) *structTag {
	s, ok := t.Lookup("raws")
	if !ok {
		return nil
	}

	parts := strings.Split(s, ",")

	if len(parts) == 0 || parts[0] == "" {
		panic("raws: invalid struct tag")
	}

	st := &structTag{
		Name: makeIndexString(parts[0]),
	}

	v := reflect.ValueOf(st).Elem()

	for i := 1; i < len(parts); i++ {
		f, ok := structTagElements[parts[i]]
		if !ok {
			panic("raws: invalid struct tag: no such flag: " + parts[i])
		}

		vf := v.FieldByIndex(f.Index)
		if vf.Bool() {
			panic("raws: invalid struct tag: duplicate flag: " + parts[i])
		}
		vf.SetBool(true)
	}

	if st.Union && (len(st.Name) != 1 || st.Name[0].String == "") {
		panic("raws: invalid struct tag: union must have simple name")
	}

	return st
}

type typeDescription struct {
	union  bool
	fields []reflect.StructField
	tags   []structTag
}

var typeDescriptions = make(map[reflect.Type]*typeDescription)
var typeDescriptionLock sync.Mutex

func getTypeDescription(t reflect.Type) *typeDescription {
	typeDescriptionLock.Lock()
	defer typeDescriptionLock.Unlock()

	return getTypeDescriptionLocked(t)
}

func getTypeDescriptionLocked(t reflect.Type) *typeDescription {
	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	if d, ok := typeDescriptions[t]; ok {
		return d
	}

	d := &typeDescription{}
	for i, n := 0, t.NumField(); i < n; i++ {
		f := t.Field(i)
		st := parseStructTag(f.Tag)
		if st != nil {
			if st.Char && f.Type.Kind() != reflect.Int32 &&
				(f.Type.Kind() != reflect.Slice || f.Type.Elem().Kind() != reflect.Int32) {
				panic("raws: invalid struct tag: char flag can only be used on rune or []rune")
			}

			d.fields = append(d.fields, f)
			d.tags = append(d.tags, *st)
		}
	}

	if len(d.tags) == 0 {
		typeDescriptions[t] = nil
		return nil
	}

	if d.tags[0].Union {
		for _, t := range d.tags {
			if !t.Union {
				panic("raws: if a struct is a union, all fields of the struct must be tagged")
			}
		}
		d.union = true
	}

	typeDescriptions[t] = d
	return d
}
