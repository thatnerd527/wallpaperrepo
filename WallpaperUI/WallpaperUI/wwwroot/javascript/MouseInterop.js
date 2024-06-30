window.addEventListener('mousemove', function (e) {
    DotNet.invokeMethod('WallpaperUI', 'GlobalMouseMove', [
        e.clientX,
        e.clientY,
        e.pageX,
        e.pageY,
        e.screenX,
        e.screenY,
        e.movementX,
        e.movementY,
        e.offsetX,
        e.offsetY
    ]);
});

window.addEventListener('mouseup', function (e) {
    DotNet.invokeMethod('WallpaperUI', 'GlobalMouseUp', null);
});