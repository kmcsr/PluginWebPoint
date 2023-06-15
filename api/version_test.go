
package api_test

import (
	"reflect"
	"testing"
	api "github.com/kmcsr/PluginWebPoint/api"
)

type (
	V = api.Version
	VC = api.VersionCond
	VCL = api.VersionCondList
)

var (
	Vnil = V{}
	VCnil = VC{}
)

func versionMemEq(a, b V)(bool){
	return reflect.DeepEqual(a, b)
}

func TestVersionFromString(t *testing.T){
	type T struct {
		S string
		V V
		E bool
	}
	data := []T{
		{ "1.2.3", V{[]int{1, 2, 3}, false, "", ""}, false },
		{ "123.456.789", V{[]int{123, 456, 789}, false, "", ""}, false },
		{ "1.0.0.0", V{[]int{1, 0, 0, 0}, false, "", ""}, false },
		{ "a1.0.0", Vnil, true },
		{ "1.b0.0", Vnil, true },
		{ "1.0.c0", Vnil, true },
		{ "1.0.0.d", Vnil, true },
		{ "1.0.0-pre", V{[]int{1, 0, 0}, false, "pre", ""}, false },
		{ "1.0.0+build", V{[]int{1, 0, 0}, false, "", "build"}, false },
		{ "1.0.0+-build", V{[]int{1, 0, 0}, false, "", "-build"}, false },
		{ "1.0.0-abcd+build", V{[]int{1, 0, 0}, false, "abcd", "build"}, false },
		{ "1.0.-", Vnil, true },
		{ "1.0.+", Vnil, true },
		{ "1.0.*", V{[]int{1, 0, -1}, true, "", ""}, false },
		{ "1.*.0.2", V{[]int{1, -1, 0, 2}, true, "", ""}, false },
		{ "1.0*", Vnil, true },
		{ "1.*0", Vnil, true },
		{ "1.*0.2", Vnil, true },
	}
	for _, d := range data {
		v, e := api.VersionFromString(d.S)
		if d.E {
			if e == nil {
				t.Errorf("Expect error when parsing version %q, but got %#v", d.S, v)
			}
		}else if e != nil {
			t.Errorf("Unexpect error when parsing version %q: %v", d.S, e)
		}else if !versionMemEq(v, d.V) {
			t.Errorf("Unexpect version %#v when parsing version %q,\n  expect %#v", v, d.S, d.V)
		}
	}
}

func TestVersionToString(t *testing.T){
	type T struct {
		V V
		S string
	}
	data := []T{
		{ V{[]int{1, 2, 3}, false, "", ""}, "1.2.3" },
		{ V{[]int{123, 456, 789}, false, "", ""}, "123.456.789" },
		{ V{[]int{1, 0, 0, 0}, false, "", ""}, "1.0.0.0" },
		{ V{[]int{1, 0, 0}, false, "pre", ""}, "1.0.0-pre" },
		{ V{[]int{1, 0, 0}, false, "", "build"}, "1.0.0+build" },
		{ V{[]int{1, 0, 0}, false, "", "-build"}, "1.0.0+-build" },
		{ V{[]int{1, 0, 0}, false, "abcd", "build"}, "1.0.0-abcd+build" },
		{ V{[]int{1, 0, -1}, true, "", ""}, "1.0.*" },
		{ V{[]int{1, -1, 0, 2}, true, "", ""}, "1.*.0.2" },
	}
	for _, d := range data {
		s := d.V.String()
		if s != d.S {
			t.Errorf("Unexpect version %q when call String() on version %#v, expect %q", s, d.V, d.S)
		}
	}
}
