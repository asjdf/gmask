package gmask

import (
	"reflect"
	"strings"
)

const tagName = "mask"

type (
	MaskFloat64Func func(value float64, arg ...string) (float64, error)
	MaskStringFunc  func(value string, arg ...string) (string, error)
	MaskIntFunc     func(value int, arg ...string) (int, error)
	MaskUintFunc    func(value uint, arg ...string) (uint, error)
	MaskAnyFunc     func(value any, arg ...string) (any, error)
)

type Masker struct {
	// [MaskName] MaskFunction
	maskFloat64FuncMap map[string]MaskFloat64Func
	maskStringFuncMap  map[string]MaskStringFunc
	maskIntFuncMap     map[string]MaskIntFunc
	maskUintFuncMap    map[string]MaskUintFunc
	maskAnyFuncMap     map[string]MaskAnyFunc
}

func New() *Masker {
	return &Masker{
		maskFloat64FuncMap: make(map[string]MaskFloat64Func),
		maskStringFuncMap:  make(map[string]MaskStringFunc),
		maskIntFuncMap:     make(map[string]MaskIntFunc),
		maskUintFuncMap:    make(map[string]MaskUintFunc),
		maskAnyFuncMap:     make(map[string]MaskAnyFunc),
	}
}

func (m *Masker) RegMaskStringFunc(maskName string, mask MaskStringFunc) *Masker {
	m.maskStringFuncMap[maskName] = mask
	return m
}

func (m *Masker) RegMaskUintFunc(maskName string, mask MaskUintFunc) *Masker {
	m.maskUintFuncMap[maskName] = mask
	return m
}

func (m *Masker) RegMaskIntFunc(maskName string, mask MaskIntFunc) *Masker {
	m.maskIntFuncMap[maskName] = mask
	return m
}

func (m *Masker) RegMaskFloat64Func(maskName string, mask MaskFloat64Func) *Masker {
	m.maskFloat64FuncMap[maskName] = mask
	return m
}

func (m *Masker) RegMaskAnyFunc(maskName string, mask MaskAnyFunc) *Masker {
	m.maskAnyFuncMap[maskName] = mask
	return m
}

func (m *Masker) Mask(target any) (ret any, err error) {
	rv, err := m.mask(reflect.ValueOf(target), reflect.Value{})
	if err != nil {
		return ret, err
	}

	return rv.Interface(), nil
}

func (m *Masker) Float64(value float64, tag ...string) (float64, error) {
	if len(tag) == 0 || len(tag[0]) == 0 {
		return value, nil
	}

	if maskFunc, exist := m.maskFloat64FuncMap[tag[0]]; exist {
		return maskFunc(value, tag[1:]...)
	}

	if ok, v, err := m.Any(value, tag...); ok {
		return v.(float64), err
	}

	return value, nil
}

func (m *Masker) String(value string, tag ...string) (string, error) {
	if len(tag) == 0 || len(tag[0]) == 0 {
		return value, nil
	}

	if maskFunc, exist := m.maskStringFuncMap[tag[0]]; exist {
		return maskFunc(value, tag[1:]...)
	}

	if ok, v, err := m.Any(value, tag...); ok {
		return v.(string), err
	}

	return value, nil
}

func (m *Masker) Int(value int, tag ...string) (int, error) {
	if len(tag) == 0 || len(tag[0]) == 0 {
		return value, nil
	}

	if maskFunc, exist := m.maskIntFuncMap[tag[0]]; exist {
		return maskFunc(value, tag[1:]...)
	}

	if ok, v, err := m.Any(value, tag...); ok {
		return v.(int), err
	}

	return value, nil
}

func (m *Masker) Uint(value uint, tag ...string) (uint, error) {
	if len(tag) == 0 || len(tag[0]) == 0 {
		return value, nil
	}

	if maskFunc, exist := m.maskUintFuncMap[tag[0]]; exist {
		return maskFunc(value, tag[1:]...)
	}

	if ok, v, err := m.Any(value, tag...); ok {
		return v.(uint), err
	}

	return value, nil
}

