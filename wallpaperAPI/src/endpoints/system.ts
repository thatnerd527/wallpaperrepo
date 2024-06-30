import { BaseEndpointClass, Environment } from "../IPC";
import { RegisterWithEndpointRegister } from "../Reflect";


export class SystemHooks extends BaseEndpointClass {
  @RegisterWithEndpointRegister(
    Environment.endpointRegister,
    "1.0",
    "WaitUntilSystemCleanup"
  )
  async WaitUntilSystemCleanup(): Promise<void> {}
}