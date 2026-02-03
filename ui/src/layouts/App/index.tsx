import { useSignal } from "@preact/signals";
import { clsx } from "clsx";
import type { ComponentChild } from "preact";
import { useEffect, useRef } from "preact/hooks";
import {
	type DragHandleDirection,
	type DragMode,
	useResize,
} from "../../hooks";

type SectionProps = {
	children: ComponentChild;
	className?: string;
	/** default: false */
	draggable?: boolean;
	/** only used when "draggable is enabled" */
	dragDirection?: DragMode;
	dragHandleDirection?: DragHandleDirection;
};

const Section = (p: SectionProps) => {
	const size = useSignal(0);
	const dragging = useSignal(false);
	const { resizeHandle, style } = useResize({
		mode: p.dragDirection || "width",
		enabled: p.draggable || false,
		draggingSignal: dragging,
		dragHandleDirection: p.dragHandleDirection || "right",
		signal: size,
		max: Infinity,
		min: 80,
		resizeHandleClass: clsx(
			`border border-red-500 ${p.dragHandleDirection}-0`,
			{
				"w-full": p.dragDirection === "height",
				"h-full": p.dragDirection === "width",
			},
		),
	});
	const sectionRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		if (sectionRef.current && p.draggable) {
			const computedS = getComputedStyle(sectionRef.current);
			if (p.dragDirection) {
				const val = computedS[p.dragDirection];
				size.value = parseInt(val, 10);
			}
		}
	}, []);

	return (
		<div
			className={`${p.className} relative`}
			style={!size.value ? {} : style}
			ref={sectionRef}
		>
			{resizeHandle}
			{p.children}
		</div>
	);
};

export const AppLayout = () => {
	return (
		<div className="h-full flex relative gap-0.5">
			<Section className="flex flex-col w-3/4 gap-0.5">
				<Section className="bg-accent-content p-4">ACTIONS</Section>
				<Section
					className="bg-accent-content h-4/5"
					draggable
					dragDirection="height"
					dragHandleDirection="bottom"
				>
					Main
				</Section>
				<Section className="flex-1/5 bg-accent-content h-1/5">Context</Section>
			</Section>
			<Section className="bg-accent-content w-1/4 relative">SIDEBAR</Section>
		</div>
	);
};
