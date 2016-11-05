package generator

import (
	"fmt"
	"go/parser"
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	testExample = `example_test.go`
)

// generatorTestSuite
type generatorTestSuite struct {
	suite.Suite
}

// SetupSuite
func (s *generatorTestSuite) SetupSuite() {
}

// TestLengthTestSuite
func TestGeneratorTestSuite(t *testing.T) {
	suite.Run(t, new(generatorTestSuite))
}

// TestNoStructInputFile
func (s *generatorTestSuite) TestNoStructFile() {
	input := `package test
	// SomeInterface
	type SomeInterface interface{

	}
	`
	g := NewGenerator()
	f, err := parser.ParseFile(g.fileSet, "TestRequiredErrors", input, parser.ParseComments)
	s.Nil(err, "Error parsing no struct input")

	output, err := g.Generate(f)
	s.Nil(err, "Error generating formatted code")
	if false { // Debugging statement
		fmt.Println(string(output))
	}
}

// TestNoFile
func (s *generatorTestSuite) TestNoFile() {
	g := NewGenerator()
	// Parse the file given in arguments
	_, err := g.GenerateFromFile("")
	s.NotNil(err, "Error generating formatted code")
}

// TestExampleFile
func (s *generatorTestSuite) TestExampleFile() {
	g := NewGenerator()
	// Parse the file given in arguments
	imported, err := g.GenerateFromFile(testExample)
	s.Nil(err, "Error generating formatted code")
	if false {
		fmt.Println(string(imported))
	}
}

func (s *generatorTestSuite) TestDuplicateRuleFailure() {
	g := NewGenerator()
	input := `package test
	// SomeStruct
	type SomeStruct struct {
		TestString      string             ` + "`valid:\"len=0,len=1\"`" + `
	}
	`
	f, err := parser.ParseFile(g.fileSet, "TestStringInput", input, parser.ParseComments)
	s.Nil(err, "Error parsing input string")

	_, err = g.Generate(f)
	s.EqualError(err, "Duplicate rules are not allowed: 'len' on field 'TestString'")
}

func (s *generatorTestSuite) TestNoTemplate() {
	g := NewGenerator()
	input := `package test
	// SomeStruct
	type SomeStruct struct {
		TestString      string             ` + "`valid:\"xyz\"`" + `
	}
	`
	f, err := parser.ParseFile(g.fileSet, "TestStringInput", input, parser.ParseComments)
	s.Nil(err, "Error parsing input string")

	output, err := g.Generate(f)
	s.Require().Nil(err, "Error generating output")

	s.Equal(len("package test\n"), len(output))
	s.Equal("package test\n", string(output))
}

func (s *generatorTestSuite) TestPointerFunc() {
	g := NewGenerator().WithPointerMethod()
	input := `package test
	// SomeStruct
	type SomeStruct struct {
		TestString      string             ` + "`valid:\"required\"`" + `
	}
	`
	f, err := parser.ParseFile(g.fileSet, "TestStringInput", input, parser.ParseComments)
	s.Nil(err, "Error parsing input string")

	output, err := g.Generate(f)
	s.Require().Nil(err, "Error generating output")

	s.Contains(string(output), "func (s *SomeStruct) Validate()")

}

func (s *generatorTestSuite) TestAddTemplateFile() {
	g := NewGenerator()
	input := `package test
	// Inner
	type Inner struct{

	}

	// SomeStruct
	type SomeStruct struct {
		TestString      *Inner             ` + "`valid:\"dummy\"`" + `
	}
	`

	err := g.AddTemplateFiles("dummy_template.tmpl")
	s.Require().Nil(err, "Error adding dummy template")
	f, err := parser.ParseFile(g.fileSet, "TestStringInput", input, parser.ParseComments)
	s.Nil(err, "Error parsing input string")

	output, err := g.Generate(f)
	s.Require().NotNil(err, fmt.Sprintf("Should have had an error generating output"))
	s.Contains(err.Error(), "generate: error formatting code")
	if false {
		fmt.Printf("Output: %s\n", string(output))
	}
}

var result bool

var EmptyStructs = []Inner{
	Inner{},
	Inner{},
	Inner{},
	Inner{},
	Inner{},
	Inner{},
	Inner{},
	Inner{},
	Inner{},
	Inner{},
}

// BenchmarkReflection is a quick test to see how much of an impact reflection has
// in the performance of an application
func BenchmarkReflectionInt(b *testing.B) {
	match := false
	for x := 0; x < b.N; x++ {
		zeroInt := reflect.Zero(reflect.ValueOf(x).Type())
		if reflect.ValueOf(x).Interface() == zeroInt.Interface() {
			match = true
		}
		match = false
	}
	result = match
}

// BenchmarkEmptyInt is a quick benchmark to determine performance of just using an empty var
// for zero value comparison
func BenchmarkEmptyInt(b *testing.B) {
	match := false
	for x := 0; x < b.N; x++ {
		var zeroInt int
		if x == zeroInt {
			match = true
		}
		match = false
	}
	result = match
}

// BenchmarkReflectionStruct is a quick test to see how much of an impact reflection has
// in the performance of an application
func BenchmarkReflectionStruct(b *testing.B) {
	match := false
	count := len(EmptyStructs)

	for x := 0; x < b.N; x++ {
		uut := EmptyStructs[x%count]
		zeroInner := reflect.Zero(reflect.ValueOf(uut).Type())
		if reflect.ValueOf(uut).Interface() == zeroInner.Interface() {
			match = true
		}
		match = false
	}
	result = match
}

// BenchmarkEmptyStruct is a quick benchmark to determine performance of just using an empty var
// for zero value comparison
func BenchmarkEmptyStruct(b *testing.B) {
	match := false
	count := len(EmptyStructs)

	for x := 0; x < b.N; x++ {
		var zeroExample Inner
		if EmptyStructs[x%count] == zeroExample {
			match = true
		}
		match = false
	}
	result = match
}

var Strings = []string{
	"",
	"a",
	"abcdefghijklmnopqrstuvwxyz",
	"1234567890",
	"qwerty",
	"zxcvb",
	"chickens",
	"cows",
	"trains",
}

// BenchmarkReflectionString is a quick test to see how much of an impact reflection has
// in the performance of an application
func BenchmarkReflectionString(b *testing.B) {
	match := false
	count := len(Strings)

	for x := 0; x < b.N; x++ {
		uut := Strings[x%count]
		zeroInner := reflect.Zero(reflect.ValueOf(uut).Type())
		if reflect.ValueOf(uut).Interface() == zeroInner.Interface() {
			match = true
		}
		match = false
	}
	result = match
}

// BenchmarkEmptyString is a quick benchmark to determine performance of just using an empty var
// for zero value comparison
func BenchmarkEmptyString(b *testing.B) {
	match := false
	count := len(Strings)

	for x := 0; x < b.N; x++ {
		var zeroExample string
		if Strings[x%count] == zeroExample {
			match = true
		}
		match = false
	}
	result = match
}
