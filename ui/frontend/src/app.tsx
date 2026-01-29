import { h } from 'preact';
import { useState } from "preact/hooks";

export function App() {
  const [count, setCount] = useState(0)

  return (
    <div id="App">
      <div>
        <h2>{count}</h2>
        <button disabled={count <= 0} onClick={_ => setCount(c => c - 1)}>-</button>
        <button onClick={_ => setCount(c => c + 1)}>+</button>
      </div>
    </div>
  )
}
