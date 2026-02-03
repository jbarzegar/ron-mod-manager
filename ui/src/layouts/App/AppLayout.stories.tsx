import { AppLayout } from "./";

export default {
	title: "Layouts/App",
	component: AppLayout,
	decorators: [
		// biome-ignore lint/suspicious/noExplicitAny: I will deal with TS types another day
		(Story: any) => (
			<div
				class="absolute top-0 left-0 right-0 bottom-0"
				style={{ height: "100vh" }}
			>
				<Story />
			</div>
		),
	],
};

export const Main = {};
