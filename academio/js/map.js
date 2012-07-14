var paper;
var gridSz = 80;
var rectSz = 40;

function getInfo(concept) {
   var elem = $(concept).find('.num');
   var num = parseInt(elem.text()) - 1;
   var x = parseInt(elem.attr('x'));
   var y = parseInt(elem.attr('y'));
   return {
      x: (x !== -1 ? x : num % 5),
      y: (y !== -1 ? y : Math.floor(num / 5)),
      deps: JSON.parse(elem.attr('deps')),
      num: num + 1,
      a: $(concept).find('a'),
   }
}

function drawItem(info, i) {
   var item = info[i];
   var r = paper.rect(item.x*gridSz, item.y*gridSz, rectSz, rectSz, 4);
   r.attr({ fill: "#ddd", stroke: "", cursor: "pointer" });
   var t = paper.text(item.x*gridSz + rectSz/2, item.y*gridSz + rectSz/2, item.num);
   t.attr({ font: "16px Open Sans", fill: "#666", cursor: "pointer" });
   r.click(function () {
      fragments.follow(item.a);
   });
}

function drawDeps(info, i) {
   var item = info[i];
   for (var j = 0; j < item.deps.length; j++) {
      var path = "M" + (item.x*gridSz + rectSz/2) + " " + (item.y*gridSz + rectSz/2);
      var k = item.deps[j];
      path += "L" + (info[k].x*gridSz + rectSz/2) + " " + (info[k].y*gridSz + rectSz/2);
      var p = paper.path(path);
      p.attr({ "stroke-width": 3, stroke: "#eee" });
   }
}

function drawMap() {
   if (paper === undefined) {
      paper = Raphael("map", 500, 500);
   }
   $('#map').append(paper.canvas);
   paper.clear();
   var info = [];
   $('.topic .concept').each(function () {
      info.push(getInfo(this));
   });
   $(info).each(function (i) { drawDeps(info, i); });
   $(info).each(function (i) { drawItem(info, i); });
}

fragments.onUpdate(drawMap);