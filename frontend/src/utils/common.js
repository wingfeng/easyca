function debounce(fn, delay) {
    var timer = null
    var delay = delay || 500
    return function() {
        var args = arguments
        var that = this
        clearTimeout(timer)
        timer = setTimeout(() => {
            fn.apply(that, args)
        }, delay);
    }
}

export {
    debounce
}