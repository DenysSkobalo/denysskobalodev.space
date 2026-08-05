import './wasm_exec.js';
import wasmModule from '../build/main.wasm';

let go = new Go();
let wasmInstance = null;

async function initWasm() {
  if (!wasmInstance || go.exited) {
    go = new Go();
    const instance = await WebAssembly.instantiate(wasmModule, go.importObject);
    wasmInstance = instance;
    go.run(wasmInstance);
  }
}

export default {
  async fetch(request, env, ctx) {
    try {
      await initWasm();

      const reqBody = await request.text();
      const reqPayload = {
        method: request.method,
        url: request.url,
        body: reqBody,
      };

      const rawResponse = globalThis.handleHttpRequest(reqPayload);
      if (!rawResponse) {
        throw new Error("Empty response from WASM execution");
      }

      const parsed = typeof rawResponse === 'string' ? JSON.parse(rawResponse) : rawResponse;

      return new Response(parsed.body || '', {
        status: parsed.status || 200,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': 'https://denysskobalodev.space',
          'Access-Control-Allow-Credentials': 'true',
          'Access-Control-Allow-Headers': 'Content-Type, Authorization, Cookie',
          'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
        },
      });
    } catch (err) {
      return new Response(JSON.stringify({ error: err.message, stack: err.stack }), {
        status: 500,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': 'https://denysskobalodev.space',
          'Access-Control-Allow-Credentials': 'true',
        },
      });
    }
  },
};
