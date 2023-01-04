package fdisk

import (
	"encoding/json"
	"fmt"

	"github.com/bhendo/go-powershell"
	"github.com/bhendo/go-powershell/backend"
)

type wmiObject map[string]interface{}
type wmiObjects []wmiObject

func makeWmiObject(rawData interface{}) (result wmiObject, err error) {
	switch t := rawData.(type) {
	case map[string]interface{}:
		return t, nil
	case map[interface{}]interface{}:
		result = make(map[string]interface{})
	}
	return result, nil
}

func (obj wmiObject) GetString(key string) string {
	if value, ok := obj[key]; ok {
		if s, ok := value.(string); ok {
			return s
		}
	}
	return ""
}

func (obj wmiObject) GetInt64(key string) int64 {
	if value, ok := obj[key]; ok {
		switch s := value.(type) {
		case int:
			return int64(s)
		case int64:
			return s
		case float64:
			return int64(s)
		case float32:
			return int64(s)
		}
	}
	return -1
}

func (obj wmiObject) Dump() {
	for key, value := range obj {
		fmt.Printf("%s: %v\n", key, value)
	}
}

func (objs wmiObjects) GetObject(key string, value interface{}) (result wmiObject, err error) {
	for _, obj := range objs {
		if obj[key] == value {
			if result != nil {
				return nil, fmt.Errorf("Found multiple objects with %s = %v", key, value)
			}
			result = obj
		}
	}
	if result == nil {
		return nil, fmt.Errorf("Found no object with %s = %v", key, value)
	}
	return result, nil
}

type powerShell struct {
	shell powershell.Shell
}

func NewPowerShell() (ps *powerShell, err error) {
	ps = &powerShell{}
	ps.shell, err = powershell.New(&backend.Local{})
	if err != nil {
		return nil, err
	}
	return ps, nil
}

func (ps *powerShell) Close() {
	ps.shell.Exit()
}

func (ps *powerShell) GetArray(cmd string) (result wmiObjects, err error) {
	stdout, _, err := ps.shell.Execute(cmd + " | ConvertTo-Json")
	if err != nil {
		return nil, err
	}
	var rawData interface{}
	json.Unmarshal([]byte(stdout), &rawData)
	switch t := rawData.(type) {
	case []interface{}:
		for _, obj := range t {
			wmiObj, err := makeWmiObject(obj)
			if err != nil {
				return nil, err
			}
			result = append(result, wmiObj)
		}
	case map[string]interface{}:
		wmiObj, err := makeWmiObject(t)
		if err != nil {
			return nil, err
		}
		result = append(result, wmiObj)
	default:
		return nil, fmt.Errorf("Could not parse powershell output")
	}

	return result, nil
}

func (ps *powerShell) GetObject(cmd string) (wmiObject, error) {
	arr, err := ps.GetArray(cmd)
	if err != nil {
		return nil, err
	}
	if len(arr) != 1 {
		return nil, fmt.Errorf("Found %d objects, but I need exactly 1!", len(arr))
	}
	return arr[0], nil
}

func makeWmiQuery(query string, args ...interface{}) string {
	if len(args) > 0 {
		query = fmt.Sprintf(query, args...)
	}
	return "Get-WmiObject -query \"" + query + "\""
}

func (ps *powerShell) WmiQueryArray(query string, args ...interface{}) (wmiObjects, error) {
	return ps.GetArray(makeWmiQuery(query, args...))
}

func (ps *powerShell) WmiQueryObject(query string, args ...interface{}) (wmiObject, error) {
	return ps.GetObject(makeWmiQuery(query, args...))
}
