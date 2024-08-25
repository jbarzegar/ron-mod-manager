import { Meta, StoryObj } from "storybook-solidjs";
import { action } from "@storybook/addon-actions";
import { Mod } from "@/panels/mods";
import { Archive } from "@/panels/archives";
import { v4 as uuid } from "uuid";

import App from "./App";
import { createStore, produce } from "solid-js/store";
import { Actions } from "./panels/toolbar";

const meta = {
  title: "App",
  component: App,
} as Meta<typeof App>;

type Story = StoryObj<typeof App>;

const mockMods = new Array(5).fill(null).map<Mod>((_, i) => ({
  active: i % 2 === 0,
  name: `Mod ${i}`,
  id: i.toString(),
}));

const mockArchives = new Array(10).fill(null).map<Archive>((_, i) => ({
  id: `${i}`,
  name: `Mock Archive ${i}`,
  installed: i % 2 === 0,
}));

export const Mocked: Story = {
  args: {
    toolbar: {
      handleToolbarClick: action("Toolbar item clicked"),
    },
    mods: {
      items: mockMods,
      handleUpdate: action("Update Mod"),
      handleDelete: action("Delete Mod"),
    },
    archives: {
      items: mockArchives,
      handleInstall: action("Install archive"),
      handleUninstall: action("Uninstall archive"),
    },
  },
};

interface MockedState {
  mods: Mod[];
  archives: Archive[];
}

export const Stateful: Story = () => {
  const [store, setStore] = createStore<MockedState>({
    mods: [],
    archives: [],
  });

  console.log(store.mods);

  return (
    <App
      toolbar={{
        handleToolbarClick: (t) => {
          switch (t) {
            case Actions.AddMod:
              setStore(
                "mods",
                produce((s: Mod[]) => {
                  s.push({
                    id: uuid(),
                    active: false,
                    name: `Random mod ${uuid()}`,
                  });
                }),
              );
              break;
          }
        },
      }}
      mods={{
        items: store.mods,
        handleUpdate: (index, updated) => {
          setStore(
            "mods",
            produce((s: Mod[]) => {
              s[index] = { ...s[index], ...updated };
            }),
          );
        },
        handleDelete: (index) =>
          setStore(
            "mods",
            produce((s) => {
              s.splice(index, 1);
            }),
          ),
      }}
      archives={{
        items: store.archives,
        handleInstall: action("Install archive"),
        handleUninstall: action("Uninstall archive"),
      }}
    />
  );
};

export default meta;
