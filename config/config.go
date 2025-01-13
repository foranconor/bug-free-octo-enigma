package config

import (
	"fmt"
	"os"

	"github.com/kr/pretty"
	lua "github.com/yuin/gopher-lua"
)

type Params struct {
	Numeric map[string]float64
	Text    map[string]string
}

type Section struct {
	Kind       string
	Steps      int
	StartWidth float64
	EndWidth   float64
	Parameters map[string]float64
}

type Stair struct {
	Sections   []Section
	Parameters Params
}

type StairError struct {
	Kind    string
	Message string
	Bad     []string
}

func (e *StairError) Error() string {
	return e.Message
}

func LoadStair(path string) (*Stair, error) {
	// check if the file exists
	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	// run the file
	L := lua.NewState()
	defer L.Close()
	err = L.DoFile(path)
	if err != nil {
		return nil, err
	}
	// get the top of the stack, what is ultimatly returned from the lua script
	ret := L.Get(-1)
	if ret.Type() != lua.LTTable {
		return nil, &StairError{
			Kind:    "config",
			Message: "script must return a table",
		}
	}
	tab := ret.(*lua.LTable)
	bad := basicCheck(tab)
	if len(bad) > 0 {
		return nil, &StairError{
			Kind:    "config",
			Message: "all terminal config elements must be text or numbers",
			Bad:     bad,
		}
	}
	pretty.Println(tab)
	return nil, nil
}

func basicCheck(table *lua.LTable) []string {
	bad := make([]string, 0)

	accepted := make(map[lua.LValueType]bool)
	accepted[lua.LTString] = true
	accepted[lua.LTNumber] = true

	table.ForEach(func(k, v lua.LValue) {
		if v.Type() == lua.LTTable {
			bad = append(bad, basicCheck(v.(*lua.LTable))...)
		} else {
			_, ok := accepted[v.Type()]
			if !ok {
				pretty.Println(k, v.Type())
				bad = append(bad, fmt.Sprintf("%s = %s", k.String(), v.Type().String()))
			}
		}
	})
	return bad
}
