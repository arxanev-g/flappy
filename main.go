package main

import (
	"math"
	"math/rand"
	"strconv"
	"syscall/js"
	"time"

	"honnef.co/go/js/dom/v2"
)

var (
	window  dom.Window
	ctx     js.Value
	birdImg js.Value

	width  float64
	height float64

	gravity = 0.5
	jump    = -10.0

	bird              *Bird
	pipes             []*Pipe
	score             int
	gameOver          bool
	imgLoaded         = make(chan struct{})
	pipeSpawnInterval = 150
	frameCounter      = 0
)

func main() {
	window = dom.GetWindow()
	document := js.Global().Get("document")
	body := document.Call("getElementsByTagName", "body").Index(0)
	str := `<canvas id="myCanvas" width=` + strconv.Itoa(window.InnerWidth()) + " height=" + strconv.Itoa(window.InnerHeight()) + `></canvas>`
	body.Set("innerHTML", str)

	rand.Seed(time.Now().UnixNano())

	updateSizes()
	bird = NewBird(width, height)

	canvas := document.Call("getElementById", "myCanvas")
	ctx = js.Global().Get("document").
		Call("getElementById", "myCanvas").
		Call("getContext", "2d")

	loadBirdImage()
	<-imgLoaded

	window.AddEventListener("resize", false, func(dom.Event) {
		updateSizes()
		updateCanvasSize(&canvas)
	})

	setupEventListeners()

	var renderFrame js.Func
	renderFrame = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if !gameOver {
			update()
		}
		draw()
		js.Global().Call("requestAnimationFrame", renderFrame)
		return nil
	})
	js.Global().Call("requestAnimationFrame", renderFrame)

	select {}
}

func updateSizes() {
	width = float64(window.InnerWidth())
	height = float64(window.InnerHeight())
}

func updateCanvasSize(canvas *js.Value) {
	canvas.Set("Width", (int(width)))
	canvas.Set("Height", (int(height)))
}

func loadBirdImage() {
	birdImg = js.Global().Get("Image").New()
	birdImg.Set("onload", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		close(imgLoaded)
		return nil
	}))
	birdImg.Set("src", "gopher.png")
}

func setupEventListeners() {
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if gameOver {
			resetGame()
			return nil
		}
		bird.Jump(jump * (height / 600))
		return nil
	})

	js.Global().Get("document").Call("addEventListener", "click", handler)
	js.Global().Get("document").Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if args[0].Get("code").String() == "Space" {
			handler.Invoke()
		}
		return nil
	}))
}

func update() {

	frameCounter++

	bird.Update(gravity*(height/600), 1.0)
	if frameCounter >= pipeSpawnInterval {
		pipes = append(pipes, NewPipe(width, height, width*0.003))
		frameCounter = 0
	}

	for i := 0; i < len(pipes); i++ {
		pipes[i].Update()

		if pipes[i].IsOffscreen() {
			pipes = append(pipes[:i], pipes[i+1:]...)
			i--
			score++
			continue
		}
		if !gameOver && pipes[i].CheckCollision(bird) {
			gameOver = true
		}
	}

	if bird.Y < 0 || bird.Y+bird.Height > height {
		gameOver = true
	}
}

func draw() {
	ctx.Call("clearRect", 0, 0, width, height)

	if !gameOver {
		// Рисуем игровые объекты
		ctx.Call("drawImage", birdImg, bird.X, bird.Y, bird.Width, bird.Height)

		// Рисуем трубы
		for _, pipe := range pipes {
			pipe.Draw(ctx)
		}

		// Рисуем текущий счет
		fontSize := math.Max(width*0.08, 28)
		ctx.Set("font", js.ValueOf(fontSize).String()+"px JetBrains Mono")
		ctx.Set("fillStyle", "black")
		ctx.Call("fillText", score, width*0.05, height*0.1)
	} else {
		// Рисуем только экран Game Over
		gameOverFontSize := math.Max(width*0.12, 40)
		ctx.Set("font", js.ValueOf(gameOverFontSize).String()+"px Arial")
		ctx.Set("fillStyle", "red")

		// Центрируем текст
		textX := width / 2
		textY := height / 2

		// Основной текст
		ctx.Call("fillText", "GAME OVER", textX, textY)
		ctx.Set("font", js.ValueOf(gameOverFontSize*0.6).String()+"px Arial")
		ctx.Call("fillText", "Your account "+strconv.Itoa(score), textX, textY+height*0.1)
		// Инструкция
		ctx.Set("font", js.ValueOf(gameOverFontSize*0.6).String()+"px Arial")
		ctx.Call("fillText", "Click to restart", textX, textY+height*0.15)
	}
}

func resetGame() {
	bird.Y = height / 2
	bird.VY = 0
	pipes = nil
	score = 0
	gameOver = false
}
