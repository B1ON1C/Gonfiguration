package configuration

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

var (
	errorConfigurationNil = errors.New("error: configuration is nil")
	errorCantOpenFile     = errors.New("error: can't open configuration file")
	errorMapContainsEntry = errors.New("error: map contains that key yet")
	errorKeyNotFound      = errors.New("error: can't found key-value on configuration map")
)

// Configuration hold info like path, map of values
type Configuration struct {
	ConfigurationValues map[string]string
	ValidateConfigLine  func(line string) bool
	FilePath            string
}

// isConfigurationLine internal func for return true if line is correct like key=value
func (c *Configuration) isConfigurationLine(line string) bool {
	if c.ValidateConfigLine != nil {
		return c.ValidateConfigLine(line)
	}
	return !strings.HasPrefix(line, "#") && line != "" && strings.Contains(line, "=")
}

// LoadConfiguration save all key-value into the map
func (c *Configuration) LoadConfiguration(path string, validateFunction func(line string) bool) (*Configuration, error) {
	newConfiguration := new(Configuration)
	newConfiguration.ConfigurationValues = make(map[string]string)
	newConfiguration.ValidateConfigLine = validateFunction
	newConfiguration.FilePath = path

	reloadError := c.Reload()
	if reloadError != nil {
		return nil, reloadError
	}
	return newConfiguration, nil
}

// GetConfigParamAsString return the string value of the given key
func (c *Configuration) GetConfigParamAsString(key string, def string) (string, error) {
	if c == nil || c.ConfigurationValues == nil {
		return "", errorConfigurationNil
	}

	val := c.ConfigurationValues[strings.ToLower(key)]
	if !c.ContainsKey(strings.ToLower(key)) {
		return def, errorKeyNotFound
	}
	return val, nil
}

// GetConfigParamAsInt return the int value of the given key
func (c *Configuration) GetConfigParamAsInt(key string, def int) (int, error) {
	if c == nil || c.ConfigurationValues == nil {
		return 0, errorConfigurationNil
	}

	val := c.ConfigurationValues[strings.ToLower(key)]
	if !c.ContainsKey(strings.ToLower(key)) {
		return def, errorKeyNotFound
	}
	intval, _ := strconv.Atoi(val)
	return intval, nil
}

// GetConfigParamAsBool return the bool param of the given key
func (c *Configuration) GetConfigParamAsBool(key string, def bool) (bool, error) {
	if c == nil || c.ConfigurationValues == nil {
		return false, errorConfigurationNil
	}

	val := c.ConfigurationValues[strings.ToLower(key)]
	if !c.ContainsKey(strings.ToLower(key)) {
		return def, errorKeyNotFound
	}
	return (val == "true" || val == "1"), nil
}

// AddNewKeyValueEntry add new entry to the config map, return error if it exists yet
func (c *Configuration) AddNewKeyValueEntry(key, value string) error {
	if c == nil || c.ConfigurationValues == nil {
		return errorConfigurationNil
	}

	if !c.ContainsKey(strings.ToLower(key)) {
		c.ConfigurationValues[strings.ToLower(key)] = value
		return nil
	}
	return errorMapContainsEntry
}

// UpdateOrAddEntry add or update the value of the key
func (c *Configuration) UpdateOrAddEntry(key, value string) error {
	if c == nil || c.ConfigurationValues == nil {
		return errorConfigurationNil
	}
	c.ConfigurationValues[strings.ToLower(key)] = value
	return nil
}

// DeleteEntry delete a entry from the map
func (c *Configuration) DeleteEntry(key string) error {
	if c == nil || c.ConfigurationValues == nil {
		return errorConfigurationNil
	}
	delete(c.ConfigurationValues, strings.ToLower(key))
	return nil
}

// ContainsKey return true if the key exists on the map
func (c *Configuration) ContainsKey(key string) bool {
	if c == nil || c.ConfigurationValues == nil {
		return false
	}
	_, contains := c.ConfigurationValues[strings.ToLower(key)]
	return contains
}

// GetConfigurationLen return the len of the config map
func (c *Configuration) GetConfigurationLen() (int, error) {
	if c == nil || c.ConfigurationValues == nil {
		return 0, errorConfigurationNil
	}
	return len(c.ConfigurationValues), nil
}

// Reload clear the map and load all entrys of the config file
func (c *Configuration) Reload() error {
	c.Clear()

	file, err := os.Open(c.FilePath)
	if err != nil {
		return errorCantOpenFile
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if c.isConfigurationLine(line) {
			code := strings.Split(line, "=")[0]
			value := strings.Split(line, "=")[1:]
			c.ConfigurationValues[strings.ToLower(code)] = strings.Join(value, "=")
		}
	}
	return nil
}

// Clear delete all entrys from the map
func (c *Configuration) Clear() error {
	if c == nil || c.ConfigurationValues == nil {
		return errorConfigurationNil
	}

	for key := range c.ConfigurationValues {
		delete(c.ConfigurationValues, strings.ToLower(key))
	}
	return nil
}

// ChangePathAndReload change the path of the config file and reload
func (c *Configuration) ChangePathAndReload(newPath string) error {
	c.FilePath = newPath
	reloadError := c.Reload()

	return reloadError
}

// Dispose delete all entry of the map and set it to nil
func (c *Configuration) Dispose() {
	c.ConfigurationValues = nil
	c.ValidateConfigLine = nil
	c.FilePath = ""
}

// GetMap return the map of the configuration
func (c *Configuration) GetMap() map[string]string {
	return c.ConfigurationValues
}
