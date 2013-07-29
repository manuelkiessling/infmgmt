basePath = '../';

files = [
  JASMINE,
  JASMINE_ADAPTER,
  'htdocs/js/vendor/moment.js',
  'htdocs/js/vendor/angular.js',
  'htdocs/js/vendor/angular-*.js',
  'test/lib/angular/angular-mocks.js',
  'htdocs/js/app/**/*.js',
  'test/unit/**/*.js'
];

autoWatch = true;

browsers = ['Chrome'];

junitReporter = {
  outputFile: 'test_out/unit.xml',
  suite: 'unit'
};
