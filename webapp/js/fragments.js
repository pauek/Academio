
var fragments = {};

fragments.db = {};
if (window.localStorage) {
   fragments.db = window.localStorage;
}

fragments.log = function(message) {
   console.log(message ? "fragments: " + message : "");
}

fragments.replace = function(item) {
   if (item.Html === "") {
      fragments.log('hit("' + item.Id + '")');
      item.Html = fragments.db[item.Id];
   } else {
      fragments.log('replace("' + item.Id + '")');
   }
   var elem = $('div[fragment="' + item.Id + '"]');
   // save children
   var children = {};
   $(elem).find('div[fragment]').each(function () {
      var id = $(this).attr("fragment");
      children[id] = $(this).detach();
   });
   // change elem
   $(elem).replaceWith(item.Html);
   // put children back
   for (var id in children) {
      if (children.hasOwnProperty(id)) {
         $(elem)
            .find('div[fragment="' + id + '"]')
            .replaceWith(children[id]);
      }
   }
}

fragments.assemble = function(list, $elem) {
   var msg = "assemble({";
   $elem.html('<div fragment="' + list[0].Id + '"></div>');
   for (var i = 0; i < list.length; i++) {
      var item = list[i];
      msg += " " + item.Id;
      fragments.replace(item);
      fragments.db[item.Id] = item.Html;      
   }
   msg += " })";
   fragments.log(msg);
   window.scrollTo(0, 0); // Go to top
}

fragments.load = function(url) {
   var stamp = fragments.db[url];
   stamp = (stamp ? '"' + JSON.parse(stamp) + '"' : "null");
   fragments.log('load("' + url + '", ' + stamp + ')');
   $.ajax({
      url: url,
      headers: { "FragmentsSince": stamp },
      dataType:"json",
      success: function(page) {
         fragments.db[url] = JSON.stringify(page.Stamp); 
         document.title = page.Title;
         fragments.log("Message is '" + page.Message + '"');
         if (page.Message !== "") {
            $('#message').html(page.Message);
            $('#message').show('fast');
         } else {
            $('#message').hide();
         }
         fragments.assemble(page.Navbar, $('#navbar'));
         fragments.assemble(page.Body, $('#body'));
         fragments.replaceLinks();
         _gaq.push(['_trackPageview', url]); // Google Analytics
      },
   });
}

fragments.follow = function(link) {
   var url = $(link).attr("href")
   fragments.log("follow: " + url);
   fragments.load(url);
   history.pushState(url, null, url);
}

fragments.replaceLinks = function() {
   fragments.log("replaceLinks()");
   $('a[ajx]').click(function (ev) {
      ev.preventDefault();
      fragments.follow(this);
   });
   fragments._notify();
   fragments.log();
}

fragments._fns = []

fragments.ready = function(fn) {
   fragments._fns.push(fn);
}

fragments._notify = function() {
   var fns = fragments._fns;
   for (var i = 0; i < fns.length; i++) {
      fns[i]();
   }
}

// Bugs:
// - Chrome: Al clicar a la '/' manualmente (el icono de 'academ.io')
//           el botón 'back' muestra el texto AJAX... (en Chrome)

fragments._onpopstate = function(ev) {
   fragments.log("_onpopstate(" + JSON.stringify(ev.state) + ")")
   fragments.log();
   if (ev.state === null) {
      history.replaceState("reload", null, document.location.pathname);
   } else if (ev.state === "reload") {
      document.location.reload();
   } else {
      fragments.load(ev.state);
   }
}

$(document).ready(function () {
   history.replaceState("reload", null, document.location.pathname);
   fragments.log("ready()");
   if (history && history.pushState) {
      fragments.replaceLinks();
      onpopstate = fragments._onpopstate;
   }
})
