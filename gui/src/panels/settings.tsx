import { createEffect, Show } from "solid-js";

// TODO: Determine what needs to go into settings
// I realize that I'm not sure how to get settings _into_ the app yet.
// But here's the possible settings I want to support
// - Checkbox to associate nxm ready or not links to the app (juries out on how to support that)
// - Text input to set path to ready or not mod (Or should be a file system pick thing)
// - Text input to set the mod path instance directory (needed?)
// - Button to reset mod database (destructive action)

export interface SettingsPageProps {
  open: boolean;
  onSettingsOpen(): Promise<never>;
  onSettingsClosed(error?: Error): Promise<never>;
}
export const Settings = (p: SettingsPageProps) => {
  createEffect((prev) => {
    if (p.open && !prev) {
      p.onSettingsOpen();
    } else if (prev && !p.open) {
      p.onSettingsClosed();
    }

    return p.open;
  });

  return (
    <Show when={p.open}>
      <h1>Hello from settings page</h1>
    </Show>
  );
};
