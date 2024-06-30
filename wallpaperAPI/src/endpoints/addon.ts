import {BaseEndpointClass, Environment} from '../IPC';
import {RegisterClientFunctionWithER} from '../Reflect';
export class AddonEndpoints extends BaseEndpointClass {
    @RegisterClientFunctionWithER(
    Environment.endpointRegister,
    "GetEnvironmentInfo"
  )
  GetEnvironmentInfo(): any {
    let queries = new URL(window.location.href).searchParams;
    let assembled: { [key: string]: string } = {};
    for (let query of queries) {
      assembled[query[0]] = query[1];
    }
    return assembled;
  }
}

export function AddonDummy() {

}