package shorten

import (
	"math"
	"errors"
)

type Dictionary struct{
	dict string
	dictLen uint
}

func MakeDictionary(dict string) *Dictionary{
	x := new(Dictionary)

	x.dict = dict
	x.dictLen = uint(len(dict))
	return x
}

func (d Dictionary) Shorten(val uint) (o string, err error){
	if val == 0{
		return string(d.dict[0]), nil
	}
	for {
		remainder := val %d.dictLen
		val = uint(val/d.dictLen)
		o += string(d.dict[remainder])

		if val == 0{
			return o, nil
		}
	}
	return "", errors.New("Cannot shorten.")
}

func (d Dictionary) valueOf(digit rune, position uint) uint{
	for i, chr := range d.dict{
		if chr == digit{
			return uint(i) * uint(math.Pow(float64(d.dictLen), float64(position)))
		}
	}
	panic("Improper value passed to ValueOf; rune '" + string(digit) + "' should be present in '" + d.dict + "'.")
}


func (d Dictionary) Lengthen(v string) (o uint){
	o=0
	for i, chr := range v{
		o += d.valueOf(chr, uint(i))
	}
	return
}
