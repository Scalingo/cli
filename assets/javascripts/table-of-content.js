define(['jquery'], function ($) {

  function slugify(s) {
    return s
      .toLowerCase()
      .replace(/[^\w ]+/g,'')
      .replace(/ +/g,'-')
  }

  if (window.location.pathname !== '/404.html' && document.title.indexOf('404 Not found') == -1) {
    var sidebar_nav = $('.sidebar-nav')
    var headers = $('.content').find(':header:not(h1)')
    var ul = $('<ul>').addClass('list-unstyled')

    headers.each(function (index, el) {
      // make id from cotent
      var innerText = $(el).text()
      var id = slugify(innerText)

      // inject id
      el.id = id

      var a = $('<a>')
        .attr('href','#' + id)
        .text(innerText)

      var li = $('<li>')
        .append(a)

      ul.append(li)
    })

    sidebar_nav.prepend(ul)
    sidebar_nav.prepend('<strong>Table of content<strong>')
  }

})
