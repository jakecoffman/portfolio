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

app.run(function($rootScope, $location, $anchorScroll, $routeParams, $window) {
	//when the route is changed scroll to the proper element.
	$rootScope.$on('$routeChangeSuccess', function(newRoute, oldRoute) {
		$location.hash($routeParams.scrollTo);
		$anchorScroll();
	});

	$rootScope.$on('$viewContentLoaded', function(event) {
		$window._gaq.push(['_trackPageview', $location.path()]);
	});
});

app.controller('MainCtl', ['$scope', '$http', function($scope, $http){
	$scope.email = "";
	$scope.message = "";
	$scope.subject = "";
	$scope.contacted = "";

	$scope.contacto = function(email, subject, message){
		if(email === "" && message === "" && subject === "") {
			$scope.contacted = "Try entering something first";
			return;
		}
		$scope.contacted = "";

		var data = {
				email: email,
				subject: subject,
				message: message
			};

		console.log(data);

		$http({
			method: "POST",
			url: "/contact",
			data: data
		}).success(function(){
			$scope.contacted = "Message sent";
		}).error(function(data){
			$scope.contacted = data;
		})
	};

	$scope.logos = logos;
}]);

app.controller('ProjectCtl', ['$scope', '$routeParams', function($scope, $routeParams){
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

	$scope.project = $routeParams.project;
	$scope.all = ['alarm', 'desk', 'flask_tutorial', 'gameoflife', 'gorunner', 'jobs', 'taskpy'];
}]);

var logos = [{
	image: "python.png",
	alt: "python"
},{
	image: "js.png",
	alt: "javascript"
},{
	image: "ruby.gif",
	alt: "ruby"
},{
	image: "go.png",
	alt: "go"
},{
	image: "java.png",
	alt: "java"
},{
	image: "angularjs.png",
	alt: "angularjs"
},{
	image: "django.png",
	alt: "django"
},{
	image: "git.png",
	alt: "git"
}]