'use strict';

describe('infmgmt controllers', function() {

  describe('OverviewCtrl', function() {
    var scope, ctrl, $httpBackend;

    var vmhosts = {
      "1":{
        "Id":"1",
        "DnsName":"vmhost1.example.com",
        "Vmguests":{
          "abc":{
            "Id":"abc",
            "Name":"virtual1.example.com",
            "State":"running"
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

      expect(JSON.stringify(scope.vmhosts)).toEqual(JSON.stringify(vmhosts));
    });
  });

});
