'use strict';

function calculatePercentage(total, part) {
  return Math.round(((100 * part) / total))
}

function OverviewCtrl($scope, VmhostStore) {
  var vmhosts;

  vmhosts = VmhostStore.query({}, function success() {
    for (var vmhostId in vmhosts) {
      for (var vmguestId in vmhosts[vmhostId].Vmguests) {
        vmhosts[vmhostId].Vmguests[vmguestId].memoryBlockWidth = calculatePercentage(vmhosts[vmhostId].TotalMemory, vmhosts[vmhostId].Vmguests[vmguestId].AllocatedMemory);
        if (vmhosts[vmhostId].Vmguests[vmguestId].State == "running") {
          vmhosts[vmhostId].Vmguests[vmguestId].stateColor = "#7f7";
        } else {
          vmhosts[vmhostId].Vmguests[vmguestId].stateColor = "#ddd";
        }
        console.log(vmhosts[vmhostId].Vmguests[vmguestId]);
      }
    }
    $scope.vmhosts = vmhosts;
  });

  $scope.showVmguestInfoBox = function(vmguest) {
  }

}
//OverviewCtrl.$inject = ['$scope', 'VmhostStore'];
