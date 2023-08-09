const cacheVersion = "v1";

const addResourcesToCache = async (resources) => {
  const cache = await caches.open(cacheVersion);
  await cache.addAll(resources);
};

const putInCache = async (request, response) => {
  const cache = await caches.open(cacheVersion);
  await cache.put(request, response);
};

const networkFailure = () => {
  return new Response("Network error happened", {
    status: 408,
    headers: { "Content-Type": "text/plain" },
  });
};

const fetchAndCache = async (request) => {
  try {
    // Requests and responses can only be used once, so we need to clone them
    const response = await fetch(request.clone());
    await putInCache(request, response.clone());
    return response;
  } catch (error) {
    console.error(error);
    return null;
  }
};

const networkFirst = async ({ request }) => {
  // First try the network
  const responseFromNetwork = await fetchAndCache(request);
  if (responseFromNetwork) {
    return responseFromNetwork;
  }

  // Then try the cache
  const responseFromCache = await caches.match(request);
  if (responseFromCache) {
    return responseFromCache;
  }

  return networkFailure();
};

const cacheFirst = async ({ request }) => {
  // First try the cache
  const responseFromCache = await caches.match(request);
  if (responseFromCache) {
    return responseFromCache;
  }

  // Then try the network
  const responseFromNetwork = await fetchAndCache(request);
  if (responseFromNetwork) {
    return responseFromNetwork;
  }

  return networkFailure();
};

self.addEventListener("install", (event) => {
  event.waitUntil(addResourcesToCache(["/favicon.png", "/app.css", "/app.js"]));
});

self.addEventListener("activate", (event) => {
  clients.claim();
  event.waitUntil(
    caches.keys().then((keys) =>
      Promise.all(
        keys.map((key) => {
          if (key !== cacheVersion) {
            return caches.delete(key);
          }
        }),
      ),
    ),
  );
});

const cacheOrigins = ["fonts.gstatic.com", "fonts.googleapis.com"];

self.addEventListener("fetch", (event) => {
  const request = event.request;
  const url = new URL(request.url);

  // Only handle GET requests
  if (request.method !== "GET") return;

  // Bug fix
  // https://stackoverflow.com/a/49719964
  if (request.cache === "only-if-cached" && request.mode !== "same-origin")
    return;

  // Prefer the cache for certain assets
  if (cacheOrigins.indexOf(url.host) > -1) {
    event.respondWith(cacheFirst({ request }));
    return;
  }

  // Prefer the network for HTML resources
  event.respondWith(networkFirst({ request }));
});
