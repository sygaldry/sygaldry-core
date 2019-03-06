'use strict';

import { injectable } from "inversify";
import "reflect-metadata";
import TYPES from "../types";

import { PublisherClient } from './interfaces/publisher_client';
import { BuilderClient } from './interfaces/builder_client';
import { sync } from 'cross-spawn';

export interface NpmClientConfig { }

@injectable()
export class NpmClient implements PublisherClient, BuilderClient {

  npmClientConfig: NpmClientConfig;

  constructor(npmClientConfig: NpmClientConfig) {
    this.npmClientConfig = npmClientConfig;
  }

  private commandSync(command: string, args: string[], cwd: string): any {
    return sync('npm', [ command, ...args ], {
      cwd: cwd,
    });
  }

  public install(args: string[], cwd: string): any {
    return this.commandSync('install', args, cwd);
  }

  public test(args: string[], cwd: string): any {
    return this.commandSync('test', args, cwd);
  }

  public push(args: string[], cwd: string): any {
    return this.commandSync('publish', args, cwd);
  }

}
