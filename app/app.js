(function() {
  var app = angular.module('search', ['prettify']);

  // A shared service describes the current filename.
  app.factory('currentFile', function() {
    return {
      name: 'github.com/prattmic/codesearch/cmd/search/search.go'
    };
  });

  app.directive('searchBox', function() {
    return {
      restrict: 'E',
      templateUrl: 'search.html',
      controller: ['$http', 'currentFile', function($http, currentFile) {
        var search = this;
        search.query = "";

        search.loadFile = function(main) {
          currentFile.name = search.query;
          search.query = "";
        };

        search.search = function(main) {
          $http.post('/search', search.query)
            .success(function(response) {
              console.log(response);
              if (response.length) {
                currentFile.name = response[0].Path;
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

  app.directive('displayFile', function() {
    return {
      restrict: 'E',
      templateUrl: 'file.html',
      controller: ['$scope', '$http', 'currentFile', function($scope, $http, currentFile) {
        var file = this;
        file.service = currentFile;
        file.contents = "";

        // Watch for file to change, updating contents when it does.
        $scope.$watch(function(){
          return file.service.name;
        }, function(name) {
          file.contents = "Loading...";

          $http.get('/src/' + name)
            .success(function(response) {
              file.contents = response;
            })
            .error(function(response) {
              file.contents = 'Unable to load ' + file.name + ': ' + response;
            });
        });
      }],
      controllerAs: 'file'
    }
  });
})();
