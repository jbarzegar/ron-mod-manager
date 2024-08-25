//Actions the toolbar contains:
// - add new mod
// - open ron nexus page
// - open settings page
//

import { For } from "solid-js";
import { Actions } from "../../lib/toolbar";
// enum UpdateEvents {
//   // Settings page state is set to "open"
//   SettingsOpened = "settingsOpened",
//   // Settings Page state is set to "closed"
//   SettingsClosed = "settingsClosed",
// }
//

type ActionItem = { name: Actions; icon: string; title: string };
const actions: Array<ActionItem> = [
  {
    name: Actions.AddMod,
    icon: "i-heroicons:plus-solid",
    title: "Add New Mod",
  },
  {
    name: Actions.OpenSettings,
    icon: "i-heroicons:cog",
    title: "Open settings window",
  },
  {
    name: Actions.OpenNexus,
    icon: "i-simple-icons:nexusmods ml-auto",
    title: "Open ready or not nexusmods",
  },
];

// Holds context actions
//
export interface ToolbarProps {
  onToolItemClick(type: Actions): void;
}
export const Toolbar = (p: ToolbarProps) => {
  return (
    <div class="bg-slate-600 p-3 flex gap-3">
      <For each={actions}>
        {({ name, icon, title }) => (
          <button
            class={`color-white ${icon} text-xl`}
            onClick={() => p.onToolItemClick(name)}
            title={title}
          />
        )}
      </For>
    </div>
  );
};
