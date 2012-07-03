package content

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"html/template"
)

type Dir struct {
	base, rel string
}

func (d Dir) abs() string { return filepath.Join(d.base, d.rel) }

func (d Dir) join(subdir string) Dir {
	return Dir{d.base, filepath.Join(d.rel, subdir)}
}

func (d Dir) file(filename string) string {
	return filepath.Join(d.abs(), filename)
}

// Common contains all fields that are common to content items
//
type CommonData struct {
	dir Dir
	Title string
	Doc   struct {
		Rst  string
		Html template.HTML
	}
}

func (data *CommonData) Id() string { return toID(data.dir.rel) }

func (data *CommonData) absdir() string { return data.dir.abs() }

func (data *CommonData) Data() *CommonData { return data }

func (data *CommonData) read(dir Dir) {
	data.dir = dir
	_, last := filepath.Split(data.dir.rel)
	data.Title = removeOrder(last)
	dochtml := filepath.Join(dir.abs(), "doc.html")
	if raw, err := ioutil.ReadFile(dochtml); err == nil {
		data.Doc.Html = template.HTML(string(raw))
	}
	docrst := filepath.Join(dir.abs(), "doc.rst")
	if raw, err := ioutil.ReadFile(docrst); err == nil {
		data.Doc.Rst = string(raw)
	}
}

// Item

type ItemGroup interface {
	Get(id string) Item
	EachChild(fn func(id string, item Item))
}

type Item interface {
	ItemGroup
	Data() *CommonData
}

// Concept

type Concept struct {
	CommonData
	Depends []string
	VideoID string
}

func (c *Concept) Get(id string) Item { 
	return nil 
}

func (c *Concept) EachChild(fn func(string, Item)) {}

func (c *Concept) read(dir Dir) *Concept {
	c.CommonData.read(dir)
	if vid, err := ioutil.ReadFile(dir.file("video")); err == nil {
		c.VideoID = strings.Replace(string(vid), "\n", "", -1)
	}
	if deps, err := ioutil.ReadFile(dir.file("depends")); err == nil {
		for _, dep := range strings.Split(string(deps), "\n") {
			if dep != "" {
				c.Depends = append(c.Depends, toID(dep))
			}
		}
	}
	return c
}

// Group

type groupItem struct{ 
	Id string 
	Item Item 
}

type Group struct {
	ItemMap map[string]int
	Items []groupItem
}

func (g *Group) Add(id string, item Item) {
	if g.ItemMap == nil {
		g.ItemMap = make(map[string]int)
	}
	g.ItemMap[id] = len(g.Items)
	g.Items = append(g.Items, groupItem{id, item})
}

func (g Group) Get(id string) Item {
	if i, ok := g.ItemMap[id]; ok {
		return g.Items[i].Item
	}
	return nil
}

func (g Group) EachChild(fn func(string, Item)) {
	for _, git := range g.Items {
		fn(git.Id, git.Item)
	}
}

// Topic

type Topic struct {
	CommonData
	Group
	// Map of concepts?
}

func (t *Topic) read(dir Dir) *Topic {
	t.CommonData.read(dir)
	eachSubDir(dir.abs(), func(subdir string) {
		t.Add(toID(subdir), new(Concept).read(dir.join(subdir)))
	})
	return t
}

// Course

type Course struct {
	CommonData
	Group
	basedir string
}

func (c *Course) read(dir Dir) *Course {
	c.CommonData.read(dir)
	eachSubDir(dir.abs(), func(subdir string) {
		c.Add(toID(subdir), new(Topic).read(dir.join(subdir)))
	})
	return c
}

// CourseList

var courseList Group

func Read() {
	path := os.Getenv("ACADEMIO_PATH")
	if path == "" {
		log.Fatalf("Empty ACADEMIO_PATH")
	}
	for _, root := range filepath.SplitList(path) {
		err := eachSubDir(root, func(subdir string) {
			courseList.Add(toID(subdir), new(Course).read(Dir{root, subdir}))
		})
		if err != nil {
			log.Printf("ReadContent: Cannot read '%s': %s", root, err)
		}
	}
}

func Get(id string) (item Item) {
	var g ItemGroup = courseList
	for {
		i := strings.Index(id, ".")
		if i == -1 {
			return g.Get(id)
		}
		g = g.Get(id[:i])
		if g == nil {
			return nil
		}
		id = id[i+1:]
	}
	panic("unreachable")
	return nil
}

func Type(item Item) string {
	switch item.(type) {
	case *Course: return "course"
	case *Topic: return "topic"
	case *Concept: return "concept"
	}
	return ""
}

func Show() {
	courseList.EachChild(func(id string, course Item) {
		fmt.Printf("%s %s\n", id, course.Data().Title)
		course.EachChild(func (id string, topic Item) {
			fmt.Printf("   %s \"%s\"\n", id, topic.Data().Title)
			topic.EachChild(func (id string, concept Item) {
				fmt.Printf("      %s \"%s\"\n", id, concept.Data().Title)
			})
		})
		fmt.Println()
	})
}
