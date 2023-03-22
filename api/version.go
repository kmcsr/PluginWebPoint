
package api

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	PluginIdRe = regexp.MustCompile("[0-9a-z_]{1,64}")
	VersionRe = regexp.MustCompile(`^[-+0-9A-Za-z*]+(?:\.[-+0-9A-Za-z*]+)*`)
)

type Version struct {
	Major int
	Minor int
	Patch int
	Desc  string
}

var _ json.Unmarshaler = (*Version)(nil)
var _ json.Marshaler = (*Version)(nil)

func VersionFromString(data string)(v Version, err error){
	if !VersionRe.MatchString(data) {
		err = fmt.Errorf("Format error for version %q, not match version regexp", data)
		return
	}
	v.Major = 0
	v.Minor = 0
	v.Patch = 0
	i := strings.IndexAny(data, "-+")
	if i >= 0 {
		data, v.Desc = data[:i], data[i + 1:]
	}
	if i = strings.IndexByte(data, '.'); i < 0 {
		v.Major, err = strconv.Atoi(data)
		return
	}
	if v.Major, err = strconv.Atoi(data[:i]); err != nil {
		return
	}
	data = data[i + 1:]
	if i = strings.IndexByte(data, '.'); i < 0 {
		v.Minor, err = strconv.Atoi(data)
		return
	}
	if v.Minor, err = strconv.Atoi(data[:i]); err != nil {
		return
	}
	if v.Patch, err = strconv.Atoi(data[i + 1:]); err != nil {
		return
	}
	return
}

func (v *Version)UnmarshalJSON(data []byte)(err error){
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return
	}
	var v0 Version
	if v0, err = VersionFromString(s); err != nil {
		return
	}
	*v = v0
	return
}

func (v Version)String()(s string){
	s = strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Patch)
	if len(v.Desc) > 0 {
		s += "-" + v.Desc
	}
	return
}

func (v Version)Value()(driver.Value, error){
	return v.String(), nil
}

func (v Version)MarshalJSON()([]byte, error){
	return ([]byte)("\"" + v.String() + "\""), nil
}

func (v *Version)Scan(d any)(err error){
	s, ok := d.(string)
	if !ok {
		if b, ok := d.([]byte); !ok {
			return fmt.Errorf("Unexpect type %T for Version, expect string or bytes", d)
		}else{
			s = (string)(b)
		}
	}
	var v0 Version
	if v0, err = VersionFromString(s); err != nil {
		return
	}
	*v = v0
	return
}

type Cond uint8

const (
	InvaildCond Cond = iota
	EQ // =
	LT // <
	LE // <=
	GT // >
	GE // >=
	EX // ^
	TD // ~
)

func CondFromString(s string)(Cond){
	switch s {
	case "=":
		return EQ
	case "<":
		return LT
	case "<=":
		return LE
	case ">":
		return GT
	case ">=":
		return GE
	case "^":
		return EX
	case "~":
		return TD
	}
	return InvaildCond
}

func (c Cond)String()(string){
	switch c {
	case EQ:
		return "="
	case LT:
		return "<"
	case LE:
		return "<="
	case GT:
		return ">"
	case GE:
		return ">="
	case EX:
		return "^"
	case TD:
		return "~"
	}
	panic("Invaild condition")
}

type VersionCond struct {
	Cond Cond
	Major int
	Minor int
	Patch int
	Desc  string
}

var VersionMatchAny = VersionCond{
	Cond: EQ,
	Major: -1,
	Minor: -1,
	Patch: -1,
}

