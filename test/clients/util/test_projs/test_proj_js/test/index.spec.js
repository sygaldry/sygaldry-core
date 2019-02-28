
const index = require('../src/index.js');
const expect = require('chai').expect;

describe('index', () => {
  it('#todaysDate should return a string', () => {
    var todaysDate = index.todaysDate();
    expect(todaysDate).to.be.a('string');
  });
});
