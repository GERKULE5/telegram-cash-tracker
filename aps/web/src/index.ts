import { serve } from "bun";
import index from "./index.html";

const apiURL = (process.env.API_BASE_URL || "http://localhost:8080").replace(
  /\/$/,
  "",
);

const server = serve({
  port: process.env.WEB_PORT || 3000,
  routes: {
    "/api/admin/spendings": proxy,
    "/api/admin/spendings/:id": proxy,
    "/*": index,
  },
  development: process.env.NODE_ENV !== "production" && {
    hmr: true,
    console: true,
  },
});

async function proxy(request: Request) {
  const url = new URL(request.url);
  const path = url.pathname.replace(/^\/api/, "");
  const targetURL = `${apiURL}${path}${url.search}`;
  const hasBody = request.method !== "GET" && request.method !== "HEAD";

  try {
    return await fetch(targetURL, {
      method: request.method,
      headers: { "Content-Type": "application/json" },
      body: hasBody ? await request.text() : undefined,
      signal: AbortSignal.timeout(10_000),
    });
  } catch {
    return Response.json(
      { error: "The expense service is unavailable. Please try again later." },
      { status: 502 },
    );
  }
}

console.log(`Web server running on ${server.url}`);
