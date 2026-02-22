import { signal } from "@preact/signals";
import type { RefObject } from "preact";
import { useEffect } from "preact/hooks";
import type { contextmenu } from "../../../wailsjs/go/models";
import * as app from "../../../wailsjs/go/ui/App";
import * as runtime from "../../../wailsjs/runtime";

// Store log messages in a signal
const logMessages = signal<contextmenu.LogEntry[]>([]);

/**
 * subscribeForNewMessages watches for new messages coming from the go side
 * new messages will get appended to the client-side log state
 */
const subscribeForNewMessages = () =>
	runtime.EventsOn("log_event_new_messages", (data: contextmenu.LogEntry[]) => {
		logMessages.value.push(...data);
	});

const useSetupLogMenu = () => {
	// Get initial logs & subscribe to events for additional Logs
	// on mount
	useEffect(() => {
		// caller component might get re-mounted
		// don't reset log messages if it happens
		if (logMessages.value.length <= 0) {
			app.SetupLogs().then((x) => {
				logMessages.value = x;
			});
		}

		subscribeForNewMessages();
	}, []);
};

interface LogMenuProps {
	logRef: RefObject<HTMLUListElement>;
}
export const LogMenu = (p: LogMenuProps) => {
	useSetupLogMenu();

	return (
		<div class="flex flex-col relative overflow-y-scroll">
			<ul class="mt-1 mb-20" ref={p.logRef}>
				{logMessages.value.map((x) => (
					<li
						key={x.index}
						class="border border-slate-500 text-sm p-1 even:bg-slate-800"
					>
						<span class="mr-4 pl-2">{x.index + 1}</span>{" "}
						<span>{x.message}</span>
					</li>
				))}
			</ul>
		</div>
	);
};
