import { Environment } from "./IPC";
import { Addon } from "./types";

export async function GetAddons(): Promise<Addon[]> {
    let response = await fetch("http://localhost:" + Environment.controlPort + "/addons");
    return await response.json();
}

export async function GetAddonOrigin(clientid: string): Promise<string> {
    let url = new URL(
      "http://localhost:" + Environment.controlPort + "/getaddonorigin"
    );
    url.searchParams.append("clientID", clientid);
    let response = await fetch(url.toString());
    return await response.text();
}


export async function BootstrapOriginsToAddons() {
    let addons = await GetAddons();
    for (let addon of addons) {
        let origin = await GetAddonOrigin(addon.clientID);
        Environment.originToAddons[origin] = addon.clientID;
    }
}