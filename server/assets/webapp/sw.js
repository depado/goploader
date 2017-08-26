var CACHE_NAME = 'gpldr-v1';
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
