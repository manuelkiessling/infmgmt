'use strict';

function calculatePercentage(total, part) {
  return ((100 * part) / total)
}

function OverviewCtrl($scope, VmhostStore) {
  var vmhosts;

  vmhosts = VmhostStore.query({}, function success() {
    for (var vmhostId in vmhosts) {
      console.log(vmhosts[vmhostId]);
      for (var vmguestId in vmhosts[vmhostId].Vmguests) {
        vmhosts[vmhostId].Vmguests[vmguestId].memoryBlockWidth = calculatePercentage(vmhosts[vmhostId].TotalMemory, vmhosts[vmhostId].Vmguests[vmguestId].AllocatedMemory);
      }
    }
    $scope.vmhosts = vmhosts;
  });
}
//OverviewCtrl.$inject = ['$scope', 'VmhostStore'];
