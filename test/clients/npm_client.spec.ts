
import { NpmClient } from '../../src/clients/npm_client';
import { expect } from 'chai';

describe('NpmClient', () => {

  const npmClient = new NpmClient({});

  it('#install should perform an npm install', () => {
    var response = npmClient.install([], `${process.cwd()}/test/clients/util/test_projs/test_proj_js`);
    expect(response.status).to.equal(0);
  }).timeout(600000);

  it('#test should perform an npm test', () => {
    var response = npmClient.test([], `${process.cwd()}/test/clients/util/test_projs/test_proj_js`);
    expect(response.status).to.equal(0);
  }).timeout(600000);
}); 

