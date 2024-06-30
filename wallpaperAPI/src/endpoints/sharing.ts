import {BaseEndpointClass,  generateGUID, Environment} from '../IPC';
import { RegisterWithEndpointRegister } from '../Reflect';
import * as Type from '../types';
import {SharingIntent} from '../types';
var sharingregistry: Type.SharingRegistration[] = [];

export class SharingAndIPC extends BaseEndpointClass {
    @RegisterWithEndpointRegister(Environment.endpointRegister, "1.0", "SharingRegistration")
    async SharingRegistration(
        intent: string,
        target: string,
        filter: (intent: SharingIntent) => boolean,
        sendTo: (intent: SharingIntent) => void
    ): Promise<string> {
        let sharingregister = new Type.SharingRegistration();
        sharingregister.intent = intent;
        sharingregister.target = target;
        sharingregister.filter = filter;
        sharingregister.sendTo = sendTo;

        let registrationID = generateGUID();
        sharingregister.registrationID = registrationID;
        sharingregister.sendTo = async (intent: Type.SharingIntent) => {
            this.source.postMessage({type: "SharingIntent", intent: intent, registrationID: registrationID});
        }
        sharingregistry.push(sharingregister);
        return registrationID;
    }

    @RegisterWithEndpointRegister(Environment.endpointRegister, "1.0", "SharingIntent")
    async SharingIntent(intent: string, data: string, target: string): Promise<Type.SharingRegistration[]> {
        let sharingintent = new Type.SharingIntent();
        let registrations = sharingregistry.filter((registration) => registration.intent == sharingintent.intent && registration.target == sharingintent.target);
        let results = [];
        sharingintent.intent = intent;
        sharingintent.target = target;
        sharingintent.data = data;
        for (let registration of registrations) {
            if (registration.filter(sharingintent)) {
              results.push(registration);
            }
        }
        return results;
    }

}

export function SharingDummy() {}