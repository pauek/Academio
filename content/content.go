package content

import (
	"fmt"
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
	Root, Rel string
}

func (d Dir) abs() string { return filepath.Join(d.Root, d.Rel) }

func (d Dir) join(subdir string) Dir {
	return Dir{d.Root, filepath.Join(d.Rel, subdir)}
}

func (d Dir) File(filename string) string {
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

func (data *CommonData) Id() string { return ToID(data.dir.Rel) }

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
	_, last := filepath.Split(data.dir.Rel)
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
	if vid, err := ioutil.ReadFile(dir.File("video")); err == nil {
		c.VideoID = strings.Replace(string(vid), "\n", "", -1)
	}
	if deps, err := ioutil.ReadFile(dir.File("depends")); err == nil {
		for _, dep := range strings.Split(string(deps), "\n") {
			if dep != "" {
				c.Depends = append(c.Depends, ToID(dep))
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
		subitem := SubItem{
			Id:    ToID(d.Rel),
			Title: removeOrder(subdir),
			dir:   d.abs(),
		}
		g.Add(subitem)
	})
}

// Topic

type XY struct{ X, Y int }

type Topic struct {
	CommonData
	Group
	Coords []XY
	Deps [][]int
	// Map of concepts?
}

type Info struct {
	Index int
	Item SubItem
	Coords XY
	Deps string
}

func (t *Topic) ChildrenInfo() (info []Info) {
	for i := range t.Items {
		deps := "["
		for i, d := range t.Deps[i] {
			if i > 0 {
				deps += ", "
			}
			deps += fmt.Sprintf("%d", d)
		}
		deps += "]";

		info = append(info, Info{
			Index: i + 1,
			Item: t.Items[i],
			Coords: t.Coords[i],
			Deps: deps,
		})
	}
	return
}

func (t *Topic) Type() string { return "topic" }

func (t *Topic) index(id string) int {
	for i, subitem := range t.Items {
		if subitem.Id == id {
			return i
		}
	}
	return -1
}

func (t *Topic) read(dir Dir) Item {
	t.CommonData.read(dir)
	t.Group.read(dir)
	t.Coords = make([]XY, len(t.Group.Items))
	t.Deps = make([][]int, len(t.Group.Items))
	for i, subitem := range t.Items {
		path := filepath.Join(subitem.dir, "xy")
		xy := XY{-1, -1}
		if file, err := os.Open(path); err == nil {
			fmt.Fscanf(file, "%d %d", &xy.X, &xy.Y)
			file.Close()
		}
		path = filepath.Join(subitem.dir, "depends")
		di := []int{}
		if depends, err := ioutil.ReadFile(path); err == nil {
			for _, dep := range strings.Split(string(depends), "\n") {
				// fmt.Printf("dep = '%s'\n", ToID(dep))
				if j := t.index(ToID(dep)); j != -1 {
					di = append(di, j)
				}
			}
		}
		t.Deps[i] = di
		t.Coords[i] = xy
	}
	return t
}

// Course

type Course struct {
	CommonData
	Group
	Photo string // path for the photo
}

func (c *Course) Type() string { return "course" }

func (c *Course) read(dir Dir) Item {
	c.CommonData.read(dir)
	c.Group.read(dir)
	if _, err := os.Stat(dir.File("photo.png")); err == nil {
		c.Photo = dir.File("photo.png")
	}
	return c
}

// Interface

func Get(id string) (item Item) {
	dir := ToDir(id)
	if dir.Root == "" {
		return nil
	}
	switch numLevels(dir.Rel) {
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
