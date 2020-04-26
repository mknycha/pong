package main

import (
	"fmt"
	"math"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// TODO
// Make the ball go faster and faster with each bounce
// Make the paddle unable to go out of the screen
// PvP ?

const windowWidth = 800
const windowHeight = 600
const paddleConvexityEffectMultiplier = 150
const paddleVelocityEffectMultiplier = 15

// This kind of enum in GO
type gameState int

const (
	start gameState = iota
	play
)

// till here
var state = start

var nums = [][]byte{
	{
		1, 1, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 1, 1,
	},
	{
		1, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		1, 1, 1,
		1, 0, 0,
		1, 1, 1,
	},
	{
		1, 1, 1,
		0, 0, 1,
		0, 1, 1,
		0, 0, 1,
		1, 1, 1,
	},
}

func drawNumber(pos pos, color color, size int, num int, pixels []byte) {
	startX := int(pos.x) - (size*3)/2
	startY := int(pos.y) - (size*5)/2

	for i, v := range nums[num] {
		if v == 1 {
			for y := startY; y < startY+size; y++ {
				for x := startX; x < startX+size; x++ {
					setPixel(x, y, color, pixels)
				}
			}
		}
		startX += size
		// Once we have drawn three squares, we are going down by one square and left by three squares
		if (i+1)%3 == 0 {
			startY += size
			startX -= size * 3
		}
	}
}

type color struct {
	r, g, b byte
}

// The position is relative to the left upper conrner of the screen
type pos struct {
	x, y float32
}

type ball struct {
	// pos    pos //this is composition, the x would be referred to as ball.pos.x
	pos    // this birngs one struct into another, this allows us to refer to ball.x. It copies all the functions too!
	radius float32
	xv     float32
	yv     float32
	color  color
}

func (ball *ball) draw(pixels []byte) {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x+x), int(ball.y+y), ball.color, pixels)
			}
		}
	}
}

// Returns center of the screen
func getCenter() pos {
	return pos{float32(windowWidth) / 2, float32(windowHeight) / 2}
}

func (ball *ball) update(leftPaddle *paddle, rightPaddle *paddle, elapsedTime float32) {
	ball.x += ball.xv * elapsedTime
	ball.y += ball.yv * elapsedTime

	if ball.y-ball.radius < 0 { // bounce from the bottom of the screen
		ball.yv = -ball.yv
		ball.y = ball.radius
	}
	if ball.y+ball.radius > windowHeight { // bounce from the top of the screen
		ball.yv = -ball.yv
		ball.y = windowHeight - ball.radius
	}

	if ball.x-ball.radius < 0 {
		rightPaddle.score++
		ball.pos = getCenter()
		state = start
	} else if ball.x+ball.radius > windowWidth {
		leftPaddle.score++
		ball.pos = getCenter()
		state = start
	}

	if ball.x-ball.radius < leftPaddle.x+leftPaddle.w/2 {
		if ball.y < leftPaddle.y+leftPaddle.h/2 && ball.y > leftPaddle.y-leftPaddle.h/2 {
			ball.xv = -ball.xv
			ball.x = leftPaddle.x + leftPaddle.w/2.0 + ball.radius
			// handle bouncing angles differently closer the paddle edges (to the outside)
			ball.yv += (ball.y - leftPaddle.y) * elapsedTime * paddleConvexityEffectMultiplier
			// pass velocity of the moving paddle to the bar
			ball.yv += leftPaddle.yv * elapsedTime * paddleVelocityEffectMultiplier
		}
	}
	if ball.x+ball.radius > rightPaddle.x-rightPaddle.w/2 {
		if ball.y < rightPaddle.y+rightPaddle.h/2 && ball.y > rightPaddle.y-rightPaddle.h/2 {
			ball.xv = -ball.xv
			ball.x = rightPaddle.x - rightPaddle.w/2.0 - ball.radius
			// handle bouncing angles differently closer the paddle edges (to the outside)
			ball.yv += (ball.y - rightPaddle.y) * elapsedTime * paddleConvexityEffectMultiplier
			// pass velocity of the moving paddle to the bar
			ball.yv += rightPaddle.yv * elapsedTime * paddleVelocityEffectMultiplier
		}
	}
}

