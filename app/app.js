(function() {
  var app = angular.module('search', ['ngRoute', 'searchControllers', 'prettify']);

  app.config(['$routeProvider', '$locationProvider', function($routeProvider, $locationProvider) {
    $routeProvider.when('/', {
      templateUrl: 'welcome.html',
    });

    $routeProvider.when('/file/:name*', {
      templateUrl: 'file.html',
      controller: 'FileCtrl',
      controllerAs: 'file',
    });

    $routeProvider.otherwise({
      redirectTo: '/'
    });

    // Utilize HTML5 URL rewriting.
    $locationProvider.html5Mode(true);
  }]);
})();
