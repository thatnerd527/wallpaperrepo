import {BaseEndpointClass, Environment} from '../IPC';
import {RegisterWithEndpointRegister, RegisterClientFunctionWithER} from '../Reflect';

export class PanelManagement extends BaseEndpointClass {
  @RegisterWithEndpointRegister(
    Environment.endpointRegister,
    "1.0",
    "ClosePanel"
  )
  async ClosePanel(persistentPanelID: string) {
      DotNet.invokeMethod("WallpaperUI", "ClosePanel", persistentPanelID);
  }
  @RegisterWithEndpointRegister(
    Environment.endpointRegister,
    "1.0",
    "OpenPanel"
  )
  async OpenPanel(panelid: string) {
    DotNet.invokeMethod("WallpaperUI", "OpenPanel", panelid);
  }
  @RegisterWithEndpointRegister(
    Environment.endpointRegister,
    "1.0",
    "GetPanelSize"
  )
  async GetPanelSize(persistentPanelID: string) {
    return DotNet.invokeMethod("WallpaperUI", "GetPanelSize", persistentPanelID);
  }
  @RegisterWithEndpointRegister(
    Environment.endpointRegister,
    "1.0",
    "SetPanelSize"
  )
  async SetPanelSize(persistentPanelID: string, width: number, height: number) {
    DotNet.invokeMethod("WallpaperUI", "SetPanelSize", persistentPanelID, width, height);
  }
  @RegisterWithEndpointRegister(
    Environment.endpointRegister,
    "1.0",
    "GetPanelPosition"
  )
  async GetPanelPosition(persistentPanelID: string) {
    return DotNet.invokeMethod("WallpaperUI", "GetPanelPosition", persistentPanelID);
  }
  @RegisterWithEndpointRegister(
    Environment.endpointRegister,
    "1.0",
    "SetPanelPosition"
  )
  async SetPanelPosition(persistentPanelID: string, x: number, y: number) {
    DotNet.invokeMethod("WallpaperUI", "SetPanelPosition", persistentPanelID, x, y);

  }
  @RegisterWithEndpointRegister(
    Environment.endpointRegister,
    "1.0",
    "GetPanelVisibility"
  )
  async GetPanelVisibility(persistentPanelID: string) {
    return DotNet.invokeMethod("WallpaperUI", "GetPanelVisibility", persistentPanelID);
  }
  @RegisterWithEndpointRegister(
    Environment.endpointRegister,
    "1.0",
    "SetPanelData"
  )
  async SetPanelData(persistentPanelID: string, data: string) {
    DotNet.invokeMethod("WallpaperUI", "SetPanelData", persistentPanelID, data);
  }
  @RegisterWithEndpointRegister(
    Environment.endpointRegister,
    "1.0",
    "GetPanelData"
  )
  async GetPanelData(persistentPanelID: string) {
    return DotNet.invokeMethod("WallpaperUI", "GetPanelData", persistentPanelID);
  }
  @RegisterWithEndpointRegister(
    Environment.endpointRegister,
    "1.0",
    "SetPanelHeader"
  )
  async SetPanelHeader(
    persistentPanelID: boolean,
    titlebarvisible: boolean,
    enableresize: boolean,
    enableclose: boolean,
    enabledrag: boolean
  ) {
    DotNet.invokeMethod(
      "WallpaperUI",
      "SetPanelHeader",
      persistentPanelID,
      titlebarvisible,
      enableresize,
      enableclose,
      enabledrag
    );
  }
  @RegisterWithEndpointRegister(
    Environment.endpointRegister,
    "1.0",
    "GetPanelHeader"
  )
  async GetPanelHeader(persistentPanelID: string) {
    return DotNet.invokeMethod("WallpaperUI", "GetPanelHeader", persistentPanelID);
  }
}
export function PanelDummy() {
}