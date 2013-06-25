'use strict';

angular.module('infmgmt', ['infmgmtServices']).
  config(['$routeProvider', function($routeProvider) {
    $routeProvider.
      when('/vmhosts', {templateUrl: 'partials/vmhosts.html', controller: VmhostsCtrl}).
      //when('/vmhosts/:vmhostId', {templateUrl: 'partials/vmhost-details.html', controller: VmhostDetailsCtrl}).
      when('/vmhosts/:vmhostId/vmguests', {templateUrl: 'partials/vmguests.html', controller: VmguestsCtrl}).
      otherwise({redirectTo: '/vmhosts'});
  }]);