func (m *Masker) Any(value any, tag ...string) (hit bool, output any, err error) {
	if len(tag) == 0 || len(tag[0]) == 0 {
		return false, value, nil
	}

	if maskFunc, exist := m.maskAnyFuncMap[tag[0]]; exist {
		output, err = maskFunc(value, tag[1:]...)
		return true, output, err
	}

	return false, value, nil
}

func (m *Masker) mask(rv reflect.Value, mp reflect.Value, tag ...string) (reflect.Value, error) {
	//if ok, v, err := m.maskAnyValue(tag, rv); ok {
	//	return v, err
	//}
	switch rv.Type().Kind() {
	case reflect.Interface:
		return m.maskInterface(rv, mp, tag...)
	case reflect.Ptr:
		return m.maskPtr(rv, mp, tag...)
	case reflect.Struct:
		return m.maskStruct(rv, mp)
	case reflect.Array:
		return m.maskSlice(rv, mp, tag...)
	case reflect.Slice:
		if rv.IsNil() {
			return reflect.Zero(rv.Type()), nil
		}
		return m.maskSlice(rv, mp, tag...)
	//case reflect.Map:
	//	return m.maskMap(rv, mp, tag)
	case reflect.Float32, reflect.Float64:
		return m.maskFloat(rv, mp, tag...)
	case reflect.String:
		return m.maskString(rv, mp, tag...)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return m.maskInt(rv, mp, tag...)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return m.maskUint(rv, mp, tag...)
	default:
		if mp.IsValid() {
			mp.Set(rv)
			return mp, nil
		}
		return rv, nil
	}
}

func (m *Masker) maskInterface(rv reflect.Value, _ reflect.Value, tag ...string) (reflect.Value, error) {
	if rv.IsNil() {
		return reflect.Zero(rv.Type()), nil
	}

	mp := reflect.New(rv.Type()).Elem()
	rv2, err := m.mask(reflect.ValueOf(rv.Interface()), reflect.Value{}, tag...)
	if err != nil {
		return reflect.Value{}, err
	}
	mp.Set(rv2)

	return mp, nil
}

func (m *Masker) maskPtr(rv reflect.Value, _ reflect.Value, tag ...string) (reflect.Value, error) {
	if rv.IsNil() {
		return reflect.Zero(rv.Type()), nil
	}

	mp := reflect.New(rv.Type().Elem())
	rv2, err := m.mask(rv.Elem(), mp.Elem(), tag...)
	if err != nil {
		return reflect.Value{}, err
	}
	mp.Elem().Set(rv2)

	return mp, nil
}

func (m *Masker) maskStruct(rv reflect.Value, mp reflect.Value) (reflect.Value, error) {
	if rv.IsZero() {
		return reflect.Zero(rv.Type()), nil
	}

	rt := rv.Type()
	if !mp.IsValid() {
		mp = reflect.New(rt).Elem()
	}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		// skip private field
		if field.PkgPath != "" {
			continue
		}

		args := strings.Split(field.Tag.Get(tagName), ",")
		switch field.Type.Kind() {
		case reflect.String:
			s, err := m.String(rv.Field(i).String(), args...)
			if err != nil {
				return reflect.Value{}, err
			}
			mp.Field(i).SetString(s)
		default:
			rvf, err := m.mask(rv.Field(i), mp.Field(i), args...)
			if err != nil {
				return reflect.Value{}, err
			}
			mp.Field(i).Set(rvf)
		}
	}

	return mp, nil
}

