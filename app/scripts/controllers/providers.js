'use strict';

/**
 * @ngdoc function
 * @name getmelangecomApp.controller:ProvidersCtrl
 * @description
 * # ProvidersCtrl
 * Controller of the getmelangecomApp
 */
angular.module('getmelangecomApp')
  .controller('ProvidersCtrl', ['$scope', 'mlgTrackers', 'mlgServers',
    function ($scope, mlgTrackers, mlgServers) {
      $scope.trackers = mlgTrackers.query();
      $scope.servers = mlgServers.query();
  }]);
