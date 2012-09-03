#!/usr/bin/env python

import sys, os, subprocess, codecs, re
import docutils.core
from docutils import nodes
from docutils.parsers.rst.roles import register_generic_role
from docutils.writers import html4css1

### Docutils

# inline directive :concept:`Subject/Topic/Concept`
# - takes an ID and generates the apropriate link.

class concept(nodes.TextElement): pass
class math(nodes.TextElement): pass

register_generic_role('concept', concept)
register_generic_role('math', math)

## Highlighting in docutils (using pygments)

def normalize(cid):
   p = subprocess.Popen(["dir2id", cid], stdout=subprocess.PIPE)
   output = ""
   for line in p.stdout:
      output += line
   return output[:-1]

class HTMLTranslator(html4css1.HTMLTranslator):
   def __init__(self, document):
      html4css1.HTMLTranslator.__init__(self, document)

   # hack to be able to select tables in CSS (1)
   def visit_table(self, node):
      self.body.append('<div class="table">')
      html4css1.HTMLTranslator.visit_table(self, node)

   # hack to be able to select tables in CSS (2)
   def depart_table(self, node):
      html4css1.HTMLTranslator.depart_table(self, node)
      self.body.append('</div>')

   def visit_image(self, node):
      global itemid
      uri = node.attributes['uri']
      path = 'blabla'
      node.attributes['uri'] = itemid + '/' + uri
      self.body.append('<div class="image">')
      html4css1.HTMLTranslator.visit_image(self, node)
      
   def depart_image(self, node):
      html4css1.HTMLTranslator.depart_image(self, node)
      self.body.append('</div>')

   def visit_concept(self, node):
      path = node.astext()
      url = normalize(path)
      name = path.split('/')[-1]
      self.body.append('<a ajx href="/%s">%s</a>' % (url, name))
      raise nodes.SkipNode # avoid depart_...
   
   def visit_math(self, node):
      print "'" + node.rawsource + "'"
      g = re.match(':math:`(.*)`', node.rawsource, flags=re.DOTALL)
      if g: 
         math = g.group(1)
         print math
         self.body.append('\\(%s\\)' % math)
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

PATH = os.getenv('ACADEMIO_PATH')
for path in PATH.split(':'):
   itempath = sys.argv[1]
   itemid = normalize(itempath)
   print itemid
   rst = path + '/' + itempath + '/doc.rst'
   dirr, base = os.path.split(rst)
   print dirr
   name, ext = os.path.splitext(base)
   html = os.path.join(dirr, name + '.html')
   print html
   text = read_file(rst)
   dic = docutils.core.publish_parts(text, writer = HTMLWriter())
   f = codecs.open(html, 'w', 'utf8')
   f.write(dic['body'])