func VersionCondFromString(s string)(v VersionCond, err error){
	v.Major = -1
	v.Minor = -1
	v.Patch = -1
	v.Cond = EQ
	switch s[0] {
	case '=':
		if s[1] == '=' {
			s = s[2:]
		}else{
			s = s[1:]
		}
		v.Cond = EQ
	case '^':
		s = s[1:]
		v.Cond = EX
	case '~':
		s = s[1:]
		v.Cond = TD
	case '<':
		if s[1] == '=' {
			s = s[2:]
			v.Cond = LE
		}else{
			s = s[1:]
			v.Cond = LT
		}
	case '>':
		if s[1] == '=' {
			s = s[2:]
			v.Cond = GE
		}else{
			s = s[1:]
			v.Cond = GT
		}
	}
	if !VersionRe.MatchString(s) {
		err = fmt.Errorf("Format error for version %q, not match version regexp", s)
		return
	}
	i := strings.IndexAny(s, "-+")
	if i >= 0 {
		s, v.Desc = s[:i], s[i + 1:]
	}
	if s == "*" || s == "x" {
		return
	}
	if i = strings.IndexByte(s, '.'); i < 0 {
		v.Major, err = strconv.Atoi(s)
		return
	}
	if v.Major, err = strconv.Atoi(s[:i]); err != nil {
		return
	}
	if s = s[i + 1:]; s == "*" || s == "x" {
		return
	}
	if i = strings.IndexByte(s, '.'); i < 0 {
		v.Minor, err = strconv.Atoi(s)
		return
	}
	if v.Minor, err = strconv.Atoi(s[:i]); err != nil {
		return
	}
	if s = s[i + 1:]; s == "*" || s == "x" {
		return
	}
	if v.Patch, err = strconv.Atoi(s); err != nil {
		return
	}
	return
}

func (v *VersionCond)UnmarshalJSON(data []byte)(err error){
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return
	}
	var v0 VersionCond
	if v0, err = VersionCondFromString(s); err != nil {
		return
	}
	*v = v0
	return
}

func (v VersionCond)String()(s string){
	s = v.Cond.String()
	if v.Major < 0 {
		s += "*"
	}else{
		s += strconv.Itoa(v.Major) + "."
		if v.Minor < 0 {
			s += "*"
		}else{
			s += strconv.Itoa(v.Minor) + "."
			if v.Patch < 0 {
				s += "*"
			}else{
				s += strconv.Itoa(v.Patch)
			}
		}
	}
	if len(v.Desc) > 0 {
		s += "-" + v.Desc
	}
	return
}

func (v VersionCond)MarshalJSON()([]byte, error){
	return ([]byte)("\"" + v.String() + "\""), nil
}

func (v *VersionCond)Scan(d any)(err error){
	s, ok := d.(string)
	if !ok {
		if b, ok := d.([]byte); !ok {
			return fmt.Errorf("Unexpect type %T for Version, expect string or bytes", d)
		}else{
			s = (string)(b)
		}
	}
	var v0 VersionCond
	if v0, err = VersionCondFromString(s); err != nil {
		return
	}
	*v = v0
	return
}

func (v VersionCond)Value()(driver.Value, error){
	return v.String(), nil
}

func (vc VersionCond)IsMatch(v Version)(bool){
	if vc.Major == -1 {
		return true
	}
	switch vc.Cond {
	case EQ:
		return (vc.Major == v.Major &&
			(vc.Minor == -1 || (vc.Minor == v.Minor &&
				(vc.Patch == -1 || vc.Patch == v.Patch))))
	case LT:
		return vc.Major < v.Major || (vc.Major == v.Major &&
			(vc.Minor == -1 || vc.Minor < v.Minor || (vc.Minor == v.Minor &&
				(vc.Patch == -1 || vc.Patch < v.Patch))))
	case LE:
		return vc.Major < v.Major || (vc.Major == v.Major &&
			(vc.Minor == -1 || vc.Minor < v.Minor || (vc.Minor == v.Minor &&
				(vc.Patch == -1 || vc.Patch <= v.Patch))))
	case GT:
		return vc.Major > v.Major || (vc.Major == v.Major &&
			(vc.Minor == -1 || vc.Minor > v.Minor || (vc.Minor == v.Minor &&
				(vc.Patch == -1 || vc.Patch > v.Patch))))
	case GE:
		return vc.Major >= v.Major || (vc.Major == v.Major &&
			(vc.Minor == -1 || vc.Minor >= v.Minor || (vc.Minor == v.Minor &&
				(vc.Patch == -1 || vc.Patch >= v.Patch))))
	case EX:
		return vc.Major >= v.Major
	case TD:
		return (vc.Major == v.Major && vc.Minor >= v.Minor)
	default:
		panic("Invaild condition")
	}
}
