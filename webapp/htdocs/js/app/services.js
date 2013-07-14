'use strict';

angular.module('infmgmtServices', ['ngResource']).
  factory('VmhostStore', function($resource) {
    return $resource('/webservice/vmhosts/:vmhostId', {}, {
      query: {
        method: 'GET',
        params: {
          vmhostId: ''
        },
        isArray: false
      }
    });
  });
