import './wasm_exec.js';
import wasmModule from '../build/main.wasm';

let wasmInstance = null;

async function bootstrapWasm() {
  if (wasmInstance) return wasmInstance;

  const go = new Go();
  const instance = await WebAssembly.instantiate(wasmModule, go.importObject);
  go.run(instance);
  wasmInstance = instance;
  return wasmInstance;
}

export default {
  async fetch(request, env, ctx) {
    await bootstrapWasm();
    
    return globalThis.handleHttpRequest(request, env, ctx);
  }
};
