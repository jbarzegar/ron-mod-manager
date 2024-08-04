import { Archive, ArchivesPanel, ArchivesViewProps } from "./panels/archives";
import { Mod, ModsPanel, ModItemProps, ModsPanelProps } from "./panels/mods";

interface AppProps {
  mods: {
    items: Mod[];
    handleDelete: ModsPanelProps["onDelete"];
    handleUpdate: ModsPanelProps["onUpdate"];
  };
  archives: {
    items: Archive[];
    handleInstall: ArchivesViewProps["onArchiveInstall"];
    handleUninstall: ArchivesViewProps["onArchiveUninstall"];
  };
}

function App({ mods, archives }: AppProps) {
  return (
    <div class="h-full absolute top-0 bottom-0 left-0 right-0">
      <div class="flex h-full gap-3">
        <ModsPanel
          mods={mods.items}
          onDelete={mods.handleDelete}
          onUpdate={mods.handleUpdate}
        />
        <ArchivesPanel
          archives={archives.items}
          onArchiveInstall={archives.handleInstall}
          onArchiveUninstall={archives.handleUninstall}
        />
      </div>
    </div>
  );
}

export default App;
