import './wasm_exec.js';
import mod from '../build/main.wasm';

const go = new Go();
let instance;

async function initWasm() {
  if (!instance) {
    instance = await WebAssembly.instantiate(mod, go.importObject);
    go.run(instance);
  }
}

export default {
  async fetch(request, env, ctx) {
    await initWasm();
    return globalThis.handleHttpRequest(request, env, ctx);
  }
};
