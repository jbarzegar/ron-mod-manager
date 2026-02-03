/** biome-ignore-all lint/suspicious/noAssignInExpressions: math must be done */
import { computed, type Signal } from "@preact/signals";
import { clsx } from "clsx";
import { useMemo } from "preact/hooks";

export type DragMode = "width" | "height";
export type DragHandleDirection = "left" | "right" | "top" | "bottom";

type UseResize = {
	enabled: boolean;
	mode: DragMode;
	min: number;
	max: number;
	signal: Signal<number>;
	draggingSignal: Signal<boolean>;
	dragHandleDirection: DragHandleDirection;
	resizeHandleClass: string;
};
export const useResize = ({
	enabled = true,
	resizeHandleClass,
	signal,
	draggingSignal,
	mode,
	dragHandleDirection: dragDirection,
	max = Number.MAX_SAFE_INTEGER,
	min,
}: UseResize) =>
	useMemo(() => {
		const onMouseDown = (e: MouseEvent) => {
			if (!enabled) return;

			draggingSignal.value = true;

			const { pageX: startX, pageY: startY } = e;
			const startingVal = signal.value;

			const updater = (e: MouseEvent) => {
				const calcDelta = () => {
					switch (mode) {
						case "height":
							return startingVal + e.pageY - startY;
						case "width":
							return startingVal + e.pageX - startX;
					}
				};
				// update the signal
				signal.value = Math.max(Math.min(max, calcDelta()));
				return;
			};

			// setup listener to compute and update the width
			window.addEventListener("mousemove", updater);

			// setup listener which will remove the update listener
			window.addEventListener(
				"mouseup",
				() => window.removeEventListener("mousemove", updater),
				{ once: true },
			);

			// prevent any other interaction during resize
			e.preventDefault();
			e.stopPropagation();
		};

		// this is the trick, computed signal which we can then
		// pass directly to the style prop
		const style = computed(() => `${mode}: ${signal.value}px`);

		const resizeHandle = !enabled ? null : (
			<button
				type="button"
				class={clsx(`absolute ${resizeHandleClass}`, {
					"cursor-row-resize": mode === "height",
					"cursor-col-resize": mode === "width",
				})}
				onMouseDown={onMouseDown}
				onMouseUp={() => (draggingSignal.value = false)}
			/>
		);

		return { style, resizeHandle, onMouseDown };
	}, [signal, min, max]);
