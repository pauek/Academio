package content

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`(^|/)[0-9]+\. ?`)

// Remove front-numbers from directory names:
//
//   "1. Blah/2. Blih" => "Blah/Blih"
//   "004. The Title" => "The Title"
//
func removeOrder(dir string) string {
	s := re.ReplaceAllString(dir, "/")
	if len(s) > 0 && s[0] == '/' {
		s = s[1:]
	}
	return s
}

var replacer *strings.Replacer

func init() {
	bef := "àèìòùáéíóúÀÈÌÒÙÁÉÍÓÚäëïöüÄËÏÖÜñÑ +.:-()"
	aft := "aeiouaeiouAEIOUAEIOUaeiouAEIOUnN p     "
	R := []string{}
	for i, b := range strings.Split(bef, "") {
		R = append(R, b)
		R = append(R, aft[i:i+1])
	}
	replacer = strings.NewReplacer(R...)
}

// Convert a directory name to an ID
// 
func ToID(dir string) (id string) {
	id = removeOrder(dir)                  // remove "^[0-9]+. "
	id = replacer.Replace(id)              // remove accents + map chars
	id = strings.Title(id)                 // Make A Title
	id = strings.Replace(id, " ", "", -1)  // remove spaces
	id = strings.Replace(id, "/", ".", -1) // '/' -> '.'
	return id
}

// Determine the ID of the parent
//
func parentID(id string) (pid string) {
	if len(id) > 0 {
		parts := strings.Split(id, ".")
		pid = strings.Join(parts[:len(parts)-1], ".")
	}
	return
}

// Find a dir from the roots that translates to ID
// 
func ToDir(ID string) (dir Dir) {
	for _, root := range roots {
		found := ""
		walkDirs(root, func(reldir string, level int) {
			if ToID(reldir) == ID && found == "" {
				found = reldir
			}
		})
		if found != "" {
			return Dir{root, found}
		}
	}
	return
}

// Remove the root from an absolute path
//
func removeRoot(abspath string) string {
	for _, root := range roots {
		if root == abspath {
			return ""
		}
		if rel, err := filepath.Rel(root, abspath); err == nil {
			if rel[0] != '.' || rel[1] != '.' {
				return rel
			}
		}
	}
	panic("unreachable")
	return ""
}

// Determine the number of levels in a relative dir
//
func numLevels(reldir string) (lv int) {
	if len(reldir) > 0 {
		lv = 1 + strings.Count(reldir, "/")
	}
	return
}

// Determine the number of levels in a relative dir
//
func numLevelsID(id string) (lv int) {
	if len(id) > 0 {
		lv = 1 + strings.Count(id, ".")
	}
	return
}

// Do something for each subdirectory
//
func eachSubDir(dir string, fn func(dir string)) error {
	fileinfo, err := ioutil.ReadDir(dir) // sorted by name
	if err != nil {
		return err
	}
	for _, linfo := range fileinfo {
		// since ReadDir uses Lstat, we have to 
		// repeat with os.Stat to follow the symlinks
		info, _ := os.Stat(filepath.Join(dir, linfo.Name())) // sure ignore err?		
		if info.IsDir() && info.Name()[0] != '.' && info.Name()[0] != '_' {
			fn(info.Name())
		}
	}
	return nil
}

// Walk a 3-level hierarchy of directories
//
func walkDirs(path string, fn func(dir string, level int)) {
	// need to declare since it's recursive
	var walk func(dir string, level int)

	walk = func(dir string, level int) {
		if level > 0 {
			fn(dir, level)
		}
		eachSubDir(filepath.Join(path, dir), func(subdir string) {
			reldir := filepath.Join(dir, subdir)
			if level < 3 {
				walk(reldir, level+1)
			}
		})
	}

	walk("", 0)
}
