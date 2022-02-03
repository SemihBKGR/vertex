const yourColor = "green"
const opponentColor = "blue"

let gp

document.getElementById("board-div").style.visibility = "hidden"
document.getElementById("queue-div").style.visibility = "hidden"

document.getElementById("start-button").addEventListener("click", function (_) {
    console.log(gameBoard.blocks)
    document.getElementById("start-div").style.visibility = "hidden"
    document.getElementById("queue-div").style.visibility = "visible"
    connect()
    addOnMessageReceiveListener(function (message) {
        if (message.action === actionJoined) {
            document.getElementById("queue-button").textContent = "Leave"
        } else if (message.action === actionLeft) {
            document.getElementById("queue-button").textContent = "Join"
        } else if (message.action === actionMatched) {
            gp = message.data[dataPlayer]
            document.getElementById("queue-div").style.visibility = "hidden"
            document.getElementById("board-div").style.visibility = "visible"
        } else if (message.action === actionMoved) {
            const x = message.data[dataMoveX]
            const y = message.data[dataMoveY]
            const p = message.data[dataPlayer]
            const s = message.data[dataScore]
            let block = gameBoard.blocks[y][x]
            if (p === gp) {
                block.s = 1
                fillBlock(gameBoard, block, yourColor, 3)
                yourScore.innerText = "Score: " + s
            } else {
                block.s = 2
                fillBlock(gameBoard, block, opponentColor, 3)
                opponentScore.innerText = "Score: " + s
            }
        }
    })
})

document.getElementById("queue-div").addEventListener("click", function (_) {
    if (document.getElementById("queue-button").textContent === "Join") {
        sendMessage(actionJoinMessage)
    } else if (document.getElementById("queue-button").textContent === "Leave") {
        sendMessage(actionLeaveMessage)
    }
})

const canvas = document.getElementById("game-board")
const width = 20
const height = 15
const gameBoard = new GameBoard(canvas, width, height)
console.log(gameBoard.blocks)
drawGrid(gameBoard)
addGridActionListener(gameBoard, clicked)

function clicked(gameGrid, block) {
    if (block.s === 0) {
        let message = new Message(actionMove)
        message.data[dataMoveX] = block.x
        message.data[dataMoveY] = block.y
        sendMessage(message)
    }
}

const yourScore = document.getElementById("your-score")
const opponentScore = document.getElementById("opponent-score")

yourScore.innerText="Score: 0"
opponentScore.innerText="Score: 0"
