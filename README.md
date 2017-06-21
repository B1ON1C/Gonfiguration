# Gonfiguration
A text-based configuration manager written in Go for simple apps, fast and lightweight.

## Usage
You only need to create a variable and call the LoadGonfiguration method, it will return a pointer with the initialized Gonfiguration object.

## Example of a Configuration file
This code can read a configuration file like the next example:
```
# Comment1
key=value
key2=value2

# Comment2
key3=value3
```
