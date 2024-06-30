import "reflect-metadata";
import {
  EndpointRegister2,
  javaStringHashCode,
  TypeAndParam,
} from "./IPC";
import * as acorn2 from "acorn-loose";

function getFunctionBody(func: Function) {
    let code = func.toString();
    let parsed = acorn2
      .parse(code, { ecmaVersion: "latest" })
        .body.find((x) => x.type == "BlockStatement");
    return code.substring(parsed.start, parsed.end);
}



function getParamNames(func: Function): string[] {
  const fnStr = func.toString().replace(/((\/\/.*$)|(\/\*[\s\S]*?\*\/))/gm, "");
  const result = fnStr
    .slice(fnStr.indexOf("(") + 1, fnStr.indexOf(")"))
    .match(/([^\s,]+)/g);
  return result === null ? [] : result;
}

export function RegisterClientFunctionWithER(endpointregister: EndpointRegister2, clientname: string) {
    return function (target: any, propertyKey: string | symbol) {
        let func = target[propertyKey];
        let parameternames = getParamNames(func);
        let types: any[] = Reflect.getMetadata(
            "design:paramtypes",
            target,
            propertyKey
        ).map((x: any) => {
            switch (x.name) {
                case "String":
                    return "string";
                case "Number":
                    return "number";
                case "Boolean":
                    return "boolean";
                case "Function":
                    return "Function";
                default:
                    return x.name;
            }
        });
        let returntype = Reflect.getMetadata("design:returntype", target, propertyKey) == undefined ? "void" : Reflect.getMetadata("design:returntype", target, propertyKey).name;
        let body = getFunctionBody(func);
        let typesAndParams = types.map(
          (x) => new TypeAndParam(x, parameternames.shift())
        );
        endpointregister.clientEndpoints.push(
            `function ${clientname}(${typesAndParams
          .flatMap((x) => [x.param + ": " + x.type])
                .join(",")}): ${returntype == "Promise" ? "Promise<any>" : returntype} ${body}
`);

    }
}

export function RegisterWithEndpointRegister(
  endpointregister: EndpointRegister2,
  endpointversion: string,
  clientname: string
) {
  return function (target: any, propertyKey: string | symbol) {
   // console.log(Reflect.getMetadataKeys(target, propertyKey));
   // console.log(Reflect.getMetadata("design:paramtypes", target, propertyKey));
    let func = target[propertyKey];
    let parameternames = getParamNames(func);
      let classobject = target;
    let types: any[] = Reflect.getMetadata(
      "design:paramtypes",
      target,
      propertyKey
    ).map((x: any) => {
        switch (x.name) {
            case "String":
                return "string";
            case "Number":
                return "number";
            case "Boolean":
                return "boolean";
            case "Function":
                return "Function";
            default:
                return x.name;
        }
      });
      let returntype = Reflect.getMetadata("design:returntype", target, propertyKey);

    let functionrunner = async (
      source: MessageEventSource,
      origin: string,
      args: any[]
    ) => {
        let instancedclass = new classobject(source, origin);
      return await func.apply(instancedclass, args);
      };
      //parseTest(func.toString());
    let typesAndParams = types.map(
      (x) => new TypeAndParam(x, parameternames.shift())
    );
    endpointregister.registerEndpoint(
      propertyKey.toString(),
      functionrunner,
      endpointversion,
      clientname,
      javaStringHashCode(
        clientname + endpointversion + func.toString()
      ).toString(),
        typesAndParams,
      returntype.name
    );
  };
}
