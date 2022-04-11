package interceptor

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

const MobileRegexp = `^(?:\+?86)?1(?:3\d{3}|5[^4\D]\d{2}|8\d{3}|7(?:[0-35-9]\d{2}|4(?:0\d|1[0-2]|9\d))|9[0-35-9]\d{2}|6[2567]\d{2}|4[579]\d{2})\d{6}$`

func isMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	match, err := regexp.MatchString(MobileRegexp, mobile)
	if err != nil {
		return false
	}

	return match
}

func isNumberSlice(fl validator.FieldLevel) bool {
	s := fl.Field()
	slice, ok := s.Interface().([]int64)
	for _, v := range slice {
		if v <= 0 {
			return false
		}
	}

	return ok
}

func isStringSlice(fl validator.FieldLevel) bool {
	s := fl.Field()
	strings, ok := s.Interface().([]string)
	for _, v := range strings {
		if len(v) >= 20 {
			return false
		}
	}

	return ok
}

func isChineseChar(fl validator.FieldLevel) bool {
	s := fl.Field()
	str, ok := s.Interface().(string)
	if !ok {
		return false
	}
	return func(str string) bool {
		for _, r := range str {
			if unicode.Is(unicode.Scripts["Han"], r) {
				return true
			}
		}
		return false
	}(str)
}
