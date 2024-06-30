export class InputRequest {
  public input_MaxLength: string;
  public trackingID: string;
  public input_Type: string;
  public input_Placeholder: string;
}

export class OpenPanel {
  public loaderPanelID: string;
}

export class ClosePanel {
  public loaderPanelID: string;
  public persistentPanelID: string;
}

export class GetPanelSize {
  public loaderPanelID: string;
  public persistentPanelID: string;
}

export class SetPanelSize {
  public loaderPanelID: string;
  public persistentPanelID: string;
  public width: string;
  public height: string;
}

export class PopupRequest {
  public trackingID: string;
  public popup_URL: string;
  public popup_ClientID: string;
  public popup_AppName: string;
  public popup_Favicon: string;
  public popup_Title: string;
}

export class SharingIntent {
    public intent: string;
    public data: string;
    public target: string;
}

export class SharingRegistration {
    public intent: string;
    public target: string;
    public filter: (intent: SharingIntent) => boolean;
    public registrationID: string;
    public sendTo: (intent: SharingIntent) => void;
}

export class Addon {
  public name: string;
  public version: string;
  public description: string;
  public author: string;
  public clientID: string;
  public bootstapExecutable: string;
  public enableRestart: boolean;
}

export class OpenScopedStorage {
  public scope: string;
}

export class CloseScopedStorage {
  public scope: string;
}

export class ReadScopedStorage {
  public scope: string;
}

export class WriteScopedStorage {
  public scope: string;
  public data: string;
}

