import { Meta, StoryObj } from "storybook-solidjs";
import { action } from "@storybook/addon-actions";

import { ModsPanel, Mod } from "./mods";

const mockMods = new Array(5).fill(null).map<Mod>((_, i) => ({
  active: i % 2 === 0,
  name: `Mod ${i}`,
  id: i.toString(),
}));

const meta = {
  title: "Panels/Mods",
  component: ModsPanel,
} as Meta<typeof ModsPanel>;

type Story = StoryObj<typeof ModsPanel>;

export const MockItems: Story = {
  args: {
    mods: mockMods,
    onUpdate: action("Mod update triggered"),
    onDelete: action("Mod delete triggered"),
  },
};

export default meta;
