import { computed, useSignal } from "@preact/signals";
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
	const { resizeHandle, style } = useResize({
		mode: p.dragDirection || "width",
		enabled: p.draggable || false,
		dragHandleDirection: p.dragHandleDirection || "right",
		signal: size,
		max: Infinity,
		min: 80,
		resizeHandleClass: `border-${(p.dragHandleDirection || "right")[0]}-3`,
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
			<div className="p-5">{p.children}</div>
		</div>
	);
};

export const AppLayout = () => {
	const width = useSignal(0);

	const sectionRef = useRef<HTMLDivElement>(null);
	const { resizeHandle } = useResize({
		dragHandleDirection: "left",
		enabled: true,
		max: Infinity,
		min: 300,
		mode: "width",
		signal: width,
		resizeHandleClass: "border-l-3",
	});

	const mainViewStyle = computed((): preact.CSSProperties => {
		if (width.value) return { flexBasis: width.value };
		return { flex: 4 };
	});

	return (
		<div className="h-full flex relative">
			<div
				className="flex flex-col"
				ref={sectionRef}
				style={mainViewStyle.value}
			>
				<nav className="bg-slate-900 p-4 border-slate-700 border-b-3">
					ACTIONS
				</nav>
				<Section
					className="bg-slate-900 h-4/5"
					draggable
					dragDirection="height"
					dragHandleDirection="bottom"
				>
					Main
				</Section>
				<footer className="flex-1/5 bg-slate-900 h-1/5">Context</footer>
			</div>
			<aside className="bg-slate-900 flex-1 relative">
				{resizeHandle}
				<div className="p-5">div.w-full SIDEBAR</div>
			</aside>
		</div>
	);
};
