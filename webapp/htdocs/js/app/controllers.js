'use strict';

function calculatePercentage(total, part) {
  return Math.round(((100 * part) / total))
}

function formatMemory(memorySize) {
  return "" + Number(( memorySize / 1024 / 1024 ).toFixed(2)) + " GiB";
}

function formatDateTime(unixTimestamp) {
  return moment(new Date(unixTimestamp * 1000)).format("YYYY-MM-DD HH:mm:ss");
}

function OverviewCtrl($scope, VmhostStore, VmguestinfoFreshnessCalculator) {
  var vmhosts, infoUpdatedAtDate;
  vmhosts = VmhostStore.query({}, function success() {
    for (var vmhostId in vmhosts) {
      vmhosts[vmhostId].formattedTotalMemory = formatMemory(vmhosts[vmhostId].TotalMemory)
      //console.log(VmguestinfoFreshnessCalculator.findVmguestWithOldestInfo(vmhosts[vmhostId].Vmguests));
      for (var vmguestId in vmhosts[vmhostId].Vmguests) {
        vmhosts[vmhostId].Vmguests[vmguestId].memoryBlockWidth = calculatePercentage(vmhosts[vmhostId].TotalMemory, vmhosts[vmhostId].Vmguests[vmguestId].AllocatedMemory);
        vmhosts[vmhostId].Vmguests[vmguestId].formattedAllocatedMemory = formatMemory(vmhosts[vmhostId].Vmguests[vmguestId].AllocatedMemory)
        if (vmhosts[vmhostId].Vmguests[vmguestId].State == "running") {
          vmhosts[vmhostId].Vmguests[vmguestId].stateColor = "#7f7";
        } else {
          vmhosts[vmhostId].Vmguests[vmguestId].stateColor = "#ddd";
        }
        vmhosts[vmhostId].Vmguests[vmguestId].infoUpdatedAt = formatDateTime(vmhosts[vmhostId].Vmguests[vmguestId].InfoUpdatedAt);
      }
    }
    $scope.vmhosts = vmhosts;
  });
}
//OverviewCtrl.$inject = ['$scope', 'VmhostStore', 'VmguestinfoFreshnessCalculator'];
