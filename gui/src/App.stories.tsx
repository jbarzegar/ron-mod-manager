import { Meta, StoryObj } from "storybook-solidjs";
import { action } from "@storybook/addon-actions";
import { Mod } from '@/panels/mods'
import { Archive } from '@/panels/archives'

import App from "./App";

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
    mods: {
      items: mockMods,
      handleUpdate: action("Update Mod"),
      handleDelete: action("Delete Mod")
    },
    archives: {
      items: mockArchives,
      handleInstall: action("Install archive"),
      handleUninstall: action("Uninstall archive")
    }
  }
};

export default meta;
