import type { Meta } from "@storybook/preact-vite";
import { AppLayout } from "./";

export default {
	title: "Layouts/App",
	component: AppLayout,
	decorators: [
		(Story) => (
			<div
				class="absolute top-0 left-0 right-0 bottom-0"
				style={{ height: "100vh" }}
			>
				<Story />
			</div>
		),
	],
} satisfies Meta;

export const Main = {};
