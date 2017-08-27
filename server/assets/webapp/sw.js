var CACHE_NAME = 'gpldr-v2';
var urlsToCache = [
    "/simple",
    "/static/jquery.min.js",
    "/static/clipboard.min.js",
    "/static/toastr.min.js",
    "/static/custom.js",
    "/static/milligram.min.css",
    "/static/toastr.css",
    "/static/style.css"
];

self.addEventListener('install', function(event) {
  event.waitUntil(
    caches.open(CACHE_NAME).then(function(cache) {
        return cache.addAll(urlsToCache);
    })
  );
});

self.addEventListener('fetch', function(event) {
  event.respondWith(
    caches.match(event.request)
      .then(function(response) {
        if (response) {
          return response;
        }
        return fetch(event.request);
      }
    )
  );
});

