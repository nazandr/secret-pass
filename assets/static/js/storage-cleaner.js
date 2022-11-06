fetch('/api/lifespan', {
    method: 'GET',
}).then(function (response) {
    return response.json();
}).then(function (data) {
    clearStorage(data.lifespan);
}).catch(function (error) {
    clearStorage(3600000);
    console.log(error);
});

function clearStorage(lifespan) {
    var toRemove = Array();
    for (let i = 0; i < localStorage.length; i++) {
        let key = localStorage.key(i);
        var value = JSON.parse(localStorage.getItem(key));
        if (value.time < Date.now() - lifespan) {
            toRemove.push(key);
        };
    };
    for (let i = 0; i < toRemove.length; i++) {
        localStorage.removeItem(toRemove[i]);
    };
};
