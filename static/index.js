class GameBoard {
    canvas
    context
    width
    height
    blockMatrix
    blockWidth
    blockHeight
    constructor(canvas, width, height) {
        this.canvas = canvas
        this.width = width
        this.height = height
        this.context = canvas.getContext("2d")
        this.blockMatrix = []
        this.blockWidth = canvas.width / width
        this.blockHeight = canvas.height / height
        for (let i = 0; i < height; i++) {
            let blockArray = []
            for (let j = 0; j < width; j++) {
                blockArray.push(new Block(j, i, 0))
            }
            this.blockMatrix.push(blockArray)
        }
    }
}

class Block {
    x
    y
    s
    constructor(x, y, s) {
        this.x = x;
        this.y = y;
        this.s = s;
    }
}

function drawGrid(gameBoard) {
    gameBoard.context.beginPath()
    for (let i = gameBoard.blockWidth; i < gameBoard.canvas.width; i += gameBoard.blockWidth) {
        gameBoard.context.moveTo(i, 0)
        gameBoard.context.lineTo(i, gameBoard.canvas.height)
    }
    for (let i = gameBoard.blockHeight; i < gameBoard.canvas.height; i += gameBoard.blockHeight) {
        gameBoard.context.moveTo(0, i)
        gameBoard.context.lineTo(gameBoard.canvas.width, i)
    }
    gameBoard.context.stroke()
}

function addGridActionListener(gameBoard, clickCallback) {
    let hoverBlock = undefined
    gameBoard.canvas.addEventListener("mousemove", function (event) {
        if (hoverBlock !== undefined && hoverBlock.s === 0) {
            fillBlock(gameBoard, hoverBlock, "white")
        }
        drawGrid(gameBoard)
        let mousePosition = getMousePosition(gameBoard.canvas, event)
        hoverBlock = getBlock(gameBoard, mousePosition)
        if (hoverBlock !== undefined && hoverBlock.s === 0) {
            fillBlock(gameBoard, hoverBlock, "black")
        }
    })
    gameBoard.canvas.addEventListener("click", function (_) {
        let block = gameBoard.blockMatrix[hoverBlock.y][hoverBlock.x]
        if (block !== undefined && block.s === 0) {
            block.s = 1
            clickCallback(gameBoard, block)
        }
    })
    gameBoard.canvas.addEventListener("mouseleave", function (_) {
        hoverBlock = undefined
    })
}

function getBlock(gameBoard, mousePosition) {
    let x = Math.floor(mousePosition.x / gameBoard.blockWidth)
    let y = Math.floor(mousePosition.y / gameBoard.blockHeight)
    if (gameBoard.width > x && gameBoard.height > y) {
        return gameBoard.blockMatrix[y][x]
    }
    return undefined
}

function fillBlock(gameBoard, block, fillStyle) {
    gameBoard.context.fillStyle = fillStyle
    gameBoard.context.fillRect(block.x * gameBoard.blockWidth, block.y * gameBoard.blockHeight, gameBoard.blockWidth, gameBoard.blockHeight)
}

function getMousePosition(canvas, event) {
    let rect = canvas.getBoundingClientRect();
    return {
        x: (event.clientX - rect.left) / (rect.right - rect.left) * canvas.width,
        y: (event.clientY - rect.top) / (rect.bottom - rect.top) * canvas.height
    };
}

// -------------------------------------------------------
const canvas = document.getElementById("game-canvas")
const width = 20
const height = 15
const gameBoard = new GameBoard(canvas, width, height)
drawGrid(gameBoard)
addGridActionListener(gameBoard, clicked)

function clicked(gameGrid, block) {
    console.log(block)
}

// -------------------------------------------------------
