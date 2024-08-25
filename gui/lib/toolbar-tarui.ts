import { invoke } from "@tauri-apps/api";
import { IToolbar } from "./toolbar";
import * as dialog from "@tauri-apps/api/dialog";

// Implements the toolbar in the context of a Tauri app
export class TauriToolbar implements IToolbar {
  // Add mod will open a file manager and emit an event to tauri
  // @returns path of zip file
  async addMod(): Promise<string | undefined> {
    // Open
    const selectedPath = await dialog.open({
      filters: [{ name: "Archive", extensions: ["zip"] }],
    });

    if (!selectedPath) {
      console.warn("No thing");
      return;
    }

    await invoke("add_mod", { modPath: selectedPath });
  }
  openNexus(): Promise<string> {
    throw new Error("Method not implemented.");
  }
  openSettings(): Promise<string> {
    throw new Error("Method not implemented.");
  }
}
