package raws // import "github.com/BenLubar/dfide/raws"

import (
	"encoding"
	"reflect"
	"strconv"
)

func (w *Writer) Serialize(v interface{}) error {
	return w.SerializeValue(reflect.ValueOf(v))
}

func (w *Writer) SerializeValue(v reflect.Value) error {
	return w.serializeValueOrUnion(w.object, v)
}

func (w *Writer) serializeValueOrUnion(tagName string, v reflect.Value) error {
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return nil
	}

	v = reflect.Indirect(v)

	if v.Kind() == reflect.Slice {
		for i, n := 0, v.Len(); i < n; i++ {
			if err := w.serializeValueOrUnion(tagName, v.Index(i)); err != nil {
				return err
			}
		}
		return nil
	}

	if v.Kind() == reflect.Bool {
		if !v.Bool() {
			return nil
		}
		return w.WriteTag([]string{tagName})
	}

	if desc := getTypeDescription(v.Type()); desc.union {
		found := false
		for i, f := range desc.fields {
			field := v.FieldByIndex(f.Index)
			if field.Kind() != reflect.Ptr || !field.IsNil() {
				if found {
					found = false
					break
				}
				found = true

				if err := w.serializeValue(desc.tags[i].Name[0].String, reflect.Indirect(field)); err != nil {
					return err
				}
			}
		}
		if !found {
			panic("raws: exactly one field in a union must be non-nil")
		}
		return nil
	}

	return w.serializeValue(tagName, v)
}

func (w *Writer) serializeValue(tagName string, v reflect.Value) error {
	desc := getTypeDescription(v.Type())

	tag := []string{tagName}
	for i, t := range desc.tags {
		if len(t.Name) == 1 && t.Name[0].String == "" {
			f := v.FieldByIndex(desc.fields[i].Index)
			for len(tag) <= t.Name[0].Index {
				tag = append(tag, "")
			}
			var err error
			tag[t.Name[0].Index], err = w.serializeField(f, t)
			if err != nil {
				return err
			}
		}
	}

	err := w.WriteTag(tag)
	if err != nil {
		return err
	}

	w.Indent++
	defer func() {
		w.Indent--
	}()

	seen := make(map[string]bool)
	for i, t := range desc.tags {
		if t.Name[0].String != "" && !seen[t.Name[0].String] {
			seen[t.Name[0].String] = true
			f := v.FieldByIndex(desc.fields[i].Index)
			if len(t.Name) == 1 {
				err = w.serializeValueOrUnion(t.Name[0].String, f)
			} else {
				err = w.serializeSubValue(t.Name[0].String, v, desc)
			}
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (w *Writer) serializeSubValue(tagName string, v reflect.Value, desc *typeDescription) error {
	tag := []string{tagName}

	var err error
	var any, wasSlice bool

	for i, t := range desc.tags {
		if len(t.Name) == 2 && t.Name[0].String == tagName && t.Name[1].String == "" {
			f := v.FieldByIndex(desc.fields[i].Index)
			if f.Kind() == reflect.Ptr && f.IsNil() {
				continue
			}
			if reflect.Indirect(f).Kind() == reflect.Slice && t.Name[1].Index == 1 {
				wasSlice = true
				if any {
					panic("raws: a single sub-value cannot contain slices and scalar values")
				}
				f = reflect.Indirect(f)
				tag := append(tag, "")
				for i, n := 0, f.Len(); i < n; i++ {
					tag[1], err = w.serializeField(f.Index(i), t)
					if err != nil {
						return err
					}

					if err := w.WriteTag(tag); err != nil {
						return err
					}
				}
				continue
			}
			if wasSlice {
				panic("raws: a single sub-value cannot contain slices and scalar values")
			}
			for len(tag) <= t.Name[1].Index {
				tag = append(tag, "")
			}
			tag[t.Name[1].Index], err = w.serializeField(f, t)
			if err != nil {
				return err
			}
			any = true
		}
	}

	if !any {
		return nil
	}

	return w.WriteTag(tag)
}

func (w *Writer) serializeField(v reflect.Value, t structTag) (string, error) {
	if m, ok := v.Addr().Interface().(encoding.TextMarshaler); ok {
		b, err := m.MarshalText()
		return string(b), err
	}

	switch v.Kind() {
	case reflect.String:
		return v.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if t.Char {
			return FromChar(rune(v.Int()))
		}
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), nil
	}
	panic("raws: unhandled type: " + v.Kind().String())
}

func (w *Writer) SerializeAll(v interface{}) error {
	return w.SerializeAllValue(reflect.ValueOf(v))
}

func (w *Writer) SerializeAllValue(v reflect.Value) error {
	v = reflect.Indirect(v)

	if v.Kind() != reflect.Slice {
		panic("raws: SerializeAll must be called on a slice")
	}

	for i, n := 0, v.Len(); i < n; i++ {
		if i != 0 {
			if err := w.w.WriteByte('\n'); err != nil {
				return err
			}
		}
		if err := w.SerializeValue(v.Index(i)); err != nil {
			return err
		}
	}

	return nil
}
