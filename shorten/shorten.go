/*
Package shorten provides functions for converting unsigned integers
to short strings using baseN alphabets.

Example use:
	import (
		"fmt"
		"github.com/TShadwell/go-shorten/shorten"
	)

	const VALUE = 14794393443
	const DICT = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!\"£$%^&*()abcdefghijklmnopqrstuvwxyz[]{};:@'~#?/<>,.\\|`¬"

	func main(){
		myDict := shorten.MakeDictionary(DICT)

		outString, _ := myDict.Shorten(MYINT)

		convertedBack := myDict.Lengthen(outString)

		fmt.Println("Original number:", MYINT, ".")
		fmt.Println("Shortens to:", outString, ".")
		fmt.Println("Lengthens to:", convertedBack, ".")
	}

Outputs:

	Original number: 14794393443.

	Shortens to: 2dD)BC.

	Lengthens to: 14794393443.

This package also includes a Rearrange method, which is not designed to be done quickly. 
It is for making dictionaries that have particular phrases that eagle-eyed users can notice. Yes, an easter-egg function.
	import (
		"fmt"
		"github.com/TShadwell/go-shorten/shorten"
	)

	const PHRASE = "Lorem ipsum dolor sit amet" 
	const DICT = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!\"£$%^&*()abcdefghijklmnopqrstuvwxyz,."
	func main(){
		myDict := shorten.MakeDictionary(DICT)

		newDict, _ := myDict.Rearrange(
			PHRASE,
			shorten.Convertcase,
			shorten.ConvertSpecial,
		)

		fmt.Println(newDict)
	}

Outputs:
	Lorem_ipsuM-dOl0R.SIt,aETABCDFGHJKNPQUVWXYZ123456789!"£$%^&*()bcfghjknqvwxyz
*/
package shorten

import (
	"errors"
	"math"
	"unicode"
	"strings"
)

type Dictionary struct {
	dict    string
	dictLen uint
}

//MakeDictionary constructs a Dictionary object that
//is used for base conversion.
func MakeDictionary(dict string) *Dictionary {
	x := new(Dictionary)

	x.dict = dict
	x.dictLen = uint(len(dict))
	return x
}

//Shorten converts val to the base given by the length of its
//dictionary.
func (d Dictionary) Shorten(val uint) (o string, err error) {
	if val == 0 {
		return string(d.dict[0]), nil
	}
	for {
		remainder := val % d.dictLen
		val = uint(val / d.dictLen)
		o += string(d.dict[remainder])

		if val == 0 {
			return o, nil
		}
	}
	return "", errors.New("Cannot shorten.")
}

func (d Dictionary) valueOf(digit rune, position uint) uint {
	for i, chr := range d.dict {
		if chr == digit {
			return uint(i) * uint(math.Pow(float64(d.dictLen), float64(position)))
		}
	}
	panic("Improper value passed to ValueOf; rune '" + string(digit) + "' should be present in '" + d.dict + "'.")
}

//Converts v to an unsigned integer using the dictionary.
func (d Dictionary) Lengthen(v string) (o uint) {
	o = 0
	for i, chr := range v {
		o += d.valueOf(chr, uint(i))
	}
	return
}

//Returns string value of d.
func (d Dictionary) String()string{
	return d.dict
}

type rearrangeError uint8

const (
	OK rearrangeError = iota

	//Causes Rearrange to pass onto the next RearrangeMethod
	NOT_PRESENT

	//There are still possible values, c++ to get the next.
	ANOTHER_PRESENT

	//Causes Rearrange to fail with this error.
	FAIL
)

var rearrMap = map[rearrangeError]string{
	OK: "OK",
	NOT_PRESENT:"NOT_PRESENT",
	ANOTHER_PRESENT:"ANOTHER_PRESENT",
	FAIL:           "FAIL",
}

func (r rearrangeError) Error() string {
	return "Rearrange error of type " + rearrMap[r] + "."
}
func (r rearrangeError) String() string {
	return "const (" + rearrMap[r] + ");"
}

/*
A Rearranger converts from one rune to another in its domain, returning a rearrangeError
if something goes wrong or if the letter is not in its domain.

It may also take a count as a uint which sends any alternate conversions back.
*/
type Rearranger func(rune, uint8) (rune, rearrangeError)

//Fail stops rearranging immediately and returns the original dictionary and FAIL.
//
//Convertcase inverts letter case.
//
//Convertspecial maps e to 3, Σ, £ etc.
var (
	Fail           Rearranger = fail
	Convertcase    Rearranger = convertcase
	ConvertSpecial Rearranger = convertSpecial
)

func fail(n rune, count uint8) (rune, rearrangeError) {
	return 0, FAIL
}
func convertcase(n rune, c uint8) (out rune, err rearrangeError) {
	if !unicode.IsLetter(n) {
		err = NOT_PRESENT
		return
	}
	if unicode.IsUpper(n) {
		out = unicode.ToLower(n)
	} else if unicode.IsLower(n) {
		out = unicode.ToUpper(n)
	} else {
		err = NOT_PRESENT
		return
	}

	if out == n {
		err = NOT_PRESENT
	}

	return
}

var specialMap = map[rune][]rune{
	'a': {'4', '@'},
	'b': {'8'},
	'c': {'('},

	'e': {'Σ', '£', '3'},

	'g': {'9'},

	'i': {'1', '!', '|'},

	'l': {'1', '|'},

	'o': {'0', 'ϴ'},

	's': {'5'},
	't': {'7'},

	' ':{'_','-','.',',','=','+'},
}

func convertSpecial(n rune, c uint8) (out rune, err rearrangeError) {
	if specialMap[n]== nil {
		err = NOT_PRESENT
		return
	}
	if int(c)!=len(specialMap[n])-1{
		err = ANOTHER_PRESENT
	}
	out = specialMap[n][c]

	return
}

//Rearrange rearranges the characters in a dictionary such that the left hand side spells
//the text in s as best as possible.
//
//If a rune in s is not present in the Dictionary or a rune is duplicated in s,
//error will not be nil, returning a rearrangeError of NOT_PRESENT or DUPLICATE_RUNE respectively.
//
//If a rune is not present, Rearrange will attempt to compensate, given Rearrangers.
//Rearrangers are favoured in the order they are given, when each fails to convert,
//the function moves onto the next.
func (d *Dictionary) Rearrange(s string, rearrangers ...Rearranger) (*Dictionary, error) {
	var (
		out       = d.dict
		prefix    string
		canBeUsed = func(r rune) bool {
			return strings.ContainsRune(out, r)
		}
	)
	for _, v := range s {
		fulfilled := false

		var gotWorkingLetter = func(r rune) {
			prefix += string(r)
			out =  strings.Replace(out, string(r), "", 1)
			fulfilled = true
		}
		if canBeUsed(v){
			gotWorkingLetter(v)
			continue
		}
		for _, r := range rearrangers {

			state := ANOTHER_PRESENT

			for count := 0; state == ANOTHER_PRESENT; count++ {
				var response rune
				response, state = r(v, uint8(count))
				if (state == ANOTHER_PRESENT || state == OK) && canBeUsed(response) {
					gotWorkingLetter(response)
					break
				} else if state == FAIL {
					return d, FAIL
				}
			}

			if fulfilled {
				break
			}
		}

	}

	return MakeDictionary(prefix + out), OK
}
