/** biome-ignore-all lint/suspicious/noAssignInExpressions: math must be done */
import { computed, type Signal } from "@preact/signals";
import { clsx } from "clsx";
import type { TargetedKeyboardEvent } from "preact";
import { useEffect, useMemo, useRef } from "preact/hooks";

export type DragMode = "width" | "height";
export type DragHandleDirection = "left" | "right" | "top" | "bottom";

type UseResize = {
	enabled: boolean;
	mode: DragMode;
	min: number;
	max: number;
	signal: Signal<number>;
	dragHandleDirection?: DragHandleDirection;
	resizeHandleClass?: string;
};
export const useResize = ({
	enabled = true,
	resizeHandleClass,
	signal,
	mode,
	dragHandleDirection: dragDirection,
	max = Number.MAX_SAFE_INTEGER,
	min,
}: UseResize) =>
	useMemo(() => {
		const onMouseDown = (e: MouseEvent) => {
			if (!enabled) return;

			const { pageX: startX, pageY: startY } = e;
			const startingVal = signal.peek();

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

		const handleKeyPress = (key: string, sizeStep: number) => {
			const s = signal.peek();
			const inc = () => s + sizeStep;
			const dec = () => s - sizeStep;

			// TODO: Try to figure out if I can make this mess of switch/cases cleaner
			if (mode === "height" && dragDirection === "bottom") {
				switch (key) {
					case "ArrowDown":
						return inc();
					case "ArrowUp":
						return dec();
				}
			}

			if (mode === "width") {
				switch (key) {
					case "ArrowRight":
						return inc();
					case "ArrowLeft":
						return dec();
				}
			}

			return s;
		};
		/**
		 *
		 * Enables the panes to be resized with keyboard bindings
		 */
		const onKeyDown = (evt: TargetedKeyboardEvent<HTMLButtonElement>) => {
			const sizeStep = evt.shiftKey ? 30 : 10;
			signal.value = handleKeyPress(evt.key, sizeStep);
		};

		const style = computed(() => `${mode}: ${signal.value}px`);

		const resizeHandle = !enabled ? null : (
			<button
				type="button"
				class={clsx(
					`absolute border ${resizeHandleClass}`,
					{
						"cursor-row-resize": mode === "height",
						"cursor-col-resize": mode === "width",
						"w-full": mode === "height",
						"h-full": mode === "width",
					},
					`border-slate-700 hover:border-accent ${dragDirection}-0`,
				)}
				style={{ transition: "border-color ease 0.3s" }}
				onMouseDown={onMouseDown}
				onKeyDown={onKeyDown}
			/>
		);

		return { style, resizeHandle, onMouseDown };
	}, [signal, min, max]);

/**
 * useComputedInitialSize sets a signal value in the number of pixels computed by the effect that will run on mount
 * @param size signal that represents a number of pixels
 * @param m which direction the size represents
 * @returns ref that is used to compute the given size
 */
export function useComputedInitialSize<T extends Element>(
	size: Signal<number>,
	m: DragMode,
) {
	const r = useRef<T>(null);
	useEffect(() => {
		if (r.current) {
			const computedS = getComputedStyle(r.current);
			const val = computedS[m];
			size.value = parseInt(val, 10);
		}
	}, []);

	return r;
}
