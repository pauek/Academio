
var academio = {};

academio.updateMap = function () {
   var ctx;
   var xmin, ymin;
   var xcurr = 0, ycurr = 0;
   var items = [];

   // Get positions
   $('.map .concept').each(function () {
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

   // Move each element
   $('.map .concept').each(function (i) {
      var it = items[i];
      it.px = (it.x - xmin) * 60;
      it.py = (it.y - ymin) * 60;
      $(this).css({ position: "absolute", left: it.px, top: it.py });
   });

   // Paint links
   if (typeof ctx === "undefined") {
      ctx = $('#deps')[0].getContext('2d');
   }
   ctx.clearRect(0, 0, 500, 500);
   for (var i = 0; i < items.length; i++) {
      var deps = items[i].deps;
      for (var j = 0; j < deps.length; j++) {
         var k = deps[j];
         ctx.beginPath();
         // porqué +15??
         ctx.moveTo(items[i].px + 17, items[i].py + 17);
         ctx.lineTo(items[k].px + 17, items[k].py + 17);
         ctx.closePath();
         ctx.strokeStyle = "rgba(200, 200, 200, .3)";
         ctx.lineWidth = 4;
         ctx.stroke();
      }
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