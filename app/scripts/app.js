'use strict';

/**
 * @ngdoc overview
 * @name getmelangecomApp
 * @description
 * # getmelangecomApp
 *
 * Main module of the application.
 */
angular
  .module('getmelangecomApp', [
    'ngAnimate',
    'ngCookies',
    'ngResource',
    'ngRoute',
    'ngSanitize',
    'ngTouch'
  ])
  .config(function ($routeProvider, $locationProvider) {
    $routeProvider
      .when('/', {
        templateUrl: 'views/main.html',
        controller: 'MainCtrl'
      })
      .when('/apps', {
        templateUrl: 'views/applications.html',
        controller: 'ApplicationsCtrl'
      })
      .when('/providers', {
        templateUrl: 'views/providers.html',
        controller: 'ProvidersCtrl'
      })
      .when('/developer', {
        templateUrl: 'views/developer.html',
        controller: 'DeveloperCtrl'
      })
      .when('/developer/:article', {
        templateUrl: 'views/article.html',
        controller: 'DeveloperCtrl'
      })
      .otherwise({
        redirectTo: '/'
      });

    $locationProvider.html5Mode(false);
  });
