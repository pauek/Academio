
package content

import (
	"io/ioutil"
	"regexp"
	"strings"
)

var first = regexp.MustCompile(`^[0-9]+. ?`)
var rest = regexp.MustCompile(`/[0-9]+. ?`)

func removeOrder(dir string) string {
	dir = first.ReplaceAllString(dir, "")
	dir = rest.ReplaceAllString(dir, "/")
	return dir
}

// Convert a directory name to an ID
// 
func toID(dir string) string {
	id := removeOrder(dir)

	// remove accents + map certain characters
	bef := "àèìòùáéíóúÀÈÌÒÙÁÉÍÓÚäëïöüÄËÏÖÜñÑ +.:-()"
	aft := "aeiouaeiouAEIOUAEIOUaeiouAEIOUnN p     "
	R := []string{}
	for i, b := range strings.Split(bef, "") {
		R = append(R, b)
		R = append(R, aft[i:i+1])
	}
	r := strings.NewReplacer(R...)
	id = r.Replace(id)

	id = strings.Title(id)                 // Make A Title
	id = strings.Replace(id, " ", "", -1)  // remove spaces
	id = strings.Replace(id, "/", ".", -1) // remove '/'
	return id
}

func eachSubDir(dir string, fn func(dir string)) error {
	fileinfo, err := ioutil.ReadDir(dir) // sorted by name
	if err != nil {
		return err
	}
	for _, info := range fileinfo {
		if info.IsDir() && info.Name()[0] != '.' {
			fn(info.Name())
		}
	}
	return nil
}
