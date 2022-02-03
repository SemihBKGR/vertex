document.getElementById("board-div").style.visibility = "hidden"
document.getElementById("queue-div").style.visibility = "hidden"

document.getElementById("start-button").addEventListener("click", function (_) {
    console.log(gameBoard.blocks)
    document.getElementById("start-div").style.visibility = "hidden"
    document.getElementById("queue-div").style.visibility = "visible"
    connect()
    addOnMessageReceiveListener(function (message) {
        console.log(gameBoard.blocks)
        console.log(message)
        if (message.action === actionJoined) {
            document.getElementById("queue-button").textContent = "Leave"
        } else if (message.action === actionLeft) {
            document.getElementById("queue-button").textContent = "Join"
        } else if (message.action === actionMatched) {
            document.getElementById("queue-div").style.visibility = "hidden"
            document.getElementById("board-div").style.visibility = "visible"
        } else if(message.action===actionMoved){
            const x=message.data["x"]
            const y=message.data["y"]
            const p=message.data["p"]
            console.log(gameBoard.blocks)
            console.log(gameBoard.blocks[0].length)
            let block=gameBoard.blocks[y][x]
            if (!p){
                block.s=1
                fillBlock(gameBoard,block,"green",3)
            }else{
                block.s=2
                fillBlock(gameBoard,block,"blue",3)
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
    let message = new Message(actionMove)
    message.data["x"] = block.x
    message.data["y"] = block.y
    sendMessage(message)
}
