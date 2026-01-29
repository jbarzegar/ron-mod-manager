import { h } from 'preact';
import { useState } from "preact/hooks";
import "./app.css"

export function App() {
  const [count, setCount] = useState(0)

  return (
    <div id="App">
      <div class="p-5 flex flex-col gap-3">
        <h2 class="text-3xl">{count}</h2>
        <div className="flex gap-2">
          <button class="btn rounded-full" disabled={count <= 0} onClick={_ => setCount(c => c - 1)}>-</button>
          <button class="btn rounded-full" onClick={_ => setCount(c => c + 1)}>+</button>
        </div>
      </div>
    </div>
  )
}
