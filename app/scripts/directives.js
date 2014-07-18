'use strict';

angular.module('getmelangecomApp')
  .directive('mlgAppPanel', function() {
    return {
      restrict: 'E',
      scope: {
        size: '=size',
        object: '=object',
      },
      templateUrl: 'public/views/app-panel.html',
    };
  });
