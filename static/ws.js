const actionJoin = "join"
const actionLeave = "leave"
const actionJoined = "joined"
const actionLeft = "left"
const actionMatched = "matched"
const actionMove = "move"
const actionMoved = "moved"

class Message {
    action
    data

    constructor(action) {
        this.action = action
        this.data = new Map()
    }

}

const actionJoinMessage = new Message(action = actionJoin)
const actionLeaveMessage = new Message(action = actionLeave)

let conn
let onMessageReceiveListener = []

function connect() {
    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function (evt) {
            // on connection close
        };
        conn.onmessage = function (event) {
            if (event.data.includes("\n")) {
                let messages = event.data.split("\n")
                for (const m of messages) {
                    const message = JSON.parse(m)
                    for (const listener of onMessageReceiveListener) {
                        listener(message)
                    }
                }
            } else {
                const message = JSON.parse(event.data)
                for (const listener of onMessageReceiveListener) {
                    listener(message)
                }
            }
        }
    } else {
        // browser does not support websocket
    }
}

function sendMessage(message) {
    conn.send(JSON.stringify(message))
}

function addOnMessageReceiveListener(listener) {
    onMessageReceiveListener.push(listener)
}
