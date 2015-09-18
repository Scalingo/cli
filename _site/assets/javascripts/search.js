(function() {
  $(function() {
    var htmlEscape, spinner;
    spinner = {
      nbRequests: 0,
      isSpinning: false,
      start: function() {
        if (!this.isSpinning) {
          $("#search-icon").html("<i class=\"fa fa-spinner fa-spin\"></i>");
          return this.isSpinning = true;
        }
      },
      stop: function() {
        if (this.nbRequests === 0) {
          $("#search-icon").html("<i class=\"fa fa-search\"></i>");
          return this.isSpinning = false;
        }
      }
    };
    htmlEscape = function(str) {
      return String(str).replace(/&/g, '&amp;').replace(/"/g, '&quot;').replace(/'/g, '&#39;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
    };
    $('#st-search-input').swiftypeSearch({
      engineKey: 'rQW8Hh49XhMVApNATHZL',
      resultContainingElement: "#search-results",
      renderFunction: function(document_type, item) {
        var hl, title;
        title = htmlEscape(item.title).split("-", 2)[1];
        hl = item.highlight.body || item.highlight.sections || item.highlight.title;
        return '<div class="st-result"> <h3 class="title no-top"> <a href="' + item.url + '" class="st-search-result-link">' + title + '</a> </h3> <p>' + hl + '</p> </div>';
      },
      preRenderFunction: function(data) {
        return $("#search-modal").modal('show');
      }
    });
    return $('#st-search-input').swiftype({
      engineKey: 'rQW8Hh49XhMVApNATHZL',
      onRemoteComplete: function(data) {
        spinner.nbRequests -= 1;
        return spinner.stop();
      },
      beforeRemoteCall: function() {
        spinner.nbRequests += 1;
        return spinner.start();
      }
    });
  });

}).call(this);
