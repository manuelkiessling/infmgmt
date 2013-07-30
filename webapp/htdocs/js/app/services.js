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
  }).
  factory('VmguestinfoFreshnessCalculator', function() {
    return {
      findVmguestWithOldestInfo: function(vmguests) {
        var oldestTimestamp = null;
        var oldestVmguest = null;
        for (var vmguestId in vmguests) {
          if (oldestTimestamp == null || vmguests[vmguestId].InfoUpdatedAt < oldestTimestamp) {
            oldestTimestamp = vmguests[vmguestId].InfoUpdatedAt;
            oldestVmguest = vmguests[vmguestId];
          }
        }
        return oldestVmguest;
      },
    };
  });
