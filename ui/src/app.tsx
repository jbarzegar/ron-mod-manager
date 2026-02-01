import { signal } from "@preact/signals";
import "./app.css";

const count = signal(0);

const Counter = () => {
	return (
		<div class="p-5 flex flex-col gap-3">
			<h2 class="text-3xl">{count}</h2>
			<div class="flex gap-2">
				<button
					type="button"
					class="btn rounded-full"
					disabled={count.value <= 0}
					onClick={(_) => count.value--}
				>
					-
				</button>
				<button
					type="button"
					class="btn rounded-full"
					onClick={(_) => count.value++}
				>
					+
				</button>
			</div>
		</div>
	);
};

export function App() {
	return (
		<div id="App">
			<Counter />
		</div>
	);
}
