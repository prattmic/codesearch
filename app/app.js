(function() {
  var app = angular.module('search', ['prettify']);

  app.controller('MainController', function() {
    this.file = "github.com/prattmic/codesearch/cmd/search/search.go";
  });
})();
