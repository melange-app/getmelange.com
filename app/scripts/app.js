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
    'ngResource',
    'ngRoute',
  ])
  .config(function ($routeProvider, $locationProvider) {
    $routeProvider
      // Index
      .when('/', {
        templateUrl: '/public/views/main.html',
        controller: 'MainCtrl'
      })
      // Standard Pages
      .when('/apps', {
        templateUrl: '/public/views/applications.html',
        controller: 'ApplicationsCtrl'
      })
      .when('/providers', {
        templateUrl: '/public/views/providers.html',
        controller: 'ProvidersCtrl'
      })
      // Developer Pages
      .when('/developer', {
        templateUrl: '/public/views/developer.html',
        controller: 'DeveloperCtrl'
      })
      .when('/developer/:article', {
        templateUrl: '/public/views/article.html',
        controller: 'DeveloperCtrl'
      })
      // Publication Pages
      .when('/developer/publish/app', {
        templateUrl: '/public/views/publish/app.html',
        controller: 'PublishCtrl'
      })
      .when('/developer/publish/:provider', {
        templateUrl: '/public/views/publish/provider.html',
        controller: 'PublishCtrl'
      })
      // Otherwise
      .otherwise({
        redirectTo: '/'
      });

    $locationProvider.html5Mode(true);
  });
