import { Meta, StoryObj } from "storybook-solidjs";
import { ArchivesPanel, Archive } from "./archives";
import { action } from "@storybook/addon-actions";

const mockArchives = new Array(10).fill(null).map<Archive>((_, i) => ({
  id: `${i}`,
  name: `Mock Archive ${i}`,
  installed: i % 2 === 0,
}));

const meta = {
  title: "Panels/Archives",
  component: ArchivesPanel,
} as Meta<typeof ArchivesPanel>;

type Story = StoryObj<typeof ArchivesPanel>;

export const MockItems: Story = {
  args: {
    archives: mockArchives,
    onArchiveInstall: action("Install archive"),
    onArchiveUninstall: action("Uninstall archive"),
  },
};

export default meta;
