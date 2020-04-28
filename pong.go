package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// TODO
// Display message in the screen like "press spacebar", or "you win"
// Add bonuses
// Change graphics - can ball make this motion "shadow" effect?
// PvP ?

const windowWidth = 800
const windowHeight = 600
const paddleConvexityEffectMultiplier = 150
const paddleVelocityEffectMultiplier = 15
const velocityAfterBounceMultiplier = 1.01
const initialBallXV = 400
const paddlePixelsRangeForCalculation = 4

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

func lerp(a float32, b float32, pct float32) float32 {
	return a + pct*(b-a)
}

type color struct {
	r, g, b byte
}

// The position is relative to the left upper conrner of the screen
type pos struct {
	x, y float32
}

// Returns center of the screen
func getCenter() pos {
	return pos{float32(windowWidth) / 2, float32(windowHeight) / 2}
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

	sdl.JoystickEventState(sdl.ENABLE)

	window, err := sdl.CreateWindow("PONG", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
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

	var joystickHandlers []*sdl.Joystick
	for i := 0; i < sdl.NumJoysticks(); i++ {
		joystickHandlers = append(joystickHandlers, sdl.JoystickOpen(i))
		defer joystickHandlers[i].Close()
	}

	pixels := make([]byte, windowWidth*windowHeight*4)

	player1 := paddle{pos{50, windowHeight / 2}, 20, 100, 300, 0, color{255, 255, 255}, 0}
	player2 := paddle{pos{windowWidth - 50, windowHeight / 2}, 20, 100, 300, 0, color{255, 255, 255}, 0}
	ball := ball{pos{300, 300}, 20, initialBallXV, 0, color{255, 255, 255}}

	keyState := sdl.GetKeyboardState()

	var (
		frameStart  time.Time
		elapsedTime float32
	)
	var joystickAxis int16

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
		for _, joystick := range joystickHandlers {
			if joystick != nil {
				joystickAxis = joystick.Axis(sdl.CONTROLLER_AXIS_LEFTY)
			}
		}
		if state == play {
			player1.update(keyState, joystickAxis, elapsedTime)
			player2.aiUpdate(&ball, elapsedTime)
			ball.update(&player1, &player2, elapsedTime)
		} else if state == start {
			if keyState[sdl.SCANCODE_SPACE] != 0 {
				if player1.score == 3 || player2.score == 3 {
					player1.score = 0
					player2.score = 0
				}
				// reset ball speed, but keep the direction
				ball.yv = 0
				if ball.xv < 0 {
					ball.xv = -initialBallXV
				} else {
					ball.xv = initialBallXV
				}
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
