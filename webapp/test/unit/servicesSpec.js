'use strict';

describe('infmgmt services', function() {
 
  describe('VmguestinfoFreshnessCalculator', function() {

    var vmguestinfoFreshnessCalculator;

    beforeEach(module("infmgmtServices"));

    beforeEach(function() {
      inject(function($injector) {
        vmguestinfoFreshnessCalculator = $injector.get('VmguestinfoFreshnessCalculator');
      });
    });
  
    it('should return the vmguest with the oldest InfoUpdatedAt timestamp', function() {
      var vmguests = {
        "12345": {
          "Id": "12345",
          "InfoUpdatedAt": 1375105482,
        },
        "67890": {
          "Id": "67890",
          "InfoUpdatedAt": 1275105482,
        },
      };
      var oldestVmguest = vmguestinfoFreshnessCalculator.findVmguestWithOldestInfo(vmguests);
      expect(oldestVmguest.Id).toEqual("67890");
    });
  
  });

});
