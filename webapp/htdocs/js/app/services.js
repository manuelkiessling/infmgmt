'use strict';

angular.module('infmgmtServices', ['ngResource']).
  factory('VmhostStore', function($resource) {
    return $resource('http://localhost:port/vmhosts/:vmhostId', {}, {
      query: {
        method: 'GET',
        params: {
          port: ':8080',
          vmhostId: ''
        },
        isArray: true
      }
    });
  }).factory('VmguestStore', function($resource) {
    return $resource('http://localhost:port/vmhosts/:vmhostId/vmguests/:vmguestId', {}, {
      query: {
        method: 'GET',
        params: {
          port: ':8080',
          vmguestId: '',
        },
        isArray: true
      }
    });
  });
