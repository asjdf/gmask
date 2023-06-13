package gmask

var defaultMasker *Masker

func init() {
	defaultMasker = New().
		RegMaskAnyFunc(MaskTypeZero, MaskZero).
		RegMaskStringFunc(MaskTypeChar, MaskCharString()).
		RegMaskStringFunc(MaskTypeRandom, MaskRandString).
		RegMaskStringFunc(MaskTypeHash, MaskHashString)
}

func Mask[T any](target T) (ret T, err error) {
	v, err := defaultMasker.Mask(target)
	if err != nil {
		return ret, err
	}

	return v.(T), nil
}

func Float64(value float64, tag ...string) (float64, error) {
	return defaultMasker.Float64(value, tag...)
}

func String(value string, tag ...string) (string, error) {
	return defaultMasker.String(value, tag...)
}

func Int(value int, tag ...string) (int, error) {
	return defaultMasker.Int(value, tag...)
}

func Uint(value uint, tag ...string) (uint, error) {
	return defaultMasker.Uint(value, tag...)
}

func Any(value any, tag ...string) (hit bool, output any, err error) {
	return defaultMasker.Any(value, tag...)
}
