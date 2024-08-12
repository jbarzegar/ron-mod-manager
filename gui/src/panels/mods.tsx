import { Panel } from "@/components/panel";
import { For } from "solid-js";

// Entity of a Mod in the frontend
export interface Mod {
  id: string;
  name: string;
  active: boolean;
}

type ModToUpdate = Partial<Pick<Mod, "name" | "active">>;

export interface ModItemProps {
  mod: Mod;
  onUpdate(update: ModToUpdate): void;
  onDelete(): void;
}
const ModItem = (p: ModItemProps) => {
  return (
    <li class="bg-stone-100 color-gray-800 p-3 flex gap-2">
      <input
        type="checkbox"
        checked={p.mod.active}
        onChange={(e) => p.onUpdate({ active: !e.currentTarget.checked })}
        readOnly
      />{" "}
      <p>{p.mod.name}</p>
      <button
        title={`Delete mod ${p.mod.name}`}
        class="i-heroicons:trash-16-solid text-red-500 ml-auto"
        onClick={p.onDelete}
      />
    </li>
  );
};

export interface ModsPanelProps {
  mods: Mod[];
  onUpdate(index: number, update: ModToUpdate): void;
  onDelete(index: number): void;
}
export const ModsPanel = (p: ModsPanelProps) => {
  return (
    <Panel title="Mods" color="green">
      <ul class="list-none gap-3 flex flex-col">
        <For each={p.mods}>
          {(x, i) => (
            <ModItem
              mod={x}
              onUpdate={(x) => p.onUpdate(i(), x)}
              onDelete={() => p.onDelete(i())}
            />

          )}
        </For>
      </ul>
    </Panel>
  );
};
