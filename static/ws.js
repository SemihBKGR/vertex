window.onload = function () {
    var conn

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            // on connection close
        };
        conn.onmessage = function (evt) {
            // evt.data
        }
    } else {
        // browser does not support websocket
    }
}