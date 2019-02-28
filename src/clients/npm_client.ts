'use strict';

import { ArtifactRepoClient } from './interfaces/publisher_client';
import * as child from 'child_process';
import { BuildToolClient } from './interfaces/builder_client';

export interface NpmClientConfig { }

export class NpmClient implements ArtifactRepoClient, BuildToolClient {

  npmClientConfig: NpmClientConfig;

  constructor(npmClientConfig: NpmClientConfig) {
    this.npmClientConfig = npmClientConfig;
  }

  public install(args: string[], cwd: string) {
    return child.spawnSync('npm', [ 'install', ...args ], {
      cwd: cwd,
    });
  }

  public test(args: string[], cwd: string) {
    return child.spawnSync('npm', [ 'test', ...args ], {
      cwd: cwd,
    });
  }

}
