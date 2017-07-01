package gonfiguration

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

var (
	errorGonfigurationNil = errors.New("gonfiguration object is not initialized")
	errorCantOpenFile     = errors.New("can't open configuration file")
	errorDuplicatedEntry  = errors.New("duplicated key")
	errorKeyNotFound      = errors.New("key doesn't exist")
)

// Gonfiguration hold info about the session
type Gonfiguration struct {
	GonfigurationValues map[string]string
	ValidateConfigLine  func(line string) bool
	Path                string
}

// isGonfigurationLine internal func for return true if line is correct like key=value
func (g *Gonfiguration) isGonfigurationLine(line string) bool {
	if g.ValidateConfigLine != nil {
		return g.ValidateConfigLine(line)
	}
	return !strings.HasPrefix(line, "#") && line != "" && strings.Contains(line, "=")
}

// New return the initialized gonfiguration object
func New(path string, validateFunction func(line string) bool) (*Gonfiguration, error) {
	newGonfiguration := new(Gonfiguration)
	newGonfiguration.GonfigurationValues = make(map[string]string)
	newGonfiguration.ValidateConfigLine = validateFunction
	newGonfiguration.Path = path
	errorLoadingGonfiguration := newGonfiguration.LoadFromPath(path)

	return newGonfiguration, errorLoadingGonfiguration
}

// LoadFromPath load the file given without clear the map
func (g *Gonfiguration) LoadFromPath(path string) error {
	g.Path = path

	file, err := os.Open(path)
	if err != nil {
		return errorCantOpenFile
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if g.isGonfigurationLine(line) {
			key := strings.Split(line, "=")[0]
			value := strings.Split(line, "=")[1:]

			contains, _ := g.Contains(key)
			if contains {
				return errorDuplicatedEntry
			}
			g.AddNew(strings.ToLower(key), strings.Join(value, "="))
		}
	}
	return nil
}

// Reload clear the map and call LoadFromPath()
func (g *Gonfiguration) Reload() error {
	if g == nil || g.GonfigurationValues == nil {
		return errorGonfigurationNil
	}

	g.Clear()
	errorLoadFromPath := g.LoadFromPath(g.Path)
	return errorLoadFromPath
}

// GetParamAsString return string value of the key.
func (g *Gonfiguration) GetParamAsString(key string, def string) (string, error) {
	if g == nil || g.GonfigurationValues == nil {
		return "", errorGonfigurationNil
	}

	val, contains := g.GonfigurationValues[strings.ToLower(key)]
	if !contains {
		return def, errorKeyNotFound
	}
	return val, nil
}

// GetParamAsStringArray return the value of the key divided by the sep
func (g *Gonfiguration) GetParamAsStringArray(key, def, sep string) ([]string, error) {
	if g == nil || g.GonfigurationValues == nil {
		return nil, errorGonfigurationNil
	}

	val, contains := g.GonfigurationValues[strings.ToLower(key)]
	if !contains {
		return strings.Split(def, sep), errorKeyNotFound
	}
	return strings.Split(val, sep), nil
}

// GetParamAsInt return int value of the key.
func (g *Gonfiguration) GetParamAsInt(key string, def int) (int, error) {
	if g == nil || g.GonfigurationValues == nil {
		return 0, errorGonfigurationNil
	}

	val, contains := g.GonfigurationValues[strings.ToLower(key)]
	if !contains {
		return def, errorKeyNotFound
	}
	intval, _ := strconv.Atoi(val)
	return intval, nil
}

// GetParamAsIntArray return the value of the key divided by the sep
func (g *Gonfiguration) GetParamAsIntArray(key, def, sep string) ([]int, error) {
	if g == nil || g.GonfigurationValues == nil {
		return nil, errorGonfigurationNil
	}

	tmp, errorGetStringArray := g.GetParamAsStringArray(key, def, sep)
	r := make([]int, len(tmp))
	for x := 0; x < len(tmp); x++ {
		num, _ := strconv.Atoi(tmp[x])
		r[x] = num
	}
	return r, errorGetStringArray
}

// GetParamAsBool return bool value of the key
func (g *Gonfiguration) GetParamAsBool(key string, def bool) (bool, error) {
	if g == nil || g.GonfigurationValues == nil {
		return false, errorGonfigurationNil
	}

	val, contains := g.GonfigurationValues[strings.ToLower(key)]
	if !contains {
		return def, errorKeyNotFound
	}
	return (val == "true" || val == "1"), nil
}

// AddNew add new entry to the config map, return error if it exists yet
func (g *Gonfiguration) AddNew(key, value string) error {
	if g == nil || g.GonfigurationValues == nil {
		return errorGonfigurationNil
	}

	contains, _ := g.Contains(key)
	if contains {
		return errorDuplicatedEntry
	}
	g.GonfigurationValues[strings.ToLower(key)] = value
	return nil
}

// Update add or update the value of the key
func (g *Gonfiguration) Update(key, value string) error {
	if g == nil || g.GonfigurationValues == nil {
		return errorGonfigurationNil
	}
	g.GonfigurationValues[strings.ToLower(key)] = value
	return nil
}

// Delete remove a entry from the map
func (g *Gonfiguration) Delete(key string) error {
	if g == nil || g.GonfigurationValues == nil {
		return errorGonfigurationNil
	}

	contains, _ := g.Contains(key)
	if !contains {
		return errorKeyNotFound
	}
	delete(g.GonfigurationValues, strings.ToLower(key))
	return nil
}

// Contains return true if the key exists on the map
func (g *Gonfiguration) Contains(key string) (bool, error) {
	if g == nil || g.GonfigurationValues == nil {
		return false, errorGonfigurationNil
	}

	_, contains := g.GonfigurationValues[strings.ToLower(key)]
	return contains, nil
}

// Clear delete all entrys from the map
func (g *Gonfiguration) Clear() error {
	if g == nil || g.GonfigurationValues == nil {
		return errorGonfigurationNil
	}

	for key := range g.GonfigurationValues {
		delete(g.GonfigurationValues, strings.ToLower(key))
	}
	return nil
}

// Len return the len of the config map
func (g *Gonfiguration) Len() (int, error) {
	if g == nil || g.GonfigurationValues == nil {
		return 0, errorGonfigurationNil
	}

	return len(g.GonfigurationValues), nil
}

// Map return the map of the Gonfiguration
func (g *Gonfiguration) Map() (*map[string]string, error) {
	if g == nil || g.GonfigurationValues == nil {
		return nil, errorGonfigurationNil
	}

	return &g.GonfigurationValues, nil
}

// Dispose remove all memory used by the object and destroy it
func (g *Gonfiguration) Dispose() {
	g.Clear()
	g.GonfigurationValues = nil
	g.ValidateConfigLine = nil
	g.Path = ""
	g = nil
}
