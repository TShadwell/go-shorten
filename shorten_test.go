package shorten


import "testing"

func TestShorten(t *testing.T){
	const dict ="ABCDEFHIJLKMNOPQRSTUVWXYZ1234567890abcefghijklmnopqrstuvwxyz"
	{
			tester := MakeDictionary(dict)
			for i:=0;i<500;i++{
				input := uint(i)
				output, _ := tester.Shorten(input)

				reversed := tester.Lengthen(output)
				if !(reversed == input){
					t.Errorf("Lengthen/Shorten test failed! Input %v, output %v, reversed %v.", input, output, reversed)
				}
			}
	}
}


func BenchmarkShorten(b *testing.B){
	b.StopTimer()
	const dict ="ABCDEFHIJLKMNOPQRSTUVWXYZ1234567890abcefghijklmnopqrstuvwxyz"
	tester := MakeDictionary(dict)
	b.StartTimer()

	for i:=0;i<b.N;i++{
		input := uint(i)
		output, _ := tester.Shorten(input)

		reversed := tester.Lengthen(output)
		if !(reversed == input){
			b.Errorf("Lengthen/Shorten test failed! Input %v, output %v, reversed %v.", input, output, reversed)
		}
	}
}
