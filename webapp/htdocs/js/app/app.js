'use strict';

angular.module('infmgmt', ['infmgmtServices']).
  config(['$routeProvider', function($routeProvider) {
    $routeProvider.
      when('/overview', {templateUrl: 'partials/overview.html', controller: OverviewCtrl}).
      otherwise({redirectTo: '/overview'});
  }]);
