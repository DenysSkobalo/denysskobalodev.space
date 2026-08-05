import './wasm_exec.js';
import wasmModule from '../build/main.wasm';

const go = new Go();
let wasmInstancePromise = null;

async function ensureWasmInitialized() {
  if (!wasmInstancePromise) {
    wasmInstancePromise = (async () => {
      const instance = await WebAssembly.instantiate(wasmModule, go.importObject);
      go.run(instance);
    })();
  }
  await wasmInstancePromise;
}

export default {
  async fetch(request, env, ctx) {
    try {
      await ensureWasmInitialized();

      if (typeof globalThis.handleHttpRequest !== 'function') {
        throw new Error("Wasm runtime initialization failed: handleHttpRequest is not attached to globalThis");
      }

      const reqBody = await request.text();

      const reqPayload = {
        method: request.method,
        url: request.url,
        body: reqBody,
      };

      const rawResponse = globalThis.handleHttpRequest(reqPayload);
      if (!rawResponse) {
        throw new Error("Go WASM returned empty response");
      }

      const parsed = typeof rawResponse === 'string' ? JSON.parse(rawResponse) : rawResponse;

      return new Response(parsed.body || '', {
        status: parsed.status || 200,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': env.ALLOWED_ORIGIN || 'https://denysskobalodev.space',
          'Access-Control-Allow-Credentials': 'true',
          'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
          'Access-Control-Allow-Headers': 'Content-Type, Authorization, Cookie',
        },
      });
    } catch (err) {
      return new Response(JSON.stringify({ error: err.message, stack: err.stack }), {
        status: 500,
        headers: { 
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': env.ALLOWED_ORIGIN || 'https://denysskobalodev.space',
          'Access-Control-Allow-Credentials': 'true'
        },
      });
    }
  },
};
