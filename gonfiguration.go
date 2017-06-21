package gonfiguration

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

var (
	errorGonfigurationNil = errors.New("error: gonfiguration is nil")
	errorCantOpenFile     = errors.New("error: can't open gonfiguration file")
	errorMapContainsEntry = errors.New("error: map contains that key yet")
	errorKeyNotFound      = errors.New("error: can't found key-value on gonfiguration map")
)

// Gonfiguration hold info like path, map of values
type Gonfiguration struct {
	GonfigurationValues map[string]string
	ValidateConfigLine  func(line string) bool
	FilePath            string
}

// isGonfigurationLine internal func for return true if line is correct like key=value
func (c *Gonfiguration) isGonfigurationLine(line string) bool {
	if c.ValidateConfigLine != nil {
		return c.ValidateConfigLine(line)
	}
	return !strings.HasPrefix(line, "#") && line != "" && strings.Contains(line, "=")
}

// LoadGonfiguration save all key-value into the map
func (c *Gonfiguration) LoadGonfiguration(path string, validateFunction func(line string) bool) (*Gonfiguration, error) {
	newGonfiguration := new(Gonfiguration)
	newGonfiguration.GonfigurationValues = make(map[string]string)
	newGonfiguration.ValidateConfigLine = validateFunction
	newGonfiguration.FilePath = path

	reloadError := c.Reload()
	if reloadError != nil {
		return nil, reloadError
	}
	return newGonfiguration, nil
}

// Reload clear the map and load all entrys of the config file
func (c *Gonfiguration) Reload() error {
	c.Clear()

	file, err := os.Open(c.FilePath)
	if err != nil {
		return errorCantOpenFile
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if c.isGonfigurationLine(line) {
			code := strings.Split(line, "=")[0]
			value := strings.Split(line, "=")[1:]

			if c.ContainsKey(strings.ToLower(code)) {
				return errorMapContainsEntry
			}
			c.GonfigurationValues[strings.ToLower(code)] = strings.Join(value, "=")
		}
	}
	return nil
}

// GetConfigParamAsString return the string value of the given key
func (c *Gonfiguration) GetConfigParamAsString(key string, def string) (string, error) {
	if c == nil || c.GonfigurationValues == nil {
		return "", errorGonfigurationNil
	}

	val := c.GonfigurationValues[strings.ToLower(key)]
	if !c.ContainsKey(strings.ToLower(key)) {
		return def, errorKeyNotFound
	}
	return val, nil
}

// GetConfigParamAsInt return the int value of the given key
func (c *Gonfiguration) GetConfigParamAsInt(key string, def int) (int, error) {
	if c == nil || c.GonfigurationValues == nil {
		return 0, errorGonfigurationNil
	}

	val := c.GonfigurationValues[strings.ToLower(key)]
	if !c.ContainsKey(strings.ToLower(key)) {
		return def, errorKeyNotFound
	}
	intval, _ := strconv.Atoi(val)
	return intval, nil
}

// GetConfigParamAsBool return the bool param of the given key
func (c *Gonfiguration) GetConfigParamAsBool(key string, def bool) (bool, error) {
	if c == nil || c.GonfigurationValues == nil {
		return false, errorGonfigurationNil
	}

	val := c.GonfigurationValues[strings.ToLower(key)]
	if !c.ContainsKey(strings.ToLower(key)) {
		return def, errorKeyNotFound
	}
	return (val == "true" || val == "1"), nil
}

// AddNewKeyValueEntry add new entry to the config map, return error if it exists yet
func (c *Gonfiguration) AddNewKeyValueEntry(key, value string) error {
	if c == nil || c.GonfigurationValues == nil {
		return errorGonfigurationNil
	}

	if !c.ContainsKey(strings.ToLower(key)) {
		c.GonfigurationValues[strings.ToLower(key)] = value
		return nil
	}
	return errorMapContainsEntry
}

// UpdateOrAddEntry add or update the value of the key
func (c *Gonfiguration) UpdateOrAddEntry(key, value string) error {
	if c == nil || c.GonfigurationValues == nil {
		return errorGonfigurationNil
	}
	c.GonfigurationValues[strings.ToLower(key)] = value
	return nil
}

// DeleteEntry delete a entry from the map
func (c *Gonfiguration) DeleteEntry(key string) error {
	if c == nil || c.GonfigurationValues == nil {
		return errorGonfigurationNil
	}
	delete(c.GonfigurationValues, strings.ToLower(key))
	return nil
}

// ContainsKey return true if the key exists on the map
func (c *Gonfiguration) ContainsKey(key string) bool {
	if c == nil || c.GonfigurationValues == nil {
		return false
	}
	_, contains := c.GonfigurationValues[strings.ToLower(key)]
	return contains
}

// GetGonfigurationLen return the len of the config map
func (c *Gonfiguration) GetGonfigurationLen() (int, error) {
	if c == nil || c.GonfigurationValues == nil {
		return 0, errorGonfigurationNil
	}
	return len(c.GonfigurationValues), nil
}

// Clear delete all entrys from the map
func (c *Gonfiguration) Clear() error {
	if c == nil || c.GonfigurationValues == nil {
		return errorGonfigurationNil
	}

	for key := range c.GonfigurationValues {
		delete(c.GonfigurationValues, strings.ToLower(key))
	}
	return nil
}

// ChangePathAndReload change the path of the config file and reload
func (c *Gonfiguration) ChangePathAndReload(newPath string) error {
	c.FilePath = newPath
	reloadError := c.Reload()

	return reloadError
}

// Dispose delete all entry of the map and set it to nil
func (c *Gonfiguration) Dispose() {
	c.GonfigurationValues = nil
	c.ValidateConfigLine = nil
	c.FilePath = ""
}

// GetMap return the map of the Gonfiguration
func (c *Gonfiguration) GetMap() map[string]string {
	return c.GonfigurationValues
}
