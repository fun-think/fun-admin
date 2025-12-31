import { createVNode } from 'vue'

const compMap = /* @__PURE__ */ new Map()
export function useCompConsumer() {
  const route = useRoute()
  const getComp = (component) => {
    if (!route.name)
      return component
    const compName = component?.type?.name
    const routeName = route.name
    if (compMap.has(routeName))
      return compMap.get(routeName)
    const node = component
    if (compName && compName === routeName) {
      compMap.set(routeName, node)
      return node
    }
    const Comp = createVNode(node)
    // Check if Comp.type is a symbol before trying to assign properties to it
    if (typeof Comp.type === 'symbol') {
      // If it's a symbol, we can't add properties to it, so we return the component as is
      return component
    }
    if (!Comp.type)
      Comp.type = {}
    Comp.type.name = routeName
    compMap.set(routeName, Comp)
    return Comp
  }
  return {
    getComp,
  }
}
