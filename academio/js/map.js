var paper = Raphael("map", 500, 500);

function drawItem() {
   var concept = this;
   var num = parseInt($(this).find('.num').text()) - 1;
   var i = num % 5, j = Math.floor(num / 5);
   var r = paper.rect(i*50, j*50, 30, 30, 4);
   r.attr({ 
      fill: "#ddd", 
      stroke: "",
      cursor: "pointer"
   });
   var t = paper.text(i*50+15, j*50+15, num + 1);
   t.attr({ 
      font: "14px Open Sans", 
      fill: "#666",
      cursor: "pointer"
   });
   r.click(function () {
      fragments.follow($(concept).find('a'));
   });
}

$(document).ready(function () {
   $('.topic .concept').each(drawItem);
});