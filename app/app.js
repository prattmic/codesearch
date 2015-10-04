(function() {
  var app = angular.module('search', ['prettify']);

  app.controller('MainController', ['$http', function($http) {
    var main = this;
    main.file = 'github.com/prattmic/codesearch/cmd/search/search.go';
    main.fileContents = "";

    $http.get('/src/' + main.file)
      .success(function(response) {
        main.fileContents = response;
      })
      .error(function(response) {
        main.fileContents = 'Unable to load ' + main.file + ': ' + response;
      });
  }]);
})();
