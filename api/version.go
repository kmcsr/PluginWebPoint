
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
	PluginIdRe = regexp.MustCompile(`[0-9a-z_]{1,64}`)
	VersionRe = regexp.MustCompile(`^((?:[xX*]|\d+)(?:\.(?:[xX*]|\d+))*)([-+][-+.0-9A-Za-z]+)?$`)
	VersionExtraRe = regexp.MustCompile(`^|[-+0-9A-Za-z]+(?:\.[-+0-9A-Za-z]+)*$`)
)

type Version struct {
	Comps []int
	HasWildcard bool

	Pre   string
	Build string
}

var _ json.Unmarshaler = (*Version)(nil)
var _ json.Marshaler = (*Version)(nil)

func split(str string, s byte)(a, b string){
	i := strings.IndexByte(str, s)
	if i < 0 {
		return str, ""
	}
	return str[:i], str[i + 1:]
}

func VersionFromString(data string)(v Version, err error){
	if !VersionRe.MatchString(data) {
		err = fmt.Errorf("Format error for %q, not match the version regexp", data)
		return
	}
	data, v.Build = split(data, '+')
	data, v.Pre = split(data, '-')
	var n int
	for _, x := range strings.Split(data, ".") {
		switch x[0] {
		case '*', 'x', 'X':
			v.HasWildcard = true
			n = -1
		default:
			if n, err = strconv.Atoi(x); err != nil {
				return
			}
		}
		v.Comps = append(v.Comps, n)
	}
	return
}

func (v Version)Get(i int)(int){
	if i >= len(v.Comps) {
		return -1
	}
	return v.Comps[i]
}

func (v Version)Less(o Version)(bool){
	max := len(v.Comps)
	if m := len(o.Comps); max < m {
		max = m
	}
	for i := 0; i < max; i++ {
		a, b := v.Get(i), o.Get(i)
		if a >= 0 && b >= 0 && a != b {
			return a < b
		}
	}
	if len(v.Pre) != 0 {
		if len(o.Pre) != 0 {
			return v.Pre < o.Pre
		}
		return !o.HasWildcard
	}
	return false
}

func (v Version)Equal(o Version)(bool){
	return !(v.Less(o) || o.Less(v))
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
	var sb strings.Builder
	for i, n := range v.Comps {
		if i != 0 {
			sb.WriteByte('.')
		}
		if n < 0 {
			sb.WriteByte('*')
		}else{
			sb.WriteString(strconv.Itoa(n))
		}
	}
	if len(v.Pre) > 0 {
		sb.WriteByte('-')
		sb.WriteString(v.Pre)
	}
	if len(v.Build) > 0 {
		sb.WriteByte('+')
		sb.WriteString(v.Build)
	}
	return sb.String()
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
	Ver  Version
}

var _ json.Unmarshaler = (*VersionCond)(nil)
var _ json.Marshaler = (*VersionCond)(nil)

var VersionMatchAny = VersionCond{
	Cond: EQ,
}

func VersionCondFromString(s string)(v VersionCond, err error){
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
	if v.Ver, err = VersionFromString(s); err != nil {
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
	s = v.Cond.String() + v.Ver.String()
	return s
}

func (v VersionCond)MarshalJSON()([]byte, error){
	return ([]byte)("\"" + v.String() + "\""), nil
}

func (v *VersionCond)Scan(d any)(err error){
	s, ok := d.(string)
	if !ok {
		if b, ok := d.([]byte); !ok {
			return fmt.Errorf("Unexpect type %T for VersionCond, expect string or bytes", d)
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
	switch vc.Cond {
	case LE:
		return !v.Less(vc.Ver)
	case GE:
		return !vc.Ver.Less(v)
	case LT:
		return vc.Ver.Less(v)
	case GT:
		return v.Less(vc.Ver)
	case EQ:
		return vc.Ver.Equal(v)
	case EX:
		return !v.Less(vc.Ver) && vc.Ver.Get(0) == v.Get(0)
	case TD:
		return !v.Less(vc.Ver) && vc.Ver.Get(0) == v.Get(0) && vc.Ver.Get(1) == v.Get(1)
	}
	panic("Unexpect version condition")
}

type VersionCondList []VersionCond

var _ json.Unmarshaler = (*VersionCondList)(nil)
var _ json.Marshaler = (*VersionCondList)(nil)

func VersionCondListFromString(s string)(v VersionCondList, err error){
	ss := strings.Split(s, " ")
	v = make(VersionCondList, len(ss))
	for i, s := range ss {
		if v[i], err = VersionCondFromString(s); err != nil {
			return
		}
	}
	return
}

func (v *VersionCondList)UnmarshalJSON(data []byte)(err error){
	var s string
	if err = json.Unmarshal(data, &s); err != nil {
		return
	}
	var v0 VersionCondList
	if v0, err = VersionCondListFromString(s); err != nil {
		return
	}
	*v = v0
	return
}

func (v VersionCondList)String()(s string){
	var sb strings.Builder
	sb.Grow(len(v) * 8)
	for i, vc := range v {
		if i != 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(vc.String())
	}
	return sb.String()
}
func (v VersionCondList)MarshalJSON()([]byte, error){
	return ([]byte)("\"" + v.String() + "\""), nil
}

func (v *VersionCondList)Scan(d any)(err error){
	s, ok := d.(string)
	if !ok {
		if b, ok := d.([]byte); !ok {
			return fmt.Errorf("Unexpect type %T for VersionCondList, expect string or bytes", d)
		}else{
			s = (string)(b)
		}
	}
	var v0 VersionCondList
	if v0, err = VersionCondListFromString(s); err != nil {
		return
	}
	*v = v0
	return
}

func (v VersionCondList)Value()(driver.Value, error){
	return v.String(), nil
}

func (vl VersionCondList)IsMatch(v Version)(bool){
	for _, vc := range vl {
		if !vc.IsMatch(v) {
			return false
		}
	}
	return true
}

