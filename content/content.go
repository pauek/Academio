package content

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var roots []string

func init() {
	// Get roots
	pathlist := os.Getenv("ACADEMIO_PATH")
	if pathlist == "" {
		log.Fatalf("Empty ACADEMIO_PATH")
	}
	roots = filepath.SplitList(pathlist)
}

// Dir

type Dir struct {
	root, rel string
}

func (d Dir) abs() string { return filepath.Join(d.root, d.rel) }

func (d Dir) join(subdir string) Dir {
	return Dir{d.root, filepath.Join(d.rel, subdir)}
}

func (d Dir) file(filename string) string {
	return filepath.Join(d.abs(), filename)
}

// Common contains all fields that are common to content items
//
type CommonData struct {
	dir   Dir
	Title string
	Doc   struct {
		Rst  string
		Html template.HTML
	}
}

func (data *CommonData) Id() string { return toID(data.dir.rel) }

func (data *CommonData) Path() (ids []string) {
	acum := ""
	parts := strings.Split(data.Id(), ".")
	for i, part := range parts {
		if i < len(parts)-1 {
			acum += "." + part
			ids = append(ids, acum[1:])
		}
	}
	return
}

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

type SubItem struct {
	Id, Title, dir string
}

type ItemGroup interface {
	Children() []SubItem
}

type Item interface {
	ItemGroup
	Data() *CommonData
	Type() string
	read(dir Dir) Item
}

// Concept

type Concept struct {
	CommonData
	Depends []string
	VideoID string
}

func (c *Concept) Children() []SubItem {
	return nil
}

func (c *Concept) Type() string { return "concept" }

func (c *Concept) read(dir Dir) Item {
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

type Group struct {
	ItemMap map[string]int
	Items   []SubItem
}

func (g *Group) Add(subitem SubItem) {
	if g.ItemMap == nil {
		g.ItemMap = make(map[string]int)
	}
	g.ItemMap[subitem.Id] = len(g.Items)
	g.Items = append(g.Items, subitem)
}

func (g *Group) Children() []SubItem {
	return g.Items
}

func (g *Group) read(dir Dir) {
	eachSubDir(dir.abs(), func(subdir string) {
		d := dir.join(subdir)
		g.Add(SubItem{
			Id:    toID(d.rel),
			Title: removeOrder(subdir),
			dir:   d.abs(),
		})
	})
}

// Topic

type Topic struct {
	CommonData
	Group
	// Map of concepts?
}

func (t *Topic) Type() string { return "topic" }

func (t *Topic) read(dir Dir) Item {
	t.CommonData.read(dir)
	t.Group.read(dir)
	return t
}

// Course

type Course struct {
	CommonData
	Group
	basedir string
}

func (c *Course) Type() string { return "course" }

func (c *Course) read(dir Dir) Item {
	c.CommonData.read(dir)
	c.Group.read(dir)
	return c
}

// Interface

func Get(id string) (item Item) {
	dir := toDir(id)
	if dir.root == "" {
		return nil
	}
	switch numLevels(dir.rel) {
	case 1:
		item = new(Course)
	case 2:
		item = new(Topic)
	case 3:
		item = new(Concept)
	default:
		return nil
	}
	return item.read(dir)
}

func Courses() *Group {
	g := new(Group)
	for _, root := range roots {
		g.read(Dir{root, ""})
	}
	return g
}