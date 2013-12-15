var app = angular.module('portfolio', ['ngRoute']);

app.config(['$routeProvider', function($routeProvider){
	$routeProvider.when("/", {
		templateUrl: "/templates/index.html"
	}).when("/projects", {
		templateUrl: "/templates/projects.html",
		controller: 'ProjectCtl'
	}).when("/projects/:project", {
		templateUrl: function(params) { return "/templates/projects/" + params.project + ".html" },
		controller: 'ProjectCtl'
	}).when("/contact", {
		templateUrl: "/templates/contact.html",
		controller: 'MainCtl'
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
	$scope.sending = false;

	$scope.contacto = function(){
		$scope.sending = true;
		var email = $scope.email;
		var subject = $scope.subject;
		var message = $scope.message;
		if(email === "" && message === "" && subject === "") {
			$scope.contacted = "Try entering something first";
			$scope.sending = false;
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
			$scope.sending = false;
		}).error(function(data){
			$scope.contacted = data;
			$scope.sending = false;
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
	$scope.all = [{
		page: 'alarm',
		color: '#632E9A'
	},{
		page: 'desk',
		color: '#02A200'
	}, {
		page: 'flask_tutorial',
		color: '#082AB0'
	}, {
		page: 'gameoflife',
		color: '#B9008A'
	}, {
		page: 'gorunner',
		color: '#00B8BA'
	}, {
		page: 'jobs',
		color: '#FF5B27'
	}];

	$scope.getColor = function(page) {
		for(var i=0; i<$scope.all.length; i++) {
			if($scope.all[i].page == page){
				return $scope.all[i].color;
			}
		}
		return "white";
	}
}]);

var logos = [{
	image: "python.png",
	alt: "python"
},{
	image: "git.png",
	alt: "git"
},{
	image: "angularjs.png",
	alt: "angularjs"
},{
	image: "django.png",
	alt: "django"
},{
	image: "js.png",
	alt: "javascript"
},{
	image: "java.png",
	alt: "java"
},{
	image: "go.png",
	alt: "go"
},{
	image: "ruby.gif",
	alt: "ruby"
}];