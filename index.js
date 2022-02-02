class GameGrid {
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

function drawGrid(gameGrid) {
    gameGrid.context.beginPath()
    for (let i = gameGrid.blockWidth; i < gameGrid.canvas.width; i += gameGrid.blockWidth) {
        gameGrid.context.moveTo(i, 0)
        gameGrid.context.lineTo(i, gameGrid.canvas.height)
    }
    for (let i = gameGrid.blockHeight; i < gameGrid.canvas.height; i += gameGrid.blockHeight) {
        gameGrid.context.moveTo(0, i)
        gameGrid.context.lineTo(gameGrid.canvas.width, i)
    }
    gameGrid.context.stroke()
}

function addGridActionListener(gameGrid, clickCallback) {
    let hoverBlock = undefined
    gameGrid.canvas.addEventListener("mousemove", function (event) {
        if (hoverBlock !== undefined && hoverBlock.s === 0) {
            fillBlock(gameGrid, hoverBlock, "white")
        }
        drawGrid(gameGrid)
        let mousePosition = getMousePosition(gameGrid.canvas, event)
        hoverBlock = getBlock(gameGrid, mousePosition)
        if (hoverBlock !== undefined && hoverBlock.s === 0) {
            fillBlock(gameGrid, hoverBlock, "black")
        }
    })
    gameGrid.canvas.addEventListener("click", function (_) {
        let block = gameGrid.blockMatrix[hoverBlock.y][hoverBlock.x]
        if (block !== undefined && block.s === 0) {
            block.s = 1
            clickCallback(gameGrid, block)
        }
    })
    gameGrid.canvas.addEventListener("mouseleave", function (_) {
        hoverBlock = undefined
    })
}

function getBlock(gameGrid, mousePosition) {
    let x = Math.floor(mousePosition.x / gameGrid.blockWidth)
    let y = Math.floor(mousePosition.y / gameGrid.blockHeight)
    if (gameGrid.width > x && gameGrid.height > y) {
        return gameGrid.blockMatrix[y][x]
    }
    return undefined
}

function fillBlock(gameGrid, block, fillStyle) {
    gameGrid.context.fillStyle = fillStyle
    gameGrid.context.fillRect(block.x * gameGrid.blockWidth, block.y * gameGrid.blockHeight, gameGrid.blockWidth, gameGrid.blockHeight)
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
const gameGrid = new GameGrid(canvas, width, height)
drawGrid(gameGrid)
addGridActionListener(gameGrid, clicked)

function clicked(gameGrid, block) {
    console.log(block)
}

// -------------------------------------------------------
