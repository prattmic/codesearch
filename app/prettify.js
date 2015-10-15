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

          // Enable line numbers.
          var line = true;
          if (attrs.line) {
            // Starting line number.
            line = parseInt(attrs.line, 10);
          }

          // <prettify> -> <pre> -> <span>
          formatted = element.children().children();
          while (formatted.length && !formatted.hasClass("prettify-formatted")) {
            formatted = formatted.next();
          }
          if (!formatted.length) {
              console.error("Unable to find .prettify-formatted");
              return;
          }
          formatted.html(prettyPrintOne(html, null, line));
        });
      }
    };
  });
})();
