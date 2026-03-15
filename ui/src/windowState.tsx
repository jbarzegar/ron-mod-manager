import { signal } from "@preact/signals";
import type { contextmenu } from "../wailsjs/go/models";

export const mainSectionSize = signal(0);

export const logMessages = signal<contextmenu.LogEntry[]>([]);