type paddle struct {
	pos
	w     float32
	h     float32
	speed float32
	score int
	color color
	yv    float32
}

func lerp(a float32, b float32, pct float32) float32 {
	return a + pct*(b-a)
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x - paddle.w/2)
	startY := int(paddle.y - paddle.h/2)

	// There is a reason to start with y, because it uses ram cache
	// If we load to our RAM 0, 1, 2, 3, 4, 5, 6, 7, 8 we will go through order and be in cache
	// 0, 1, 2,
	// 3, 4, 5,
	// 6, 7, 8
	for y := 0; y < int(paddle.h); y++ {
		for x := 0; x < int(paddle.w); x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}

	numX := lerp(paddle.x, getCenter().x, 0.2)
	drawNumber(pos{numX, 35}, paddle.color, 10, paddle.score, pixels)
}

func (paddle *paddle) update(keyState []uint8, controllerAxis int16, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= paddle.speed * elapsedTime
		paddle.yv = -paddle.speed
	} else if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += paddle.speed * elapsedTime
		paddle.yv = paddle.speed
	} else {
		paddle.yv = 0
	}

	if math.Abs(float64(controllerAxis)) > 1500 {
		pct := float32(controllerAxis) / 32767.0
		paddle.y += paddle.speed * pct * elapsedTime
	}
}

func (paddle *paddle) aiUpdate(ball *ball, elapsedTime float32) {
	if paddle.y < ball.y { // ball is above, paddle moves up
		paddle.y += paddle.speed * elapsedTime
		paddle.yv = paddle.speed
	} else if paddle.y > ball.y { // ball is below, paddle moves down
		paddle.y -= paddle.speed * elapsedTime
		paddle.yv = -paddle.speed
	} else {
		paddle.yv = 0
	}
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*windowWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Testing SDL2", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(windowWidth), int32(windowHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(windowWidth), int32(windowHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	var controllerHandlers []*sdl.GameController
	for i := 0; i < sdl.NumJoysticks(); i++ {
		controllerHandlers = append(controllerHandlers, sdl.GameControllerOpen(i))
		defer controllerHandlers[i].Close()
	}

	pixels := make([]byte, windowWidth*windowHeight*4)

	// go func() {
	// 	sdl.Delay(5000)
	// 	e := sdl.QuitEvent{Type: sdl.QUIT}
	// 	sdl.PushEvent(&e)
	// }()

	player1 := paddle{pos{50, windowHeight / 2}, 20, 100, 300, 0, color{255, 255, 255}, 0}
	player2 := paddle{pos{windowWidth - 50, windowHeight / 2}, 20, 100, 300, 0, color{255, 255, 255}, 0}
	ball := ball{pos{300, 300}, 20, 400, 0, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	var (
		frameStart  time.Time
		elapsedTime float32
	)
	var controllerAxis int16

	running := true
	for running {
		frameStart = time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
		for _, controller := range controllerHandlers {
			if controller != nil {
				println(controller)
				controllerAxis = controller.Axis(sdl.CONTROLLER_AXIS_LEFTY)
			}
		}
		if state == play {
			player1.update(keyState, controllerAxis, elapsedTime)
			player2.aiUpdate(&ball, elapsedTime)
			ball.update(&player1, &player2, elapsedTime)
		} else if state == start {
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.score == 3 || player2.score == 3 {
					player1.score = 0
					player2.score = 0
				}
				ball.yv = 0
				player1.y = windowHeight / 2
				player2.y = windowHeight / 2
				state = play
			}
		}

		clear(pixels)
		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		tex.Update(nil, pixels, windowWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Seconds())
		if elapsedTime < 0.005 {
			sdl.Delay(5 - uint32(elapsedTime/1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
	}
}
