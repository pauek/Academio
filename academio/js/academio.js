
var academio = {};

academio.updateMap = function () {
   // Compute minimum x and y (+ deps)
   var xmin, ymin;
   var deps = [];
   $('.map .concept').each(function () {
      var x = parseInt($(this).attr('x'));
      var y = parseInt($(this).attr('y'));
      if (typeof xmin === "undefined" || x < xmin) {
         xmin = x;
      } 
      if (typeof ymin === "undefined" || y < ymin) {
         ymin = y;
      }
      deps.push(JSON.parse($(this).attr('deps')));
   });
   // Move each element
   var pos = [];
   $('.map .concept').each(function () {
      var x = parseInt($(this).attr('x'));
      var y = parseInt($(this).attr('y'));
      var p = { x: (x - xmin) * 60, y: (y - ymin) * 60 };
      $(this).css({ position: "absolute", left: p.x, top: p.y });
      pos.push(p);
   });
   // Paint links
   var ctx = $('#deps')[0].getContext('2d');
   ctx.clearRect(0, 0, 500, 500);
   for (var i = 0; i < deps.length; i++) {
      for (var j = 0; j < deps[i].length; j++) {
         var k = deps[i][j];
         ctx.beginPath();
         // porqué -5??
         ctx.moveTo(pos[i].x+15, pos[i].y+15);
         ctx.lineTo(pos[k].x+15, pos[k].y+15);
         ctx.closePath();
         ctx.strokeStyle = "rgba(200, 200, 200, .35)";
         ctx.lineWidth = 4;
         ctx.stroke();
      }
   }
}
