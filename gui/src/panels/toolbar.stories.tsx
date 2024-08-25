import { Meta, StoryObj } from "storybook-solidjs";
import { action } from "@storybook/addon-actions";
import { Toolbar } from "./toolbar";

const meta = {
  title: "Panels/Toolbar",
  component: Toolbar,
} as Meta<typeof Toolbar>;

type Story = StoryObj<typeof Toolbar>;

export const Standard: Story = {
  args: {
    onToolItemClick: action("toolbar item clicked"),
  },
};

export default meta;
