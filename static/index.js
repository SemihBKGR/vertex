const blockYourColor = "#4961ec"
const blockOpponentColor = "#e6674e"
const blockWallColor = "#421f17"

const titleColorTurn = "#5a9d48"
const titleColorWait = "#343333"

let gp
let tp

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
            let walls = message.data[dataWalls]
            console.log(walls)
            for (const wall of walls) {
                const block = gameBoard.blocks[wall.y][wall.x]
                block.s = 1
                fillBlock(gameBoard, block, blockWallColor, 3)
            }
            gp = message.data[dataPlayer]
            if (!gp) {
                yourTitle.style.color = titleColorTurn
            } else {
                opponentTitle.style.color = titleColorTurn
            }
            tp = false
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
                fillBlock(gameBoard, block, blockYourColor, 3)
                yourScore.innerText = "Score: " + s
                yourTitle.style.color = titleColorWait
                opponentTitle.style.color = titleColorTurn
            } else {
                block.s = 2
                fillBlock(gameBoard, block, blockOpponentColor, 3)
                opponentScore.innerText = "Score: " + s
                opponentTitle.style.color = titleColorWait
                yourTitle.style.color = titleColorTurn
            }
            tp = !tp
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
    if (gp === tp) {
        if (block.s === 0) {
            let message = new Message(actionMove)
            message.data[dataMoveX] = block.x
            message.data[dataMoveY] = block.y
            sendMessage(message)
        }
    }
}

const yourTitle = document.getElementById("your-title")
const opponentTitle = document.getElementById("opponent-title")
yourTitle.style.color = titleColorWait
opponentTitle.style.color = titleColorWait

const yourScore = document.getElementById("your-score")
const opponentScore = document.getElementById("opponent-score")
yourScore.innerText = "Score: 0"
opponentScore.innerText = "Score: 0"
