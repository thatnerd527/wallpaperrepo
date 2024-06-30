import { Environment } from "./IPC";

export class Storage {
  private cacheddata: string = "";
  private scope: string = "";
  private clientid: string = "";
  private storageSocket: WebSocket;
  private closed: boolean = false;
  private changeListeners: ((data: string) => void)[] = [];

  constructor(
    scope: string,
    storageSocket: WebSocket,
    initialData: string = "",
    clientid: string = ""
  ) {
    this.scope = scope;
    this.storageSocket = storageSocket;
    this.cacheddata = initialData;
    this.clientid = clientid;
    storageSocket.addEventListener("message", (event) => {
      this.cacheddata = event.data;
      this.changeListeners.forEach((listener) => {
        listener(event.data);
      });
      this.changeListeners = [];
    });
  }

  public async getScope(): Promise<string> {
    return this.scope;
  }

  public async getClientID(): Promise<string> {
    return this.clientid;
  }

  public async read(): Promise<string> {
    return this.cacheddata;
  }

  public waitForChange(): Promise<string> {
    return new Promise((resolve, reject) => {
      this.changeListeners.push(resolve);
    });
  }

  public async write(data: string) {
    this.cacheddata = data;
    if (this.closed) {
      return;
    }
    this.storageSocket.send(data);
  }

    public async close() {
    if (this.closed) {
      return;
    }
    this.storageSocket.close();
    this.closed = true;
  }
}

export class StorageManager {
  private static storages: { [key: string]: Storage } = {};

  public static repoNameGen(scopename: string, clientid: string): string {
    return `${scopename}_${clientid}`;
  }

  public static open(reponame: string, clientid: string): Promise<Storage> {
      return new Promise((resolve, reject) => {
        if (StorageManager.storages[reponame]) {
          resolve(StorageManager.storages[reponame]);
          return;
        }
      let storageSocket = new WebSocket(
        `ws://127.0.0.1:${Environment.controlPort}/storagesocket`
      );
      storageSocket.onopen = () => {
        let cacheddata = "";
        reponame = StorageManager.repoNameGen(reponame, clientid);
        storageSocket.send(reponame);
        storageSocket.addEventListener(
          "message",
          (event) => {
            cacheddata = event.data;
            let storage = new Storage(reponame, storageSocket, cacheddata);
            StorageManager.storages[reponame] = storage;
            resolve(storage);
          },
          { once: true }
        );
      };
    });
  }

  public static close(reponame: string, clientid: string) {
    reponame = StorageManager.repoNameGen(reponame, clientid);
    StorageManager.storages[reponame].close();
    delete StorageManager.storages[reponame];
  }

  public static get(reponame: string, clientid: string): Storage {
    reponame = StorageManager.repoNameGen(reponame, clientid);
    return StorageManager.storages[reponame];
  }
}
