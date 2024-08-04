import { ParentProps } from "solid-js";
interface PanelProps {
  color?: string;
  title: string;
}
export const Panel = (props: ParentProps<PanelProps>) => (
  <div
    class={`flex-1 rounded border-solid border-2 border-${props.color}-500 p-5`}
  >
    <h2 class="text-5xl mb-3">{props.title}</h2>
    <section>{props.children}</section>
  </div>
);
