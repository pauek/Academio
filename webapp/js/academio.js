
var academio = {};

academio.updateMap = function () {
   if ($('#map').length == 0) {
      return;
   }
   
   var xmin, ymin;
   var xcurr = 0, ycurr = 0;
   var items = [];

   // Get positions
   $('#map .concept').each(function () {
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
      if (typeof ymin === "undefined" || it.y < ymin) {
         ymin = it.y;
      }
   }

   // Compute px, py
   $('#map .concept').each(function (i) {
      var it = items[i];
      it.px = (it.x - xmin) * 60;
      it.py = (it.y - ymin) * 60;
   });

   // Paint links
   var c = new fabric.Canvas('c', { 
      backgroundColor: "#eee",
      selection: false 
   });
   c.setHeight($('#map').height());
   c.setWidth($('#map').width());
   for (var i = 0; i < items.length; i++) {
      var deps = items[i].deps;
      for (var j = 0; j < deps.length; j++) {
         var k = deps[j];
         var coords = [items[i].px + 15, items[i].py + 15,
                       items[k].px + 15, items[k].py + 15];
         var line = new fabric.Line(coords, {
            fill: "rgb(200, 200, 200)",
            strokeWidth: 3,
            selectable: false
         });
         c.add(line);
      }
   }
   for (var i = 0; i < items.length; i++) {
      var rect = new fabric.Rect({
         left: items[i].px + 15, 
         top: items[i].py + 15, 
         width: 20,
         height: 20,
         fill: "#fff",
         stroke: "#555",
         strokeWidth: 2,
         selectable: false,
      });
      c.add(rect);
   }
   for (var i = 0; i < items.length; i++) {
      var text = new fabric.Text("" + (i + 1), {
         left: items[i].px + 15,
         top: items[i].py + 15,
         fontFamily: "Open Sans",
         fontSize: 12,
         stroke: "#555",
         selectable: false,
      });
      c.add(text);
   }
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