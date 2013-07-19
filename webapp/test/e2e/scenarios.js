'use strict';

describe('Infmgmt App', function() {
 describe('Overview', function() {
  
  beforeEach(function() {
    browser().navigateTo('/app/#/overview');
  });

  it('should show a list of all VM hosts', function() {
    expect(repeater('.vmhosts li.vmhost-block').count()).toBe(1);
    expect(repeater('.vmhosts li.vmhost-block', 'VM Hosts').column('vmhost.DnsName')).
      toEqual(['vmhost1']);
  });

  it('should display the formatted total memory of the vmhost', function() {
    expect(element('p.vmhost-info').text()).toBe('Total host memory: 31.39 GiB');
  });

  it('should show one memory bar for each VM host', function() {
    expect(element('ul.vmguests-memory-bar').count()).toBe(1);
  });

  it('should show a memory bar for every VM guest of the VM host', function() {
    expect(repeater('ul.vmhosts li:nth-child(1) .vmguests-memory-bar li').count()).toBe(1);
  });

  it('should generate an info-box for every VM guest of the VM host', function() {
    expect(repeater('ul.vmhosts li:nth-child(1) div.vmguest-infobox').count()).toBe(1);
  });

  it('should display the date and time the VM guest info was last updated in the info-box', function() {
    expect(element('ul.vmhosts li:nth-child(1) div.vmguest-infobox strong.last-updated-info').text()).toMatch(new RegExp('^(\\d{4})-(\\d{2})-(\\d{2}) (\\d{2}):(\\d{2}):(\\d{2})$'));
    // TODO: replace exact date test with regex
  });

 });
});
