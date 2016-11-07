package i18n

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type (
	language map[string]interface{}

	// I18N provides methods to set and get the language depending on context.
	I18N struct {
		language string
	}
)

var data map[string]language
var debug bool

// Load scans for JSON files inside the provided pathDir and initialize the language map.
func Load(pathDir string) error {
	data = make(map[string]language)

	werr := filepath.Walk(pathDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		extension := filepath.Ext(path)
		if extension != ".json" {
			return nil
		}

		file, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		name := path[len(pathDir):]
		name = filepath.ToSlash(name)
		name = strings.TrimPrefix(name, "/")
		name = strings.TrimSuffix(name, extension)

		var l map[string]interface{}
		if err := json.Unmarshal(file, &l); err != nil {
			fmt.Println("error reading json data", err.Error())
		}

		data[name] = l

		return nil
	})

	if werr != nil {
		return werr
	}

	return nil
}

// Debug print missing translation strings when enabled.
func Debug(value bool) {
	debug = value
}

// New return a new I18N structure.
func New(language string) *I18N {
	return &I18N{language: language}
}

// Print translate a string, if args are passed they are parsed using Sprintf
func (i *I18N) Print(str string, args ...interface{}) string {
	// If language or translation string is not found, return the original string
	if data[i.language] == nil || data[i.language][str] == nil {
		if debug {
			fmt.Printf("WARNING: missing translation: [%s] %s\n", i.language, str)
		}

		if len(args) > 0 && strings.Contains(str, "%") {
			return fmt.Sprintf(str, args...)
		}

		return str
	}

	if len(args) > 0 && strings.Contains(str, "%") {
		return fmt.Sprintf(data[i.language][str].(string), args...)
	}

	return data[i.language][str].(string)
}

// Plural return
func (i *I18N) Plural(value int, zero string, one string, many string) string {
	if value <= 0 {
		return i.Print(zero, value)
	}

	if value == 1 {
		return i.Print(one, value)
	}

	if value > 1 {
		return i.Print(many, value)
	}

	return ""
}

// SetLang set the language.
func (i *I18N) SetLang(language string) {
	i.language = language
}

// GetLang get the current language.
func (i *I18N) GetLang() string {
	return i.language
}
