'use strict';

describe('infmgmt controllers', function() {

  describe('OverviewCtrl', function() {
    var scope, ctrl, $httpBackend;

    var vmhosts = {
      "1":{
        "Id":"1",
        "DnsName":"vmhost1.example.com",
        "TotalMemory": 1024*1024,
        "Vmguests":{
          "abc":{
            "Id":"abc",
            "Name":"virtual1.example.com",
            "State":"running",
            "AllocatedMemory": 1024*512,
            "InfoUpdatedAt": 1375105482
          }
        }
      }
    };

    beforeEach(module("infmgmtServices"));

    beforeEach(
      inject(function(_$httpBackend_, $rootScope, $controller) {
      $httpBackend = _$httpBackend_;
      $httpBackend.expectGET('/webservice/vmhosts').
        respond(vmhosts);

      scope = $rootScope.$new();
      ctrl = $controller(OverviewCtrl, {$scope: scope});
    }));


    it('should create "vmhosts" model with 1 vmhost fetched from xhr', function() {
      $httpBackend.flush();


      var expectedVmhosts = {
        "1":{
          "Id":"1",
          "DnsName":"vmhost1.example.com",
          "TotalMemory": 1024*1024,
          "Vmguests":{
            "abc":{
              "Id":"abc",
              "Name":"virtual1.example.com",
              "State":"running",
              "AllocatedMemory": 1024*512,
              "InfoUpdatedAt": 1375105482,
              "memoryBlockWidth":50,
              "formattedAllocatedMemory":"0.5 GiB",
              "stateColor":"#7f7",
              "infoUpdatedAt":"2013-07-29 15:44:42",
            }
          },
          "formattedTotalMemory":"1 GiB",
        }
      };

      expect(JSON.stringify(scope.vmhosts)).toEqual(JSON.stringify(expectedVmhosts));
    });
  });

});
