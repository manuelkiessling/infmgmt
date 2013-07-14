'use strict';

describe('Infmgmt App', function() {
 describe('Overview', function() {
  
  beforeEach(function() {
    browser.navigateTo('http://localhost:8080/app/#/overview');
  });

  it('should show a list of all VM hosts', function() {
    expect(repeater('.vmhosts li', 'VM Hosts').column('vmhost.DnsName')).
      toEqual(['vmhost1']);
  });

 });
});
