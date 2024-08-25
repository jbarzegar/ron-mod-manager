import { Meta, StoryObj } from "storybook-solidjs";
import { action } from "@storybook/addon-actions";
import { Settings } from "./settings";

const meta = {
  title: "Panels/Settings",
  component: Settings,
} as Meta<typeof Settings>;

type Story = StoryObj<typeof Settings>;

export const Standard: Story = {
  args: {
    open: false,
    onSettingsOpen: action("Settings opened"),
    onSettingsClosed: action("Settings closed"),
  },
};

export default meta;
