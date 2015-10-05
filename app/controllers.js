(function() {
  var app = angular.module('searchControllers', []);

  app.directive('searchBox', function() {
    return {
      restrict: 'E',
      templateUrl: 'search.html',
      controller: ['$http', '$location', function($http, $location) {
        var search = this;
        search.query = "";

        search.search = function(main) {
          $http.post('/search', search.query)
            .success(function(response) {
              console.log(response);
              if (response.length) {
                // Redirect to the first result
                $location.path('/file/' + response[0].Path);
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

  app.controller('FileCtrl', ['$routeParams', '$http', function($routeParams, $http) {
    var file = this;
    file.name = $routeParams.name;
    file.contents = "Loading...";

    // Strip leading slash.
    if (file.name.length && file.name[0] === '/') {
      file.name = file.name.slice(1);
    }

    $http.get('/src/' + file.name)
      .success(function(response) {
        file.contents = response;
      })
      .error(function(response) {
        file.contents = 'Unable to load ' + file.name + ': ' + response;
      });
  }]);
})();
