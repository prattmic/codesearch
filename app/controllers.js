(function() {
  var app = angular.module('searchControllers', []);

  app.directive('searchBox', function() {
    return {
      restrict: 'E',
      templateUrl: 'search-box.html',
      controller: ['$location', function($location) {
        var search = this;
        search.query = "";

        search.search = function(main) {
          $location.path('/search').search({'query': search.query});
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

  app.controller('SearchCtrl', ['$routeParams', '$http', function($routeParams, $http) {
    var search = this;
    search.query = $routeParams.query;
    search.total = 0;
    search.results = [];
    search.error = '';

    $http.post('/api/search', search.query)
      .success(function(response) {
        search.total = response.Total;
        search.results = response.Results;
      })
      .error(function(response) {
        search.results = [];
        search.error = 'Search error: ' + response;
      });
  }]);
})();
