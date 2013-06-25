'use strict';

function VmhostsCtrl($scope, VmhostStore) {
  $scope.vmhosts = VmhostStore.query();
}
//VmhostsCtrl.$inject = ['$scope', 'VmhostStore'];

//function VmhostDetailsCtrl($scope, $routeParams, VmhostStore) {
//  $scope.vmhost = VmhostStore.get({vmhostId: $routeParams.vmhostId});
//}
////VmhostDetailCtrl.$inject = ['$scope', '$routeParams', 'VmhostStore'];

function VmguestsCtrl($scope, $routeParams, VmguestStore) {
  $scope.vmhostId = $routeParams.vmhostId;
  $scope.vmguests = VmguestStore.query({vmhostId: $routeParams.vmhostId});
}
//VmguestsCtrl.$inject = ['$scope', 'VmguestStore'];
