import { Archive, ArchivesPanel, ArchivesViewProps } from "./panels/archives";
import { Mod, ModsPanel, ModsPanelProps } from "./panels/mods";
import { Toolbar, ToolbarProps } from "./panels/toolbar";

interface AppProps {
  toolbar: {
    handleToolbarClick: ToolbarProps["onToolItemClick"];
  };
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

function App({ mods, archives, toolbar }: AppProps) {
  return (
    <div class="h-full absolute top-0 bottom-0 left-0 right-0">
      <Toolbar onToolItemClick={toolbar.handleToolbarClick} />
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
