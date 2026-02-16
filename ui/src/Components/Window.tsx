import type { ComponentChildren } from "preact";

interface WindowProps {
	children: ComponentChildren;
}

export const Window = (p: WindowProps) => {
	return (
		<div style={{ height: "100vh", overflow: "hidden" }}>{p.children}</div>
	);
};
