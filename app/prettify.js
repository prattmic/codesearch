(function() {
	var app = angular.module('prettify', [ ]);

	app.directive('prettify', function() {
		return {
			restrict: 'E',
			transclude: true,
			templateUrl: "prettify.html",
			link: function($scope, element, attrs) {
				$scope.$watch(function() {
					raw = element.find("span");
					return raw.html();
				}, function(html){
					formatted = element.find("span").next();

					formatted.html(prettyPrintOne(html));
				});
			}
		};
	});
})();
