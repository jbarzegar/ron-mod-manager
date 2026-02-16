import type { Signal } from "@preact/signals";
import clsx from "clsx";
import type { RefObject } from "preact";
import { useCallback, useEffect, useRef } from "preact/hooks";

const GroupDropDown = () => {
	return (
		<div className="dropdown w-full">
			<button tabIndex={0} type="button" class="btn btn-sm btn-wide">
				Groups
			</button>
			<ul
				tabindex={-1}
				class="dropdown-content menu bg-base-100 rounded-box z-1 w-52 p-2 shadow-sm"
			>
				<li>Item 1</li>
				<li>Item 2</li>
			</ul>
		</div>
	);
};

interface ActionBarProps {
	onJump(direction: "bottom" | "top"): void;
}
const ActionBar = (props: ActionBarProps) => {
	return (
		<div class="flex items-center px-4 w-full gap-2">
			<p class="flex-2">Filter</p>
			<div className="flex-1">
				<GroupDropDown />
			</div>
			<div class="flex-1">
				<input class="flex-1" placeholder="Filter" />
			</div>
			<div className="flex gap-2">
				<button
					type="button"
					class="btn btn-circle btn-xs"
					title="Scroll to top"
					onClick={() => props.onJump("top")}
				>
					⬆️
				</button>
				<button
					type="button"
					class="btn btn-circle btn-xs"
					title="Scroll to bottom"
					onClick={() => props.onJump("bottom")}
				>
					⬇️
				</button>
			</div>
		</div>
	);
};

const mockLogLines = new Array(200)
	.fill("")
	.map(
		(_, i) =>
			`Mock log ${i} Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.`,
	);

const LogMenu = (p: { logRef: RefObject<HTMLUListElement> }) => {
	return (
		<div class="flex flex-col relative overflow-y-scroll">
			<ul class="mt-1 mb-20" ref={p.logRef}>
				{mockLogLines.map((x, i) => (
					<li
						key={i}
						class="border border-slate-500 text-sm p-1 even:bg-slate-800"
					>
						{x}
					</li>
				))}
			</ul>
		</div>
	);
};

const computeHeight = (el: HTMLElement | Element) =>
	parseInt(getComputedStyle(el).height, 10);

/**  computeMenuSize calculates and sets the height of the present context menu
 * this style change is done with **side effects**
 *
 * TODO: This is something that should be refactored to use signals properly(I guess?)
 */
const computeMenuSize = (el: RefObject<HTMLElement>) =>
	requestAnimationFrame(() => {
		// noop on null footer
		if (!el.current) return;

		const footerHeight = computeHeight(el.current);
		const [actionBar, logMenu] = el.current.children[0].children as unknown as [
			Element,
			HTMLDivElement,
		];

		// Calculate the height of the logMenu. log menu is the
		// size of the current footer offset by the action bar
		// height
		// honestly this feels awful to write and wonder if
		// native code would be better
		const logMenuHeight = footerHeight - computeHeight(actionBar);

		// set the height
		logMenu.style.height = `${logMenuHeight}px`;
	});

const useScrollLogMenu = (logMenu: RefObject<HTMLElement>) =>
	useCallback(
		(direction: "top" | "bottom") => {
			if (!logMenu.current) return;

			logMenu.current.scrollIntoView({
				block: direction === "bottom" ? "end" : "start",
				behavior: "smooth",
			});
		},
		[logMenu.current],
	);

interface ContextMenuProps {
	parentSignal: Signal<number>;
	className?: string;
}
export const ContextMenu = (props: ContextMenuProps) => {
	const footerRef = useRef<HTMLElement>(null);
	const logMenuRef = useRef<HTMLUListElement>(null);

	// triggers menu compute when the parent is resized
	useEffect(() => {
		computeMenuSize(footerRef);
	}, [props.parentSignal.value]);

	// sets up scroll up/down handler
	const handleMenuJump = useScrollLogMenu(logMenuRef);

	return (
		<footer className={clsx(props.className, " bg-slate-900")} ref={footerRef}>
			<div class="flex flex-col">
				<ActionBar onJump={handleMenuJump} />
				<LogMenu logRef={logMenuRef} />
			</div>
		</footer>
	);
};
