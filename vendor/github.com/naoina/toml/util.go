package toml

import (
	"fmt"
	"reflect"
	"strings"
)

const fieldTagName = "toml"

// fieldCache maps normalized field names to their position in a struct.
type fieldCache struct {
	named map[string]fieldInfo // fields with an explicit name in tag
	auto  map[string]fieldInfo // fields with auto-assigned normalized names
}

type fieldInfo struct {
	index   []int
	name    string
	ignored bool
}

func makeFieldCache(rt reflect.Type) fieldCache {
	named, auto := make(map[string]fieldInfo), make(map[string]fieldInfo)
	for i := 0; i < rt.NumField(); i++ {
		ft := rt.Field(i)
		// skip unexported fields
		if ft.PkgPath != "" && !ft.Anonymous {
			continue
		}
		col, _ := extractTag(ft.Tag.Get(fieldTagName))
		info := fieldInfo{index: ft.Index, name: ft.Name, ignored: col == "-"}
		if col == "" || col == "-" {
			auto[normFieldName(ft.Name)] = info
		} else {
			named[col] = info
		}
	}
	return fieldCache{named, auto}
}

func (fc fieldCache) findField(rv reflect.Value, name string) (reflect.Value, string, error) {
	info, found := fc.named[name]
	if !found {
		info, found = fc.auto[normFieldName(name)]
	}
	if !found {
		return reflect.Value{}, "", fmt.Errorf("field corresponding to `%s' is not defined in %v", name, rv.Type())
	} else if info.ignored {
		return reflect.Value{}, "", fmt.Errorf("field corresponding to `%s' in %v cannot be set through TOML", name, rv.Type())
	}
	return rv.FieldByIndex(info.index), info.name, nil
}

func normFieldName(s string) string {
	return strings.Replace(strings.ToLower(s), "_", "", -1)
}

func extractTag(tag string) (col, rest string) {
	tags := strings.SplitN(tag, ",", 2)
	if len(tags) == 2 {
		return strings.TrimSpace(tags[0]), strings.TrimSpace(tags[1])
	}
	return strings.TrimSpace(tags[0]), ""
}
