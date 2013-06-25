'use strict';

describe('Infmgmt App', function() {
 describe('VM Hosts Lists', function() {
  
  beforeEach(function() {
    browser.navigateTo('../../htdocs/index.html#vmhosts');
  });

  it('should show a list of all VM hosts', function() {
    expect(repeater('.vmhosts li', 'VM Hosts').column('vmhost.DnsName')).
      toEqual(['vmhost1']);
  });

 });
});
