package main

import (
	"e7_pong/noise"
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

func flerp(b1 byte, b2 byte, pct float32) byte {
	return byte(float32(b1) + pct*(float32(b2)-float32(b1)))
}

func colorLerp(c1, c2 color, pct float32) color {
	return color{flerp(c1.r, c2.r, pct), flerp(c1.g, c2.g, pct), flerp(c1.b, c2.b, pct)}
}

func clamp(min, max, v int) int {
	if v < min {
		v = min
	} else if v > max {
		v = max
	}
	return v
}

func getDualGradient(c1, c2, c3, c4 color) []color {
	result := make([]color, 256)
	for i := range result {
		pct := float32(i) / float32(255)
		if pct < 0.5 {
			result[i] = colorLerp(c1, c2, pct*float32(2))
		} else {
			result[i] = colorLerp(c3, c4, pct*float32(1.5)*float32(0.5))
		}
	}
	return result
}

func getGradient(c1, c2 color) []color {
	result := make([]color, 256)
	for i := range result {
		pct := float32(i) / float32(255)
		result[i] = colorLerp(c1, c2, pct)
	}
	return result
}

func rescaleAndDraw(noise []float32, min, max float32, gradient []color, w, h int) []byte {
	result := make([]byte, w*h*4)
	scale := 255.0 / (max - min)
	offset := min * scale

	for i := range noise {
		noise[i] = noise[i]*scale - offset
		c := gradient[clamp(0, 255, int(noise[i]))]
		p := i * 4
		result[p] = c.r
		result[p+1] = c.g
		result[p+2] = c.b
	}
	return result
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

	player1Texture := loadTexture("assets/pyramid_left.png")
	player2Texture := loadTexture("assets/pyramid_right.png")
	player1 := paddle{pos{50, windowHeight / 2}, float32(player1Texture.w), float32(player2Texture.h), 300, 0, 0, color{255, 255, 255}, player1Texture}
	player2 := paddle{pos{windowWidth - 50, windowHeight / 2}, float32(player1Texture.w), float32(player2Texture.h), 300, 0, 0, color{255, 255, 255}, player2Texture}
	ballTexture := loadTexture("assets/ball.png")
	ball := ball{pos{300, 300}, float32(ballTexture.w / 2), initialBallXV, 0, ballTexture}

	keyState := sdl.GetKeyboardState()

	noise, min, max := noise.MakeNoise(noise.TURBULENCE, 0.01, 0.2, 2, 3, windowWidth, windowHeight)
	gradient := getGradient(color{252, 193, 0}, color{253, 120, 0})
	noisePixels := rescaleAndDraw(noise, min, max, gradient, windowWidth, windowHeight)

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

		for i := range noisePixels {
			pixels[i] = noisePixels[i]
		}
		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)

		tex.Update(nil, pixels, windowWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		elapsedTime = float32(time.Since(frameStart).Seconds())
		if elapsedTime < 0.005 {
			sdl.Delay(5 - uint32(elapsedTime*1000.0))
			elapsedTime = float32(time.Since(frameStart).Seconds())
		}
	}
}
