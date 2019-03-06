
export interface BuilderClient {
  install(args: string[], cwd: string): any
  test(args: string[], cwd: string): any
}
