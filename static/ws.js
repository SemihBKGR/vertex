const actionJoin = "join"
const actionLeave = "leave"
const actionJoined = "joined"
const actionLeft = "left"
const actionMatched = "matched"
const actionMove = "move"
const actionMoved = "moved"
const actionEnd = "end"
const actionEnded = "ended"

const dataMoveX = "x"
const dataMoveY = "y"
const dataPlayer = "p"
const dataScore = "s"
const dataWalls = "w"
const dataReason = "r"
const dataWinner = "winner"
const dataScoreP1 = "score-p1"
const dataScoreP2 = "score-p2"
const dataInitialMoveP1 = "imp1"
const dataInitialMoveP2 = "imp2"

const reasonResign = "resign"
const reasonDisconnect = "disconnect"

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
        conn.onopen = function (_) {
            window.onbeforeunload = function (_) {
                const data = new Map()
                data[dataReason] = reasonDisconnect
                const message = new Message(actionEnd)
                message.data = data
                sendMessage(message)
            }
        }
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
