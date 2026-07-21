const loadedScripts = new Map()

export function loadScript(id, src) {
  if (loadedScripts.has(id)) return loadedScripts.get(id)
  const existing = document.getElementById(id)
  if (existing) return Promise.resolve(existing)

  const promise = new Promise((resolve, reject) => {
    const script = document.createElement('script')
    script.id = id
    script.src = src
    script.async = true
    script.defer = true
    script.onload = () => resolve(script)
    script.onerror = () => {
      loadedScripts.delete(id)
      reject(new Error(`Could not load ${src}`))
    }
    document.head.appendChild(script)
  })
  loadedScripts.set(id, promise)
  return promise
}
