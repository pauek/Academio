
package content

import (
	"strings"
)

type ContentType int

const (
	ErrorType = ContentType(iota)
	CourseType 
	TopicType
	ConceptType
)

func (c ContentItem) Html() string {
	return string(c.Doc.Html)
}

func Type(id Id) ContentType {
	switch strings.Count(string(id), ".") {
	case 0:
		return CourseType
	case 1: 
		return TopicType
	case 2:
		return ConceptType
	}
	return ErrorType
}

func GetConcept(id Id) *Concept {
	c, ok := concepts[id]
	if !ok {
		dir, okd := dirs[id]
		if !okd {
			return nil
		}
		c = conceptFromDir(dir)
		if c != nil {
			concepts[id] = c
		}
	}
	return c
}

func GetTopic(id Id) *Topic {
	return topics[id]
}

func GetCourse(id Id) *Course {
	return courses[id]
}

