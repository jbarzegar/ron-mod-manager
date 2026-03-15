# Cursed Knowledge

### [Mar 15 2026] Preact signals can't update array with push
I was fighting with a signal stored array. Since signals just uses mutable functions I thought "hey I can just use Array.push" -- "Why is my state not updating???". Well it turns out it _doesn't_ work since whatever causes the state to reconcile doesn't get triggered. 

```ts
const x = signal([])

// doesn't work
x.value.push("Hello there")

// does work
x.value = [...x, "Hello there"]

```
