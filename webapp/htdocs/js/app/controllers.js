'use strict';

function OverviewCtrl($scope, VmhostStore) {
  $scope.vmhosts = VmhostStore.query();
}
//OverviewCtrl.$inject = ['$scope', 'VmhostStore'];
