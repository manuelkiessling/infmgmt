'use strict';

function calculatePercentage(total, part) {
  return Math.round(((100 * part) / total))
}

function formatMemory(memorySize) {
  return "" + Number(( memorySize / 1024 / 1024 ).toFixed(2)) + " GiB";
}

function OverviewCtrl($scope, VmhostStore) {
  var vmhosts, infoUpdatedAtDate;

  vmhosts = VmhostStore.query({}, function success() {
    for (var vmhostId in vmhosts) {
      vmhosts[vmhostId].formattedTotalMemory = formatMemory(vmhosts[vmhostId].TotalMemory)
      for (var vmguestId in vmhosts[vmhostId].Vmguests) {
        vmhosts[vmhostId].Vmguests[vmguestId].memoryBlockWidth = calculatePercentage(vmhosts[vmhostId].TotalMemory, vmhosts[vmhostId].Vmguests[vmguestId].AllocatedMemory);
        vmhosts[vmhostId].Vmguests[vmguestId].formattedAllocatedMemory = formatMemory(vmhosts[vmhostId].Vmguests[vmguestId].AllocatedMemory)
        if (vmhosts[vmhostId].Vmguests[vmguestId].State == "running") {
          vmhosts[vmhostId].Vmguests[vmguestId].stateColor = "#7f7";
        } else {
          vmhosts[vmhostId].Vmguests[vmguestId].stateColor = "#ddd";
        }
        infoUpdatedAtDate = new Date(vmhosts[vmhostId].Vmguests[vmguestId].InfoUpdatedAt * 1000)
        vmhosts[vmhostId].Vmguests[vmguestId].infoUpdatedAt = infoUpdatedAtDate.toString()
        // TODO: LastUpdateInfo has to be provided by backend
      }
    }
    $scope.vmhosts = vmhosts;
  });

  $scope.showVmguestInfoBox = function(vmguest) {
  }

}
//OverviewCtrl.$inject = ['$scope', 'VmhostStore'];
