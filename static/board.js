class GameBoard {
    canvas
    context
    width
    height
    blocks
    blockWidth
    blockHeight

    constructor(canvas, width, height) {
        this.canvas = canvas
        this.width = width
        this.height = height
        this.context = canvas.getContext("2d")
        this.blocks = []
        this.blockWidth = canvas.width / width
        this.blockHeight = canvas.height / height
        for (let i = 0; i < height; i++) {
            let blockArray = []
            for (let j = 0; j < width; j++) {
                blockArray.push(new Block(j, i, 0))
            }
            this.blocks.push(blockArray)
        }
    }
}

class Block {
    x
    y
    s
    m

    constructor(x, y, s) {
        this.x = x;
        this.y = y;
        this.s = s;
        this.m = false
    }
}

class MousePosition {
    x
    y

    constructor(x, y) {
        this.x = x
        this.y = y
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

function addGridActionListener(gameBoard, blockHoverColor, blockDefaultColor, blockHighlightedColor, clickCallback) {
    let hoverBlock = undefined
    gameBoard.canvas.addEventListener("mousemove", function (event) {
        if (hoverBlock !== undefined && hoverBlock.s === 0) {
            if (hoverBlock.m) {
                fillBlock(gameBoard, hoverBlock, blockHighlightedColor, 3)
            } else {
                fillBlock(gameBoard, hoverBlock, blockDefaultColor, 3)
            }
        }
        let mousePosition = getMousePosition(gameBoard.canvas, event)
        hoverBlock = getBlock(gameBoard, mousePosition)
        if (hoverBlock !== undefined && hoverBlock.s === 0) {
            fillBlock(gameBoard, hoverBlock, blockHoverColor, 3)
        }
    })
    gameBoard.canvas.addEventListener("click", function (_) {
        let block = gameBoard.blocks[hoverBlock.y][hoverBlock.x]
        if (block !== undefined) {
            clickCallback(gameBoard, block)
        }
    })
    gameBoard.canvas.addEventListener("mouseleave", function (_) {
        if (hoverBlock !== undefined && hoverBlock.s === 0) {
            fillBlock(gameBoard, hoverBlock, "white", 3)
        }
        hoverBlock = undefined
    })
}

function getBlock(gameBoard, mousePosition) {
    let x = Math.floor(mousePosition.x / gameBoard.blockWidth)
    let y = Math.floor(mousePosition.y / gameBoard.blockHeight)
    if (gameBoard.width > x && gameBoard.height > y) {
        return gameBoard.blocks[y][x]
    }
    return undefined
}

function fillBlock(gameBoard, block, fillStyle, margin) {
    gameBoard.context.fillStyle = fillStyle
    gameBoard.context.fillRect(
        block.x * gameBoard.blockWidth + margin,
        block.y * gameBoard.blockHeight + margin,
        gameBoard.blockWidth - (margin * 2),
        gameBoard.blockHeight - (margin * 2))
}

function fillMarkedBlocks(gameBoard, fillStyle, margin) {
    for (let i = 0; i < gameBoard.height; i++) {
        for (let j = 0; j < gameBoard.width; j++) {
            const block = gameBoard.blocks[i][j]
            if (block.m && block.s === 0) {
                fillBlock(gameBoard, block, fillStyle, margin)
            }
        }
    }
}

function getMousePosition(canvas, event) {
    let rect = canvas.getBoundingClientRect();
    let x = (event.clientX - rect.left) / (rect.right - rect.left) * canvas.width
    let y = (event.clientY - rect.top) / (rect.bottom - rect.top) * canvas.height
    return new MousePosition(x, y)
}

function markMovePossibilities(gameBoard, player) {
    for (let i = 0; i < gameBoard.height; i++) {
        for (let j = 0; j < gameBoard.width; j++) {
            let block = gameBoard.blocks[i][j]
            if ((block.s === 1 && !player) || (block.s === 2 && player)) {
                let blockY = i
                while (--blockY >= 0 && gameBoard.blocks[blockY][j].s === 0) {
                    gameBoard.blocks[blockY][j].m = true
                }
                blockY = i
                while (++blockY < gameBoard.height && gameBoard.blocks[blockY][j].s === 0) {
                    gameBoard.blocks[blockY][j].m = true
                }
                let blockX = j
                while (--blockX >= 0 && gameBoard.blocks[i][blockX].s === 0) {
                    gameBoard.blocks[i][blockX].m = true
                }
                blockX = j
                while (++blockX < gameBoard.width && gameBoard.blocks[i][blockX].s === 0) {
                    gameBoard.blocks[i][blockX].m = true
                }
            }
        }
    }
}

function unmarkMovePossibilities(gameBoard) {
    for (let i = 0; i < gameBoard.height; i++) {
        for (let j = 0; j < gameBoard.width; j++) {
            gameBoard.blocks[i][j].m = false
        }
    }
}
