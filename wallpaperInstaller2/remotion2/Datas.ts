export class InstallerData {
  installpath: string;
  runatstartup: boolean = true;
  createdesktopshortcuts: boolean = true;
  enableautoupdates: boolean = true;
    sendlogs: boolean = true;
  pressedbutton: boolean = false;
  installfinished: boolean = false;

  readonly safeSetState = (
    modify: (data: InstallerData) => void,
    setState: (data: InstallerData) => void
  ) => {
    let target = new InstallerData();
    Object.assign(target, this);
    modify(target);
    setState(target);
  };

  readonly safeMerge = (
    assigntarget: InstallerData,
    setState: (data: InstallerData) => void
  ) => {
    let target = new InstallerData();
    Object.assign(target, this);
    Object.assign(target, assigntarget);
    setState(target);
  };

  constructor() {
    this.installpath = "";
  }
}