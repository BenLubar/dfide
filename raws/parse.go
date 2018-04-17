package raws // import "github.com/BenLubar/dfide/raws"

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
)

func (r *Reader) Parse(v interface{}) error {
	return r.ParseValue(reflect.ValueOf(v).Elem())
}

func ensure(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr {
		return v
	}

	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
	return v.Elem()
}

func (r *Reader) ParseValue(v reflect.Value) error {
	if err := r.init(); err != nil {
		return err
	}

	tag, err := r.nextTag()
	if err != nil {
		return err
	}

	v = ensure(v)

	if desc := getTypeDescription(v.Type()); desc.union {
		for i, t := range desc.tags {
			if t.Name[0].String == tag[0] {
				return r.parseValue(v.FieldByIndex(desc.fields[i].Index), tag)
			}
		}
		return fmt.Errorf("raws: no match for object type: %q", tag[0])
	}

	return r.parseValue(v, tag)
}

func (r *Reader) parseValue(v reflect.Value, startTag []string) error {
	v = ensure(v)
	desc := getTypeDescription(v.Type())

	for i, t := range desc.tags {
		if len(t.Name) == 1 && t.Name[0].String == "" {
			f := v.FieldByIndex(desc.fields[i].Index)
			if len(startTag) <= t.Name[0].Index {
				return fmt.Errorf("raws: %q tag too short", startTag[0])
			}
			if err := r.parseField(f, t, startTag[t.Name[0].Index]); err != nil {
				return err
			}
		}
	}

	for {
		tag, err := r.nextTag()
		if err != nil {
			return err
		}

		var found bool
		for i, t := range desc.tags {
			if t.Name[0].String == tag[0] {
				var slice reflect.Value
				f := v.FieldByIndex(desc.fields[i].Index)
				if f.Kind() == reflect.Slice {
					slice, f = f, reflect.New(f.Type().Elem()).Elem()
				}
				if len(t.Name) != 1 {
					if len(t.Name) != 2 || t.Name[1].String != "" {
						panic("raws: TODO: handle complex field paths")
					}
					if len(tag) <= t.Name[1].Index {
						return fmt.Errorf("raws: %q tag too short", tag[0])
					}
					err = r.parseField(f, t, tag[t.Name[1].Index])
				} else if f.Kind() == reflect.Bool && len(tag) == 1 {
					f.SetBool(true)
					err = nil
				} else {
					err = r.parseValue(f, tag)
				}
				if err != nil {
					return err
				}

				if slice.IsValid() {
					slice.Set(reflect.Append(slice, f))
				}

				found = true
				// don't break
			}
		}
		if !found {
			r.unreadTag = tag
			return nil
		}
	}
}

func (r *Reader) parseField(v reflect.Value, t structTag, value string) error {
	v = ensure(v)

	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if t.Char {
			c, err := ToChar(value)
			v.SetInt(int64(c))
			return err
		} else {
			n, err := strconv.ParseInt(value, 10, v.Type().Bits())
			v.SetInt(n)
			return err
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if t.Char {
			c, err := ToChar(value)
			v.SetUint(uint64(c))
			return err
		} else {
			n, err := strconv.ParseUint(value, 10, v.Type().Bits())
			v.SetUint(n)
			return err
		}
	}
	panic("raws: unhandled type: " + v.Kind().String())
}

func (r *Reader) ParseAll(v interface{}) error {
	return r.ParseAllValue(reflect.ValueOf(v).Elem())
}

func (r *Reader) ParseAllValue(v reflect.Value) error {
	if v.Kind() != reflect.Slice {
		panic("raws: ParseAll must be called on a pointer to a slice")
	}

	t := v.Type().Elem()
	for {
		e := reflect.New(t).Elem()

		if err := r.ParseValue(e); err == nil {
			v.Set(reflect.Append(v, e))
		} else if err == io.EOF {
			v.Set(reflect.Append(v, e))
			return nil
		} else {
			return err
		}
	}
}
