export enum Actions {
  AddMod = "addMod",
  OpenNexus = "openNexus",
  OpenSettings = "openSettings",
}

export interface IToolbar {
  addMod(): Promise<string | undefined>;
  openNexus(): Promise<string>;
  openSettings(): Promise<string>;
}
