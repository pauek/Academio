package content

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Id string

type ContentItem struct {
	Title string
	Doc   struct {
		Rst  []byte
		Html []byte
	}
}

type Concept struct {
	ContentItem
	Depends []Id
	VideoID string
	// TODO: Exercises []string
}

type Topic struct {
	ContentItem
	Concepts []Id
	// TODO: Progress Map
}

type Course struct {
	ContentItem
	Topics []Id
}

// keep in what directory is a particular course/topic/concept
type RelDir struct {
	base, rel string
}

func (rdir RelDir) join() string {
	return filepath.Join(rdir.base, rdir.rel)
}

var dirs = make(map[Id]RelDir)

var concepts = make(map[Id]*Concept)
var topics = make(map[Id]*Topic)
var courses = make(map[Id]*Course)

var first = regexp.MustCompile(`^[0-9]+. ?`)
var rest = regexp.MustCompile(`/[0-9]+. ?`)

func removeOrder(dir string) string {
	dir = first.ReplaceAllString(dir, "")
	dir = rest.ReplaceAllString(dir, "/")
	return dir
}

// Convert a directory name to an ID
// 
func dirToID(dir string) Id {
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
	return Id(id)
}

func subDirs(dir string, fn func(dir string)) error {
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

func (item *ContentItem) read(dir string) {
	_, last := filepath.Split(dir)
	item.Title = removeOrder(last)
	dochtml := filepath.Join(dir, "doc.html")
	if data, err := ioutil.ReadFile(dochtml); err == nil {
		item.Doc.Html = data
	}
	docrst := filepath.Join(dir, "doc.rst")
	if data, err := ioutil.ReadFile(docrst); err == nil {
		item.Doc.Rst = data
	}
}

func conceptFromDir(rdir RelDir) (C *Concept) {
	dir := rdir.join()
	C = new(Concept)
	C.read(dir)
	vidfile := filepath.Join(dir, "video")
	if vid, err := ioutil.ReadFile(vidfile); err == nil {
		C.VideoID = strings.Replace(string(vid), "\n", "", -1)
	}
	depsfile := filepath.Join(dir, "depends")
	if deps, err := ioutil.ReadFile(depsfile); err == nil {
		for _, dep := range strings.Split(string(deps), "\n") {
			if dep != "" {
				C.Depends = append(C.Depends, dirToID(dep))
			}
		}
	}
	return C
}

func topicFromDir(rdir RelDir) (T *Topic) {
	dir := rdir.join()
	T = new(Topic)
	T.read(dir)
	subDirs(dir, func(subdir string) {
		id := dirToID(filepath.Join(rdir.rel, subdir))
		T.Concepts = append(T.Concepts, id)
	})
	return T
}

func courseFromDir(rdir RelDir) (C *Course) {
	dir := rdir.join()
	C = new(Course)
	C.read(dir)
	subDirs(dir, func(subdir string) {
		id := dirToID(filepath.Join(rdir.rel, subdir))
		C.Topics = append(C.Topics, id)
	})
	return C
}

func readDirectories(root string, reldir string, level int) {
	dir := filepath.Join(root, reldir)
	err := subDirs(dir, func (subdir string) {
		nextdir := filepath.Join(reldir, subdir)
		id := dirToID(nextdir)
		dirs[id] = RelDir{root, nextdir}
		if level < 2 {
			readDirectories(root, nextdir, level+1)
		}
	})
	if err != nil {
		log.Printf("readDirectories: %s", err)
	}
}

func readCourseOrTopic(id Id, reldir RelDir) {
	switch strings.Count(string(id), ".") {
	case 0:
		courses[id] = courseFromDir(reldir)
	case 1:
		topics[id] = topicFromDir(reldir)
	case 2:
		concepts[id] = conceptFromDir(reldir)
	}
}

func init() {
	path := os.Getenv("ACADEMIO_PATH")
	if path == "" {
		log.Fatalf("Empty ACADEMIO_PATH")
	}
	for _, dir := range filepath.SplitList(path) {
		readDirectories(dir, "", 0)
	}
	for id, rdir := range dirs {
		readCourseOrTopic(id, rdir)
	}

	// Debug info
	/*
		for id, rdir := range dirs {
			fmt.Printf("%#v %#v\n", id, rdir)
		}
		for id, course := range courses {
			fmt.Printf("%s: %#v %#v\n", id, course.Title, course.Topics)
		}
		for id, topic := range topics {
			fmt.Printf("%s: %#v %#v\n", id, topic.Title, topic.Concepts)
		}
	*/
	for id, c := range concepts {
		fmt.Printf("%s: %#v %#v %d\n", id, c.Title, c.Depends, len(c.Doc.Rst))
	}
}
