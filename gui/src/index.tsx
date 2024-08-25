/* @refresh reload */
import "virtual:uno.css";
import "@unocss/reset/tailwind-compat.css";
import { render } from "solid-js/web";
import App from "./App";
import { Actions } from "../lib/toolbar";
import { TauriToolbar } from "../lib/toolbar-tarui";

const root = document.getElementById("root") as HTMLElement;

const toolbarImpl = new TauriToolbar();

render(
  () => (
    <App
      mods={{
        items: [],
        async handleDelete() {},
        async handleUpdate() {},
      }}
      archives={{
        items: [],
        async handleInstall() {},
        async handleUninstall() {},
      }}
      toolbar={{
        async handleToolbarClick(type) {
          console.log({ type });
          switch (type) {
            case Actions.AddMod: {
              await toolbarImpl.addMod();
              break;
            }
          }
        },
      }}
    />
  ),
  root,
);
