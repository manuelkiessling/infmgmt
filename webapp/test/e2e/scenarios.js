'use strict';

describe('Infmgmt App', function() {
 describe('Overview', function() {
  
  beforeEach(function() {
    browser().navigateTo('/app/#/overview');
  });

  it('should show a list of all VM hosts', function() {
    expect(repeater('.vmhosts li').count()).toBe(1);
    expect(repeater('.vmhosts li', 'VM Hosts').column('vmhost.DnsName')).
      toEqual(['vmhost1']);
  });

 });
});
