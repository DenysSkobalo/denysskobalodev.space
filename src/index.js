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

      const cookieHeader = request.headers.get('Cookie') || '';
      let sessionToken = '';
      const match = cookieHeader.match(/admin_session=([^;]+)/);
      if (match) sessionToken = match[1];

      let activeSession = "";
      if (sessionToken && env.SESSION_KV) {
        activeSession = await env.SESSION_KV.get(`session:${sessionToken}`) || "";
      }

      const headersObj = {};
      for (const [key, value] of request.headers.entries()) {
        headersObj[key] = value;
      }

      const reqBody = await request.text();
      const reqPayload = {
        method: request.method,
        url: request.url,
        body: reqBody,
        headers: JSON.stringify(headersObj),
        env: {
          ADMIN_HASH: env.ADMIN_HASH || "",
          SALT: env.SALT || "denysskobalo_unique_salt",
          ACTIVE_SESSION: activeSession
        }
      };

      const rawResponse = globalThis.handleHttpRequest(reqPayload);
      if (!rawResponse) throw new Error("Empty response from WASM execution");
      const parsed = typeof rawResponse === 'string' ? JSON.parse(rawResponse) : rawResponse;

      if (parsed.kv_put && env.SESSION_KV) {
        await env.SESSION_KV.put(parsed.kv_put.key, parsed.kv_put.value, { expirationTtl: parsed.kv_put.ttl });
      }
      if (parsed.kv_delete && env.SESSION_KV) {
        await env.SESSION_KV.delete(parsed.kv_delete.key);
      }

      const responseHeaders = new Headers({
        'Content-Type': 'application/json',
        'Access-Control-Allow-Origin': env.ALLOWED_ORIGIN || 'https://denysskobalodev.space',
        'Access-Control-Allow-Credentials': 'true',
        'Access-Control-Allow-Headers': 'Content-Type, Authorization, Cookie',
        'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
      });

      if (parsed.headers) {
        for (const [key, value] of Object.entries(parsed.headers)) {
          responseHeaders.set(key, value);
        }
      }

      return new Response(parsed.body || '', {
        status: parsed.status || 200,
        headers: responseHeaders,
      });
    } catch (err) {
      return new Response(JSON.stringify({ error: err.message }), { status: 500 });
    }
  },
};
