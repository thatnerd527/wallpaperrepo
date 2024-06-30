import {
  RegisterWithEndpointRegister,
} from "./Reflect";
import { Addon } from "./types";

export function generateGUID(): string {
  return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, function (c) {
    const r = (Math.random() * 16) | 0;
    const v = c === "x" ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}

export function javaStringHashCode(str: string) {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = 31 * hash + str.charCodeAt(i);
    hash |= 0; // Convert to 32bit integer
  }
  return hash;
}


export class TypeAndParam {
  type: string;
  param: string;

  constructor(type: string, param: string) {
    this.type = type;
    this.param = param;
  }
}

export class Endpoint2 {
  _endpointID: string;
  _endpointVersion: string;
  _clientName: string;
  _runnerFunction: (
    source: MessageEventSource,
    origin: string,
    args: any[]
  ) => Promise<any>;
  _typesAndParams: TypeAndParam[];
  _returnType: string;
}

export class EndpointRegister2 {
  endpoints: Endpoint2[] = [];
  clientEndpoints: string[] = [];
  registerEndpoint(
    functionname: string,
    functionrunner: (
      source: MessageEventSource,
      origin: string,
      args: any[]
    ) => Promise<any>,
    endpointversion: string,
    clientName: string,
    endpointID: string,
    typesAndParams: TypeAndParam[],
    returntype: string
  ) {
    let endpoint = new Endpoint2();
    endpoint._endpointID = endpointID;
    endpoint._endpointVersion = endpointversion;
    endpoint._clientName = clientName;
    endpoint._runnerFunction = functionrunner;
    endpoint._typesAndParams = typesAndParams;
    endpoint._returnType = returntype;

    this.endpoints.push(endpoint);
    // console.log(this);
    window.addEventListener("message", (event) => {
      if (event.data == null || event.data == undefined) {
        return;
      }

      if (
        event.data.endpointid == null ||
        event.data.endpointid == undefined ||
        event.data.endpointid != endpoint._endpointID
      ) {
        return;
      }
      if (!Object.keys(Environment.originToAddons).includes(event.origin)) {
        return;
      }
      switch (event.data.type) {
        case "invokefunction":
          let invocationid = event.data.invocationid;
          if (
            invocationid == null ||
            invocationid == undefined ||
            invocationid == ""
          ) {
            return;
          }
          endpoint
            ._runnerFunction(event.source, event.origin, event.data.data)
            .then((result) => {
              event.source.postMessage({
                type: "invokefunctionresponse",
                endpointID: endpoint._endpointID,
                invocationid: invocationid,
                data: result,
                error: null,
              });
            })
            .catch((error) => {
              event.source.postMessage({
                type: "invokefunctionresponse",
                endpointID: endpoint._endpointID,
                invocationid: invocationid,
                data: null,
                error: error,
              });
            });
          break;
      }
    });
    return this;
  }

  generateClientSDK(): string {
    let sdk = "";
    console.log(this);

    for (let i of this.clientEndpoints) {
      sdk += i;
    }
    for (let i of this.endpoints) {
      let uniqued = generateGUID();
      sdk += `export function ${i._clientName}(${i._typesAndParams
        .flatMap((x) => [x.param + ": " + x.type])
        .join(",")}): ${
        i._returnType == "Promise" ? "Promise<any>" : i._returnType
      } {
                return new Promise((resolve, reject) => {
                    let invocationid = btoa(Math.random().toString());
                    window.parent.postMessage({
                        type: "invokefunction",
                        endpointid: ${i._endpointID},
                        invocationid: invocationid,
                        data: [${i._typesAndParams
                          .map((x) => x.param)
                          .join(",")}],
                    }, "*");
                    let callback = (event) => {
                        if (event.data == null || event.data == undefined) {
                            return;
                        }
                        if (event.data.type == "invokefunctionresponse" && event.data.invocationid == invocationid) {
                            if (event.data.error != null) {
                                reject(event.data.error);
                            } else {
                                window.removeEventListener("message", callback);
                                resolve(event.data.data);
                            }
                        }
                    };
                    window.addEventListener("message", callback);
                });
            `;
    }
    return sdk;
  }
}

export class BaseEndpointClass {
  source: MessageEventSource;
  origin: string;

  constructor(source: MessageEventSource, origin: string) {
    this.source = source;
    this.origin = origin;
  }
}



export class Environment {
  static controlPort: number = 8080;
  static originToAddons: { [key: string]: string } = {};
  static storageSocket: WebSocket;
  static addons: Addon[] = [];
  static endpointRegister: EndpointRegister2 = new EndpointRegister2();
}
