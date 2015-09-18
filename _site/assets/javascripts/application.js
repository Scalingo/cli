requirejs.config({
  waitSeconds: 200,
  paths: {
    'jquery': '//ajax.googleapis.com/ajax/libs/jquery/2.1.1/jquery.min'
  },
  shim: {
    'jquery.swiftype.autocomplete': ['jquery'],
    'jquery.swiftype.search': ['jquery'],
    'search': ['jquery', 'jquery.swiftype.autocomplete', 'jquery.swiftype.search'],
    'jquery.ba-hashchange.min': ['jquery'],
    'highlight': ['highlight.pack']
  }
})

requirejs([
  'table-of-content',
  'jquery.ba-hashchange.min',
  'highlight',
  'search'
  ])