func (m *Masker) maskSlice(rv reflect.Value, mp reflect.Value, tag ...string) (reflect.Value, error) {
	var rv2 reflect.Value

	if rt := rv.Type(); rt.Kind() == reflect.Array {
		rv2 = reflect.New(rt).Elem()
	} else {
		rv2 = reflect.MakeSlice(rv.Type(), rv.Len(), rv.Len())
	}
	for i := 0; i < rv.Len(); i++ {
		value := rv.Index(i)
		switch rv.Type().Elem().Kind() {
		case reflect.String:
			rvf, err := m.String(value.String(), tag...)
			if err != nil {
				return reflect.Value{}, err
			}
			rv2.Index(i).SetString(rvf)
		case reflect.Int:
			rvf, err := m.Int(int(value.Int()), tag...)
			if err != nil {
				return reflect.Value{}, err
			}
			rv2.Index(i).SetInt(int64(rvf))
		case reflect.Float64:
			rvf, err := m.Float64(value.Float(), tag...)
			if err != nil {
				return reflect.Value{}, err
			}
			rv2.Index(i).SetFloat(rvf)
		case reflect.Uint:
			rvf, err := m.Uint(uint(value.Uint()), tag...)
			if err != nil {
				return reflect.Value{}, err
			}
			rv2.Index(i).SetUint(uint64(rvf))
		default:
			rvf, err := m.mask(value, rv2.Index(i), tag...)
			if err != nil {
				return reflect.Value{}, err
			}
			rv2.Index(i).Set(rvf)
		}
	}

	if mp.IsValid() {
		mp.Set(rv2)
		return mp, nil
	}

	return rv2, nil
}

func (m *Masker) maskFloat(rv reflect.Value, mp reflect.Value, tag ...string) (reflect.Value, error) {
	if len(tag) == 0 {
		if mp.IsValid() {
			mp.Set(rv)
			return mp, nil
		}
		return rv, nil
	}

	fp, err := m.Float64(rv.Float(), tag...)
	if err != nil {
		return reflect.Value{}, err
	}
	if mp.IsValid() {
		mp.SetFloat(fp)
		return mp, nil
	}

	if rv.Type().Kind() != reflect.Float64 {
		return reflect.ValueOf(&fp).Elem().Convert(rv.Type()), nil
	}

	return reflect.ValueOf(&fp).Elem(), nil
}

func (m *Masker) maskString(rv, mp reflect.Value, tag ...string) (reflect.Value, error) {
	if len(tag) == 0 {
		if mp.IsValid() {
			mp.Set(rv)
			return mp, nil
		}
		return rv, nil
	}

	sp, err := m.String(rv.String(), tag...)
	if err != nil {
		return reflect.Value{}, err
	}
	if mp.IsValid() {
		mp.SetString(sp)
		return mp, nil
	}

	return valueOfString(sp), nil
}

func valueOfString(s string) reflect.Value {
	return reflect.ValueOf(&s).Elem()
}

func (m *Masker) maskInt(rv reflect.Value, mp reflect.Value, tag ...string) (reflect.Value, error) {
	if len(tag) == 0 {
		if mp.IsValid() {
			mp.Set(rv)
			return mp, nil
		}
		return rv, nil
	}

	ip, err := m.Int(int(rv.Int()), tag...)
	if err != nil {
		return reflect.Value{}, err
	}
	if mp.IsValid() {
		mp.SetInt(int64(ip))
		return mp, nil
	}

	if rv.Type().Kind() != reflect.Int {
		return reflect.ValueOf(&ip).Elem().Convert(rv.Type()), nil
	}

	return reflect.ValueOf(&ip).Elem(), nil
}

func (m *Masker) maskUint(rv reflect.Value, mp reflect.Value, tag ...string) (reflect.Value, error) {
	if len(tag) == 0 {
		if mp.IsValid() {
			mp.Set(rv)
			return mp, nil
		}
		return rv, nil
	}

	ip, err := m.Uint(uint(rv.Uint()), tag...)
	if err != nil {
		return reflect.Value{}, err
	}
	if mp.IsValid() {
		mp.SetUint(uint64(ip))
		return mp, nil
	}

	if rv.Type().Kind() != reflect.Uint {
		return reflect.ValueOf(&ip).Elem().Convert(rv.Type()), nil
	}

	return reflect.ValueOf(&ip).Elem(), nil
}
