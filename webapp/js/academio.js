
var academio = {};

academio.updateMap = function () {
   if ($('#map').length == 0) {
      return;
   }
   
   var DIST = 60, WIDTH = 25, HEIGHT = 25, MARGIN = 20;

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
   var xmin = items[0].x, ymin = items[0].y;
   var xmax = items[0].x, ymax = items[0].y;
   for (var i = 1; i < items.length; i++) {
      var it = items[i];
      if (it.x < xmin) { xmin = it.x; } 
      if (it.x > xmax) { xmax = it.x; } 
      if (it.y < ymin) { ymin = it.y; } 
      if (it.y > ymax) { ymax = it.y; } 
   }

   // Paint links
   var c = new fabric.Canvas('c', { 
      backgroundColor: "#f3f3f5",
      hoverCursor: 'pointer',
      selection: false,
   });
   var width = $('#map').width();
   var height = $('#map').height();
   c.setWidth(width);
   c.setHeight(height);

   function high(target, light) {
      var i = target._index;
      var div = $('.topic .list .concept')[i];
      var color = (light ? '#eee' : '#fff');
      $(div).css({ background: color });
      if (target != null) {
         target.getObjects()[0].setFill(color);
         c.renderAll();
      }
   }

   var over;
   c.on('mouse:move', function (opts) {
      var target = c.findTarget(opts.e, true);
      if (target) {
         if (over != target) {
            high(target, true);
            if (over != null) {
               high(over, false);
            }
            over = target;
         }
      } else if (over != null) {
         high(over, false);
         over = null;
      }
   });
   c.on('mouse:down', function (e) {
      if (over !== null) {
         var i = over._index;
         fragments.follow($('.topic .list .concept a')[i]);
      }
   });

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
      return (MARGIN + WIDTH / 2.0 + (x - xmin) * DIST) * scale + xoffset;
   }
   function ymap(y) {
      return (MARGIN + HEIGHT / 2.0 + (y - ymin) * DIST) * scale + yoffset;
   }

   // links
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
            selectable: false,
            left: xmap(items[i].x),
            top: ymap(items[i].y),
         });
         item.scaleToWidth(WIDTH * scale);
         var text = new fabric.Text("" + (i + 1), {
            selectable: false,
            left: xmap(items[i].x),
            top: ymap(items[i].y),
            fontFamily: "Open Sans",
            fontSize: HEIGHT / 1.8 * scale,
            fontWeight: "normal",
            fontStyle: "normal",
            strokeWidth: 1,
            fill: "#29303d",
         });
         var group = new fabric.Group([item, text], {
            _index: i,
            hasControls: false,
            hasBorders: false,
            lockMovementX: true,
            lockMovementY: true,
         });
         c.add(group);
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
