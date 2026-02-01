import { render } from "preact";
import { App } from "./app";

const el = document.getElementById("app");

if (!el) {
	throw new Error("could not find element to render app");
}

render(<App />, el);
