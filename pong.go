package main

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// TODO
// Frame rate independence
// Score
// Game over rate
// PvP ?
// Ai needs to be more imperfect

const windowWidth = 800
const windowHeight = 600
const thickness = 5

type color struct {
	r, g, b byte
}

// The position is relative to the left upper conrner of the screen
type pos struct {
	x, y float32
}

type score struct {
	pos
	h      float32
	w      float32
	num    int
	scored func() bool
}

func (score *score) top(pixels []byte) {
	for y := score.y - score.h/2; y < score.y-score.h/2+thickness; y++ {
		for x := score.x - score.w/2; x < score.x+score.w/2; x++ {
			setPixel(int(x), int(y), color{255, 255, 255}, pixels)
		}
	}
}

func (score *score) bottom(pixels []byte) {
	for y := score.y + score.h/2 - thickness; y < score.y+score.h/2; y++ {
		for x := score.x - score.w/2; x < score.x+score.w/2; x++ {
			setPixel(int(x), int(y), color{255, 255, 255}, pixels)
		}
	}
}

func (score *score) upperLeft(pixels []byte) {
	for y := score.y - score.h/2; y < score.y; y++ {
		for x := score.x - score.w/2; x < score.x-score.w/2+thickness; x++ {
			setPixel(int(x), int(y), color{255, 255, 255}, pixels)
		}
	}
}
func (score *score) lowerLeft(pixels []byte) {
	for y := score.y; y < score.y+score.h/2; y++ {
		for x := score.x - score.w/2; x < score.x-score.w/2+thickness; x++ {
			setPixel(int(x), int(y), color{255, 255, 255}, pixels)
		}
	}
}
func (score *score) upperRight(pixels []byte) {
	for y := score.y - score.h/2; y < score.y; y++ {
		for x := score.x + score.w/2 - thickness; x < score.x+score.w/2; x++ {
			setPixel(int(x), int(y), color{255, 255, 255}, pixels)
		}
	}
}
func (score *score) lowerRight(pixels []byte) {
	for y := score.y; y < score.y+score.h/2; y++ {
		for x := score.x + score.w/2 - thickness; x < score.x+score.w/2; x++ {
			setPixel(int(x), int(y), color{255, 255, 255}, pixels)
		}
	}
}

func (score *score) middle(pixels []byte) {
	for y := score.y - thickness/2; y < score.y+thickness/2; y++ {
		for x := score.x - score.w/2; x < score.x+score.w/2; x++ {
			setPixel(int(x), int(y), color{255, 255, 255}, pixels)
		}
	}
}

func (score *score) draw(pixels []byte) {
	switch score.num {
	case 0:
		score.bottom(pixels)
		score.lowerLeft(pixels)
		score.upperLeft(pixels)
		score.lowerRight(pixels)
		score.upperRight(pixels)
		score.top(pixels)
	case 1:
		score.lowerRight(pixels)
		score.upperRight(pixels)
	case 2:
		score.bottom(pixels)
		score.lowerLeft(pixels)
		score.upperRight(pixels)
		score.top(pixels)
		score.middle(pixels)
	case 3:
		score.bottom(pixels)
		score.lowerRight(pixels)
		score.upperRight(pixels)
		score.top(pixels)
		score.middle(pixels)
	case 4:
		score.upperLeft(pixels)
		score.lowerRight(pixels)
		score.upperRight(pixels)
		score.middle(pixels)
	case 5:
		score.top(pixels)
		score.upperLeft(pixels)
		score.lowerRight(pixels)
		score.middle(pixels)
		score.bottom(pixels)
	case 6:
		score.top(pixels)
		score.upperLeft(pixels)
		score.lowerRight(pixels)
		score.middle(pixels)
		score.lowerLeft(pixels)
		score.bottom(pixels)
	case 7:
		score.top(pixels)
		score.upperRight(pixels)
		score.lowerRight(pixels)
	case 8:
		score.top(pixels)
		score.upperLeft(pixels)
		score.upperRight(pixels)
		score.lowerRight(pixels)
		score.middle(pixels)
		score.lowerLeft(pixels)
		score.bottom(pixels)
	case 9:
		score.top(pixels)
		score.upperLeft(pixels)
		score.upperRight(pixels)
		score.lowerRight(pixels)
		score.middle(pixels)
		score.bottom(pixels)
	}
}

func (score *score) update(ball *ball) {
	if score.scored() {
		score.num++
		if score.num > 9 {
			score.num = 0
		}
	}
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

	if ball.y-ball.radius < 0 || ball.y+ball.radius > windowHeight { // bounce from the top or bottom of the screen
		ball.yv = -ball.yv
	}

	if ball.x-ball.radius < 0 || ball.x+ball.radius > windowWidth {
		ball.pos = getCenter()
	}

	if ball.x-ball.radius < leftPaddle.x+leftPaddle.w/2 {
		if ball.y < leftPaddle.y+leftPaddle.h/2 && ball.y > leftPaddle.y-leftPaddle.h/2 {
			ball.xv = -ball.xv
		}
	}
	if ball.x+ball.radius > rightPaddle.x-rightPaddle.w/2 {
		if ball.y < rightPaddle.y+rightPaddle.h/2 && ball.y > rightPaddle.y-rightPaddle.h/2 {
			ball.xv = -ball.xv
		}
	}
}

type paddle struct {
	pos
	w     float32
	h     float32
	speed float32
	color color
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
}

func (paddle *paddle) update(keyState []uint8, elapsedTime float32) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		paddle.y -= paddle.speed * elapsedTime
	} else if keyState[sdl.SCANCODE_DOWN] != 0 {
		paddle.y += paddle.speed * elapsedTime
	}
}

func (paddle *paddle) aiUpdate(ball *ball, elapsedTime float32) {
	paddle.y = ball.y
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

	pixels := make([]byte, windowWidth*windowHeight*4)

	// go func() {
	// 	sdl.Delay(5000)
	// 	e := sdl.QuitEvent{Type: sdl.QUIT}
	// 	sdl.PushEvent(&e)
	// }()

	player1 := paddle{pos{50, 100}, 20, 100, 300, color{255, 255, 255}}
	player2 := paddle{pos{windowWidth - 50, 100}, 20, 100, 300, color{255, 255, 255}}
	ball := ball{pos{300, 300}, 20, 400, 400, color{255, 255, 255}}
	player1Score := score{pos{180, 100}, 70, 40, 0, func() bool {
		return ball.x+ball.radius > windowWidth
	}}
	player2Score := score{pos{windowWidth - 180, 100}, 70, 40, 0, func() bool {
		return ball.x-ball.radius < 0
	}}

	keyState := sdl.GetKeyboardState()

	var (
		frameStart  time.Time
		elapsedTime float32
	)

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
		clear(pixels)

		player1Score.update(&ball)
		player2Score.update(&ball)
		player1.update(keyState, elapsedTime)
		player2.aiUpdate(&ball, elapsedTime)
		ball.update(&player1, &player2, elapsedTime)

		player1.draw(pixels)
		player2.draw(pixels)
		ball.draw(pixels)
		player1Score.draw(pixels)
		player2Score.draw(pixels)

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
