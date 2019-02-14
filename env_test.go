package env

import (
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestMapSource(t *testing.T) {
	e := FromMap(map[string]string{"FOO": "abc"})
	e.String("FOO", nil, "irrelevant")
	err := e.Load()
	if err != nil {
		t.Errorf("e.Load() returned error: %v", err)
	}

	e = FromMap(map[string]string{"FOO": "abc"})
	e.String("BAR", nil, "irrelevant")
	err = e.Load()
	if err == nil {
		t.Fatalf("e.Load() didn't return error for missing variable")
	}

	got := err.Error()
	want := "missing environment variables: BAR"
	if got != want {
		t.Errorf("e.Load() returned error message '%v', want '%v'", got, want)
	}
}

func TestPrefixSource(t *testing.T) {
	e := FromMap(map[string]string{"PRE_FOO": "abc"}).Prefix("PRE_")

	e.String("FOO", nil, "irrelevant")
	err := e.Load()
	if err != nil {
		t.Errorf("e.Load() returned error: %v", err)
	}

	e = FromMap(map[string]string{"FOO": "abc"}).Prefix("PRE_")
	e.String("FOO", nil, "irrelevant")
	err = e.Load()
	if err == nil {
		t.Errorf("e.Load() didn't return error for variable without prefix")
	}

	e = FromMap(map[string]string{"PRE_FOO": "abc"}).Prefix("PRE_")
	e.String("BAR", nil, "irrelevant")
	err = e.Load()
	if err == nil {
		t.Errorf("e.Load() didn't return error for missing variable")
	}

	got := err.Error()
	want := "missing environment variables: PRE_BAR"
	if got != want {
		t.Errorf("e.Load() returned error message '%v', want '%v'", got, want)
	}
}

func TestOSSource(t *testing.T) {
	os.Clearenv()
	os.Setenv("FOO", "abc")

	e := New()
	e.String("FOO", nil, "irrelevant")
	err := e.Load()
	if err != nil {
		t.Errorf("e.Load() returned error: %v", err)
	}

	e = New()
	e.String("BAR", nil, "irrelevant")
	err = e.Load()
	if err == nil {
		t.Errorf("e.Load() didn't return error for missing variable")
	}

	got := err.Error()
	want := "missing environment variables: BAR"
	if got != want {
		t.Errorf("e.Load() returned error message '%v', want '%v'", got, want)
	}
}

func TestHelp(t *testing.T) {
	e := New()
	e.String("HOST", nil, "hostname")
	e.Int("PORT", nil, "port number")
	want := `HOST -- hostname
PORT -- port number
`
	got := e.Help()
	if got != want {
		t.Errorf("e.Help() returned:\n%v\nwant:\n%v", got, want)
	}

	sub := e.Prefix("MATH_").Prefix("CONSTANTS_")
	sub.Float("PI", nil, "the ratio of a circle's circumference to its diameter")
	want = `MATH_CONSTANTS_PI -- the ratio of a circle's circumference to its diameter
`
	got = sub.Help()
	if got != want {
		t.Errorf("sub.Help() returned:\n%v\nwant:\n%v", got, want)
	}
}

func TestString(t *testing.T) {
	want := "abc"
	m := map[string]string{
		"FOO": want,
	}
	var foo string

	e := FromMap(m)
	e.String("FOO", &foo, "irrelevant")
	err := e.Load()
	if err != nil {
		t.Errorf("e.Load() returned error: %v", err)
	} else if foo != want {
		t.Errorf("e.Load() set value to %v, want %v", foo, want)
	}
}

func TestOptionalString(t *testing.T) {
	wantFoo, wantBar := "abc", "def"
	m := map[string]string{
		"FOO": wantFoo,
	}
	var foo, bar string

	e := FromMap(m)
	e.OptionalString("FOO", &foo, "xyz", "irrelevant")
	e.OptionalString("BAR", &bar, wantBar, "irrelevant")
	err := e.Load()
	if err != nil {
		t.Fatalf("e.Load() returned error: %v", err)
	}
	if foo != wantFoo {
		t.Errorf("e.Load() set value to %v, want %v", foo, wantFoo)
	}
	if bar != wantBar {
		t.Errorf("e.Load() set value to %v, want %v", bar, wantBar)
	}
}

func TestInt(t *testing.T) {
	want := 123
	m := map[string]string{
		"FOO": fmt.Sprintf("%v", want),
	}
	var foo int

	e := FromMap(m)
	e.Int("FOO", &foo, "irrelevant")
	err := e.Load()
	if err != nil {
		t.Errorf("e.Load() returned error: %v", err)
	} else if foo != want {
		t.Errorf("e.Load() set value to %v, want %v", foo, want)
	}
}

func TestOptionalInt(t *testing.T) {
	wantFoo, wantBar := 123, 456
	m := map[string]string{
		"FOO": fmt.Sprintf("%v", wantFoo),
	}
	var foo, bar int

	e := FromMap(m)
	e.OptionalInt("FOO", &foo, 999, "irrelevant")
	e.OptionalInt("BAR", &bar, wantBar, "irrelevant")
	err := e.Load()
	if err != nil {
		t.Fatalf("e.Load() returned error: %v", err)
	}
	if foo != wantFoo {
		t.Errorf("e.Load() set value to %v, want %v", foo, wantFoo)
	}
	if bar != wantBar {
		t.Errorf("e.Load() set value to %v, want %v", bar, wantBar)
	}
}

func TestFloat(t *testing.T) {
	want := 123.45
	m := map[string]string{
		"FOO": fmt.Sprintf("%v", want),
	}
	var foo float64

	e := FromMap(m)
	e.Float("FOO", &foo, "irrelevant")
	err := e.Load()
	if err != nil {
		t.Errorf("e.Load() returned error: %v", err)
	} else if foo != want {
		t.Errorf("e.Load() set value to %v, want %v", foo, want)
	}
}

func TestOptionalFloat(t *testing.T) {
	wantFoo, wantBar := 123.45, 456.78
	m := map[string]string{
		"FOO": fmt.Sprintf("%v", wantFoo),
	}
	var foo, bar float64

	e := FromMap(m)
	e.OptionalFloat("FOO", &foo, 999.99, "irrelevant")
	e.OptionalFloat("BAR", &bar, wantBar, "irrelevant")
	err := e.Load()
	if err != nil {
		t.Fatalf("e.Load() returned error: %v", err)
	}
	if foo != wantFoo {
		t.Errorf("e.Load() set value to %v, want %v", foo, wantFoo)
	}
	if bar != wantBar {
		t.Errorf("e.Load() set value to %v, want %v", bar, wantBar)
	}
}

func TestBool(t *testing.T) {
	want := true
	var foo bool

	e := FromMap(map[string]string{"FOO": "true"})
	e.Bool("FOO", &foo, "irrelevant")
	err := e.Load()
	if err != nil {
		t.Errorf("e.Load() returned error: %v", err)
	} else if foo != want {
		t.Errorf("e.Load() set value to %v, want %v", foo, want)
	}
}

func TestOptionalBool(t *testing.T) {
	wantFoo, wantBar := false, true
	m := map[string]string{
		"FOO": fmt.Sprintf("%v", wantFoo),
	}
	var foo, bar bool

	e := FromMap(m)
	e.OptionalBool("FOO", &foo, false, "irrelevant")
	e.OptionalBool("BAR", &bar, wantBar, "irrelevant")
	err := e.Load()
	if err != nil {
		t.Fatalf("e.Load() returned error: %v", err)
	}
	if foo != wantFoo {
		t.Errorf("e.Load() set value to %v, want %v", foo, wantFoo)
	}
	if bar != wantBar {
		t.Errorf("e.Load() set value to %v, want %v", bar, wantBar)
	}
}

func TestFlag(t *testing.T) {
	wantFoo, wantBar := true, false
	var foo, bar bool
	m := map[string]string{
		"FOO": "",
	}

	e := FromMap(m)
	e.Flag("FOO", &foo, "irrelevant")
	e.Flag("BAR", &bar, "irrelevant")
	err := e.Load()
	if err != nil {
		t.Fatalf("e.Load() returned error: %v", err)
	}
	if foo != wantFoo {
		t.Errorf("e.Load() set value to %v, want %v", foo, wantFoo)
	}
	if bar != wantBar {
		t.Errorf("e.Load() set value to %v, want %v", bar, wantBar)
	}
}

func TestList(t *testing.T) {
	want := []string{"de", "fr", "it"}
	m := map[string]string{"FOO": "de,fr,it"}
	var foo []string

	e := FromMap(m)
	e.List("FOO", &foo, ",", "irrelevant")
	err := e.Load()
	if err != nil {
		t.Errorf("e.Load() returned error: %v", err)
	} else if !reflect.DeepEqual(foo, want) {
		t.Errorf("e.Load() set value to %v, want %v", foo, want)
	}
}

func TestOptionalList(t *testing.T) {
	wantFoo, wantBar := []string{"de", "fr", "it"}, []string{"dk", "se"}
	m := map[string]string{"FOO": "de,fr,it"}
	var foo, bar []string

	e := FromMap(m)
	e.OptionalList("FOO", &foo, ",", []string{"ussr"}, "irrelevant")
	e.OptionalList("BAR", &bar, ",", wantBar, "irrelevant")
	err := e.Load()
	if err != nil {
		t.Fatalf("e.Load() returned error: %v", err)
	}
	if !reflect.DeepEqual(foo, wantFoo) {
		t.Errorf("e.Load() set value to %v, want %v", foo, wantFoo)
	}
	if !reflect.DeepEqual(bar, wantBar) {
		t.Errorf("e.Load() set value to %v, want %v", bar, wantBar)
	}
}

func TestSet(t *testing.T) {
	want := map[string]bool{"de": true, "it": true, "fr": true}
	m := map[string]string{"FOO": "de,it,fr"}
	var foo map[string]bool

	e := FromMap(m)
	e.Set("FOO", &foo, ",", "irrelevant")
	err := e.Load()
	if err != nil {
		t.Errorf("e.Load() returned error: %v", err)
	} else if !reflect.DeepEqual(foo, want) {
		t.Errorf("e.Load() set value to %v, want %v", foo, want)
	}
}

func TestOptionalSet(t *testing.T) {
	wantFoo := map[string]bool{"de": true, "it": true, "fr": true}
	wantBar := map[string]bool{"dk": true, "se": true}
	m := map[string]string{"FOO": "de,it,fr"}
	var foo, bar map[string]bool

	e := FromMap(m)
	e.OptionalSet("FOO", &foo, ",", map[string]bool{"sfry": true}, "irrelevant")
	e.OptionalSet("BAR", &bar, ",", wantBar, "irrelevant")
	err := e.Load()
	if err != nil {
		t.Errorf("e.Load() returned error: %v", err)
	}
	if !reflect.DeepEqual(foo, wantFoo) {
		t.Errorf("e.Load() set value to %v, want %v", foo, wantFoo)
	}
	if !reflect.DeepEqual(bar, wantBar) {
		t.Errorf("e.Load() set value to %v, want %v", bar, wantBar)
	}
}
