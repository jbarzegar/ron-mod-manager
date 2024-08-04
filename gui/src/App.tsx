import { createSignal } from "solid-js";
import logo from "./assets/logo.svg";
import { invoke } from "@tauri-apps/api/tauri";
import "./App.css";

interface Archive {
  id: string
  name: string
}

const mockArchives = new Array(10).fill(null).map<Archive>((_, i) => ({ id: `${i}`, name: `Mock Archive ${i}` }))

function App() {
  const [greetMsg, setGreetMsg] = createSignal("");
  const [name, setName] = createSignal("");

  async function greet() {
    // Learn more about Tauri commands at https://tauri.app/v1/guides/features/command
    setGreetMsg(await invoke("greet", { name: name() }));
  }

  return (
    <div class="container">

      <ul>
        {mockArchives.map(x => <li>{x.name}</li>)}
      </ul>

    </div>
  );
}

export default App;
