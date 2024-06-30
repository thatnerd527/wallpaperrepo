import {GetAddons} from '../AddonLoader';
import {BaseEndpointClass, Environment, generateGUID} from '../IPC';
import { RegisterWithEndpointRegister } from '../Reflect';
import * as Type from '../types';

export class PopupAPI extends BaseEndpointClass {
    @RegisterWithEndpointRegister(Environment.endpointRegister, "1.0", "InputRequest")
    async InputRequest(
        input_MaxLength: string,
        input_Type: string,
        input_Placeholder: string
    ): Promise<string> {
      let url = new URL(
        "http://localhost:" + Environment.controlPort + "/inputrequest"
        );
        let trackingID = generateGUID();
      url.searchParams.append("input_MaxLength", input_MaxLength);
      url.searchParams.append("trackingID", trackingID);
      url.searchParams.append("input_Type", input_Type);
      url.searchParams.append("input_Placeholder", input_Placeholder);
      let response: Response | null = null;
      while (response == null) {
        try {
          response = await fetch(url.toString());
        } catch (error) {
          console.error(error);
          console.log(`Retrying request ${trackingID}`);
        }
      }
      return await response.text();
    }

    @RegisterWithEndpointRegister(Environment.endpointRegister, "1.0", "PopupRequest")
    async PopupRequest(
        popup_URL: string,
        popup_ClientID: string,
        popup_AppName: string,
        popup_Favicon: string,
        popup_Title: string,
    ): Promise<string> {
      let url = new URL(
        "http://localhost:" + Environment.controlPort + "/popuprequest"
      );
      let trackingID = generateGUID();
      url.searchParams.append("trackingID", trackingID);
      url.searchParams.append("popup_URL", popup_URL);
      url.searchParams.append("popup_ClientID", popup_ClientID);
      url.searchParams.append("popup_AppName", popup_AppName);
      url.searchParams.append("popup_Favicon", popup_Favicon);
      url.searchParams.append("popup_Title", popup_Title);
      let response: Response | null = null;
      while (response == null) {
        try {
          response = await fetch(url.toString());
        } catch (error) {
          console.error(error);
          console.log(`Retrying request ${trackingID}`);
        }
      }
      return await response.text();
    }

}

export function PopupDummy() {

}