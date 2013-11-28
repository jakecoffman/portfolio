var app = angular.module('portfolio', ['ngRoute']);

app.config(['$routeProvider', function($routeProvider){
	$routeProvider.when("/", {
		templateUrl: "/templates/index.html"
	});
	$routeProvider.when("/projects/:project", {
		templateUrl: function(params) { return "/templates/projects/" + params.project + ".html" },
		controller: 'ProjectCtl'
	}).otherwise({
		redirectTo: "/"
	})
}]);

app.controller('ProjectCtl', ['$scope', function($scope){
	$scope.desk = [{
		href: '/static/desk/desk0.jpg'
	},{
		href: '/static/desk/desk1.jpg'
	},{
		href: '/static/desk/desk2.jpg'
	},{
		href: '/static/desk/desk5.jpg'
	},{
		href: '/static/desk/desk3.jpg'
	},{
		href: '/static/desk/desk4.jpg'
	}];
}]);