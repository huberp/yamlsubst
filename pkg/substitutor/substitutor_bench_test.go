package substitutor

import (
	"strings"
	"testing"
)

func BenchmarkSubstitute_Simple(b *testing.B) {
	yamlContent := `
name: John
age: 30
`
	input := "Hello, my name is ${.name} and I am ${.age} years old"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Substitute(input, yamlContent)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSubstitute_Nested(b *testing.B) {
	yamlContent := `
person:
  name: Alice
  address:
    city: Seattle
    state: WA
`
	input := "I live in ${.person.address.city}, ${.person.address.state}"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Substitute(input, yamlContent)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSubstitute_MultipleOccurrences(b *testing.B) {
	yamlContent := `
word: test
`
	input := strings.Repeat("${.word} ", 100)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Substitute(input, yamlContent)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSubstitute_LargeInput(b *testing.B) {
	yamlContent := `
name: John
age: 30
city: Seattle
`
	template := "Name: ${.name}, Age: ${.age}, City: ${.city}\n"
	input := strings.Repeat(template, 1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Substitute(input, yamlContent)
		if err != nil {
			b.Fatal(err)
		}
	}
}
