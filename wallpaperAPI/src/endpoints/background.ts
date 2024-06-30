import { BaseEndpointClass, Environment } from "../IPC";
import { RegisterWithEndpointRegister } from "../Reflect";
export class BackgroundManagement extends BaseEndpointClass {
  @RegisterWithEndpointRegister(Environment.endpointRegister,"1.0","SetBackground")
  async SetBackground(backgroundid: string, backgrounddata: string) {
    DotNet.invokeMethod(
      "WallpaperUI",
      "SetBackground",
      backgroundid,
      backgrounddata
    );
  }
  @RegisterWithEndpointRegister(Environment.endpointRegister,"1.0","GetBackground")
  async GetBackground() {
    return DotNet.invokeMethod("WallpaperUI", "GetBackground");
  }
  @RegisterWithEndpointRegister(Environment.endpointRegister,"1.0","SetBackgroundData")
  async SetBackgroundData(persistentbackgroundID: string,backgrounddata: string) {
    DotNet.invokeMethod(
      "WallpaperUI",
      "SetBackgroundData",
      persistentbackgroundID,
      backgrounddata
    );
  }
  @RegisterWithEndpointRegister(Environment.endpointRegister,"1.0","GetBackgroundData")
  async GetBackgroundData(persistentbackgroundID: string) {
    return DotNet.invokeMethod(
      "WallpaperUI",
      "GetBackgroundData",
      persistentbackgroundID
    );
  }
}

export function BackgroundDummy() {}
