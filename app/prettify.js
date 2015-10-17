(function() {
  var app = angular.module('prettify', [ ]);

  app.directive('prettify', function() {
    return {
      restrict: 'E',
      transclude: true,
      templateUrl: "prettify.html",
      link: function($scope, element, attrs) {
        $scope.$watch(function() {
          // Watch for the raw content to change.
          raw = element.find(".prettify-raw");
          return raw.html();
        }, function(html){
          // ... then redo the formatted content.

          // Enable line numbers.
          var line = true;
          if (attrs.line) {
            // Starting line number.
            line = parseInt(attrs.line, 10);
          }

          formatted = element.find(".prettify-formatted");
          formatted.html(prettyPrintOne(html, null, line));
        });
      }
    };
  });
})();
