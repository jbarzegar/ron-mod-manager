import { For } from "solid-js";
import { Panel } from "@/components/panel";

export interface Archive {
  name: string;
  installed: boolean;
}

export interface ArchivesViewProps {
  archives: Archive[];
  onArchiveInstall(index: number): void;
  onArchiveUninstall(index: number): void;
}

export const ArchivesPanel = (p: ArchivesViewProps) => {
  const handleDoubleClick = (index: number, installed: boolean) => (): void => {
    if (installed) p.onArchiveUninstall(index);
    else p.onArchiveInstall(index);
  };

  return (
    <Panel title="Archives" color="red">
      <ul class="flex flex-col gap-3">
        <For each={p.archives}>
          {(archive, index) => (
            <li
              class="flex gap-2"
              onDblClick={handleDoubleClick(index(), archive.installed)}
            >
              {archive.name}{" "}
              {archive.installed && (
                <p class="color-green-500 ml-3">[Installed]</p>
              )}
            </li>
          )}
        </For>
      </ul>
    </Panel>
  );
};
