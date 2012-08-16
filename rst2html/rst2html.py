#!/usr/bin/env python

import sys, os
import docutils.core
from docutils import nodes
from docutils.parsers.rst.roles import register_generic_role
from docutils.writers import html4css1

### Docutils

# inline directive :concept:`Subject/Topic/Concept`
# - takes an ID and generates the apropriate link.

class concept(nodes.TextElement): pass

register_generic_role('concept', concept)

## Highlighting in docutils (using pygments)

def normalize(cid):
   
   

class HTMLTranslator(html4css1.HTMLTranslator):
   def __init__(self, document):
      html4css1.HTMLTranslator.__init__(self, document)

   def visit_concept(self, node):
      ID = node.astext()
      url = ID
      name = ID.split('/')[-1]
      self.body.append('<a href="/%s">%s</a>' % (url, name))
      raise nodes.SkipNode

class HTMLWriter(html4css1.Writer):
   def __init__(self):
      html4css1.Writer.__init__(self)
      self.translator_class = HTMLTranslator

if len(sys.argv) != 2:
   print "usage: rst2html.py <filename>"
   sys.exit(0)

def read_file(filename, utf=False):
   try:
      if utf:
         F = codecs.open(filename, "r", "utf8")
      else:
         F = open(filename, "r")
      text = F.read()
      F.close()
      return text
   except IOError:
      return None

filename = sys.argv[1]
text = read_file(filename)
dic = docutils.core.publish_parts(text, writer = HTMLWriter())
print dic['body']
