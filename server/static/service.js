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

const networkFirst = async ({ request }) => {
  try {
    // First try the network
    const responseFromNetwork = await fetch(request.clone());
    await putInCache(request, responseFromNetwork.clone());
    return responseFromNetwork;
  } catch (error) {
    console.error(error);
  }

  // Then try the cache
  const responseFromCache = await caches.match(request);
  if (responseFromCache) {
    return responseFromCache;
  }

  return networkFailure();
};

const cacheFirst = async ({ request }) => {
  // First try to get the resource from the cache
  const responseFromCache = await caches.match(request);
  if (responseFromCache) {
    return responseFromCache;
  }

  // Next try to get the resource from the network
  try {
    const responseFromNetwork = await fetch(request.clone());
    // response may be used only once
    // we need to save clone to put one copy in cache
    // and serve second one
    putInCache(request, responseFromNetwork.clone());
    return responseFromNetwork;
  } catch (error) {
    console.error(error);
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

self.addEventListener("fetch", (event) => {
  const request = event.request;

  // Bug fix
  // https://stackoverflow.com/a/49719964
  if (
    event.request.cache === "only-if-cached" &&
    event.request.mode !== "same-origin"
  )
    return;

  //   // Prefer the cache for assets
  //   if (!request.headers.get("Accept").includes("text/html")) {
  //     event.respondWith(cacheFirst({ request }));
  //     return;
  //   }

  // Prefer the network for HTML resources
  event.respondWith(networkFirst({ request }));
});
