import {BaseEndpointClass,  Environment, } from '../IPC';
import { RegisterClientFunctionWithER, RegisterWithEndpointRegister } from '../Reflect';
import {StorageManager, Storage} from '../StorageManager';
import * as Type from '../types';

export class StorageManagement extends BaseEndpointClass {
    @RegisterWithEndpointRegister(Environment.endpointRegister, "1.0", "OpenScopedStorage")
    async OpenScopedStorage(scope: string): Promise<Storage>  {
        let addon = Environment.originToAddons[origin];
        let addonClient = Environment.addons.find((x) => x.clientID == addon);
        return StorageManager.open(scope, addonClient.clientID);
    }

    @RegisterWithEndpointRegister(Environment.endpointRegister, "1.0", "CloseScopedStorage")
    async CloseScopedStorage(scope: string) {
        let addon = Environment.originToAddons[origin];
        let addonClient = Environment.addons.find((x) => x.clientID == addon);
        StorageManager.close(scope, addonClient.clientID);
    }

    @RegisterWithEndpointRegister(Environment.endpointRegister, "1.0", "ReadScopedStorage")
    async ReadScopedStorage(scope: string): Promise<string> {
        let addon = Environment.originToAddons[origin];
        let addonClient = Environment.addons.find((x) => x.clientID == addon);
        let storage = StorageManager.get(scope, addonClient.clientID);
        return await storage.read();
    }

    @RegisterWithEndpointRegister(Environment.endpointRegister, "1.0", "WriteScopedStorage")
    async WriteScopedStorage(scope: string, data: string): Promise<void> {
        let addon = Environment.originToAddons[origin];
        let addonClient = Environment.addons.find((x) => x.clientID == addon);
        let storage = StorageManager.get(scope, addonClient.clientID);
        return await storage.write(data);
    }

    @RegisterWithEndpointRegister(Environment.endpointRegister, "1.0", "WaitForChangeScopedStorage")
    async WaitForChangeScopedStorage(scope: string): Promise<string> {
        let addon = Environment.originToAddons[origin];
        let addonClient = Environment.addons.find((x) => x.clientID == addon);
        let storage = StorageManager.get(scope, addonClient.clientID);
        return await storage.waitForChange();
    }


}


export function StorageDummy() {

}

export default null;