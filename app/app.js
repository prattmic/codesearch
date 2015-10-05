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

  app.directive('searchBox', function() {
    return {
      restrict: 'E',
      templateUrl: 'search.html',
      controller: ['$http', function($http) {
        var search = this;
        search.query = "";

        search.loadFile = function(main) {
          main.file = search.query;
          search.query = "";
        };

        search.search = function(main) {
          $http.post('/search', search.query)
            .success(function(response) {
              console.log(response);
              if (response.length) {
                main.file = response[0].Path;
              }
            })
            .error(function(response) {
              console.log(response);
            });
        };
      }],
      controllerAs: 'search'
    };
  });
})();
