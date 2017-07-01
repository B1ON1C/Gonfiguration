# Gonfiguration
A text-based configuration manager written in Go for simple apps, fast and lightweight. Feel free to use it!

## Funcs
There are a lot of funcs that you can call:

```Go
- func New(path string, validateFunction func(line string) bool) (*Gonfiguration, error)
- func (g *Gonfiguration) LoadFromPath(path string) error
- func (g *Gonfiguration) Reload() error
- func (g *Gonfiguration) GetParamAsString(key string, def string) (string, error)
- func (g *Gonfiguration) GetParamAsStringArray(key, def, sep string) ([]string, error)
- func (g *Gonfiguration) GetParamAsInt(key string, def int) (int, error)
- func (g *Gonfiguration) GetParamAsIntArray(key, def, sep string) ([]int, error)
- func (g *Gonfiguration) GetParamAsBool(key string, def bool) (bool, error)
- func (g *Gonfiguration) AddNew(key, value string) error
- func (g *Gonfiguration) Update(key, value string) error
- func (g *Gonfiguration) Delete(key string) error
- func (g *Gonfiguration) Clear() error
- func (g *Gonfiguration) Contains(key string) (bool, error)
- func (g *Gonfiguration) Len() (int, error)
- func (g *Gonfiguration) Map() (*map[string]string, error)
- func (g *Gonfiguration) Dispose() 
```

## Usage
You only need to create a variable and call the "New" method, it will return a pointer with the initialized gonfiguration object. Now you can use all the functions of the object. We assume that we have the following configuration file:
```Properties
# MySQL (This is a comment, gonfiguration ignore this lines.)
mysql.host=localhost
mysql.port=3306
mysql.allowed.address=127.0.0.1,localhost,10.0.0.1
mysql.allowed.ports=3306,1232,1111,9000

# TCP (This is a comment, gonfiguration ignore this lines.)
tcp.ip=127.0.0.1
tcp.port=30000

#Variables (This is a comment, gonfiguration ignore this lines.)
welcome.message=You can use spaces, =, @ or what you want here.
chat.enabled=1
```
We want to read it with Gonfiguration, so we should use this Go code:

```Go
// create instance
Gonfiguration, error := gonfiguration.New("./config.ini", nil)

// Here we can read values, has example:
mysqlhost, _ := Gonfiguration.GetParamAsString("mysql.host", "DEFAULT") // localhost
mysqladdress, _ := Gonfiguration.GetParamAsStringArray("mysql.allowed.address", "DEFAULT1,DEFAULT2", ",") // [127.0.0.1 localhost 10.0.0.1]
mysqlport, _ := Gonfiguration.GetParamAsInt("mysql.port", 9999) // 3306
mysqlallowedports, _ := Gonfiguration.GetParamAsIntArray("mysql.allowed.ports", "DEFAULT21,DEFAULT2", ",") // [3306 1232 1111 9000]
chatenabled, _ := Gonfiguration.GetParamAsBool("chat.allowed", true) // True

// Tips
fmt.Println(Gonfiguration.Len()) // 8
Gonfiguration.Clear()
fmt.Println(Gonfiguration.Len()) // 0
Gonfiguration.Reload()
fmt.Println(Gonfiguration.Len()) // 8
containsValue, _ := Gonfiguration.Contains("chat.asymetric")
fmt.Println(containsValue) // false 
Gonfiguration.AddNew("chat.asymetric", "true") // ADD
Gonfiguration.AddNew("chat.asymetric", "true") // ERROR, DUPLICATED KEY
Gonfiguration.Update("chat.asymetric", "false") // UPDATE
Gonfiguration.Delete("chat.asymetric") // DELETE
Gonfiguration.Delete("chat.asymetric") // ERROR, KEY NOT FOUND

// If we aren't going to use the gonfigurationExample more, we should call dispose
gonfigurationExample.Dispose()
```

If you want to use another validation func, you can do it by passing it a function that returns true if the line belongs to the configuration, or pass nil if you want to use the default. The default validator is:
```Go
// isGonfigurationLine internal func for return true if line is correct
func (c *Gonfiguration) isGonfigurationLine(line string) bool {
    if c.ValidateConfigLine != nil {
        return c.ValidateConfigLine(line)
    }
    return !strings.HasPrefix(line, "#") && line != "" && strings.Contains(line, "=")
}
```

### Note
This package is not thread-safe.
