'use strict';

function setupResource(url, $resource) {
  return $resource(url, null, {
    query: {method: "GET", isArray: true},
  });
}

angular.module('getmelangecomApp')
  .factory('mlgServers', ['$resource', function($resource) {
    return setupResource('/api/servers', $resource);
  }])
  .factory('mlgTrackers', ['$resource', function($resource) {
    return setupResource('/api/trackers', $resource);
  }]);
