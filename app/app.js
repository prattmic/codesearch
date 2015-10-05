(function() {
  var app = angular.module('search', ['prettify']);

  app.controller('MainController', ['$scope', '$http', function($scope, $http) {
    var main = this;
    main.file = 'github.com/prattmic/codesearch/cmd/search/search.go';
    main.fileContents = "";

    // Watch for file to change, updating contents when it does.
    $scope.$watch(function(){
      return main.file;
    }, function(file) {
      main.fileContents = "Loading...";

      $http.get('/src/' + file)
        .success(function(response) {
          main.fileContents = response;
        })
        .error(function(response) {
          main.fileContents = 'Unable to load ' + main.file + ': ' + response;
        });
    });
  }]);

  app.controller('SearchController', ['$http', function($http) {
    var search = this;
    search.query = "";

    search.loadFile = function(main) {
      main.file = search.query;
      search.query = "";
    };

    search.search = function() {
      $http.post('/search', search.query)
        .success(function(response) {
          console.log(response);
        })
        .error(function(response) {
          console.log(response);
        });
    };
  }]);
})();
