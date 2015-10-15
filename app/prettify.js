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
          raw = element.find("span");
          return raw.html();
        }, function(html){
          // ... then redo the formatted content.

          if (attrs.highlight) {
              html = html.replace(RegExp(attrs.highlight),
                                  "<span class=\"highlight\">$&</span>");
          }

          // Enable line numbers.
          var line = true;
          if (attrs.line) {
            // Starting line number.
            line = parseInt(attrs.line, 10);
          }

          formatted = element.find("span").next();
          formatted.html(prettyPrintOne(html, null, line));
        });
      }
    };
  });
})();
