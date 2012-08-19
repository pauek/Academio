
var academio = {};

academio.updateMap = function () {
   if ($('#map').length == 0) {
      return;
   }
   
   var DIST = 60, WIDTH = 25, HEIGHT = 25, MARGIN = 20;

   var xmin, ymin, xmax, ymax;
   var xcurr = 0, ycurr = 0;
   var items = [];

   // Get positions
   $('.topic .list .concept').each(function () {
      var x = parseInt($(this).attr('x'));
      var y = parseInt($(this).attr('y'));
      var deps = JSON.parse($(this).attr('deps'));
      items.push({ x: x, y: y, deps: deps });
   });

   // Assign positions if necessary
   for (var i = 0; i < items.length; i++) {
      var it = items[i];
      if (it.x == -1 && it.y == -1) {
         it.x = xcurr, it.y = ycurr;
         xcurr += 1;
         if (xcurr >= 5) {
            xcurr = 0, ycurr++;
         }
      }
   }

   // Compute minimum of x and y
   for (var i = 0; i < items.length; i++) {
      var it = items[i];
      if (typeof xmin === "undefined" || it.x < xmin) {
         xmin = it.x;
      } 
      if (typeof xmax === "undefined" || it.x > xmax) {
         xmax = it.x;
      }
      if (typeof ymin === "undefined" || it.y < ymin) {
         ymin = it.y;
      }
      if (typeof ymax === "undefined" || it.y > ymax) {
         ymax = it.y;
      }
   }

   // Paint links
   var c = new fabric.Canvas('c', { 
      backgroundColor: "#f3f3f5",
      selection: false 
   });
   var width = $('#map').width();
   var height = $('#map').height();
   c.setWidth(width);
   c.setHeight(height);

   var scale, xoffset, yoffset;
   var xtotal = WIDTH + (xmax - xmin) * DIST + 2 * MARGIN;
   var ytotal = HEIGHT + (ymax - ymin) * DIST + 2 * MARGIN;
   var prop = height / width, pr = ytotal / xtotal;
   if (pr < prop) { // even negative values
      scale = width / xtotal;
      xoffset = (width - xtotal * scale) / 2.0;
      yoffset = 0;
   } else {
      scale = height / ytotal;
      xoffset = (width - xtotal * scale) / 2.0;
      yoffset = (height - ytotal * scale) / 2.0;
   }
   function xmap(x) {
      return (WIDTH / 2.0 + (x - xmin) * DIST) * scale + xoffset + MARGIN * scale;
   }
   function ymap(y) {
      return (HEIGHT / 2.0 + (y - ymin) * DIST) * scale + yoffset + MARGIN * scale;
   }

   for (var i = 0; i < items.length; i++) {
      var deps = items[i].deps;
      for (var j = 0; j < deps.length; j++) {
         var k = deps[j];
         var coords = [xmap(items[i].x), ymap(items[i].y),
                       xmap(items[k].x), ymap(items[k].y)];
         var line = new fabric.Line(coords, {
            fill: "rgb(200, 200, 200)",
            strokeWidth: 3 * scale,
            selectable: false
         });
         c.add(line);
      }
   }
   fabric.loadSVGFromURL('/img/item.svg', function(objects, options) {
      // rectangles
      var svg = fabric.util.groupSVGElements(objects, options);
      for (var i = 0; i < items.length; i++) {
         var item = svg.clone();
         item.set({
            left: xmap(items[i].x),
            top: ymap(items[i].y),
            selectable: false,
         });
         item.scaleToWidth(WIDTH * scale);
         c.add(item);
      }
      // text on top
      for (var i = 0; i < items.length; i++) {
         var text = new fabric.Text("" + (i + 1), {
            left: xmap(items[i].x),
            top: ymap(items[i].y),
            fontFamily: "Open Sans",
            fontSize: HEIGHT / 1.8 * scale,
            stroke: "#29393d",
            selectable: false,
         });
         c.add(text);
      }
   });
}

academio.showVideo = function (ev) {
   var videoid = $('#vid').attr('videoid');
   $('#vid img').replaceWith('<iframe class="youtube-player" ' +
                             'width="950" ' +
                             'height="570" ' +
                             'frameborder="0" ' +
                             'src="http://www.youtube.com/embed/' + 
                             videoid + 
                             '?autoplay=1&rel=0&wmode=transparent" ' +
                             'allowfullscreen>' +
                             '</iframe>');
}