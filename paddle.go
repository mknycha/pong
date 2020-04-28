package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type paddle struct {
	pos
	w     float32
	h     float32
	speed float32
	score int
	color color
	yv    float32
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
	if keyState[sdl.SCANCODE_UP] != 0 && paddle.y > 0 {
		paddle.y -= paddle.speed * elapsedTime
		paddle.yv = -paddle.speed
	} else if keyState[sdl.SCANCODE_DOWN] != 0 && paddle.y < windowHeight {
		paddle.y += paddle.speed * elapsedTime
		paddle.yv = paddle.speed
	} else {
		paddle.yv = 0
	}

	if math.Abs(float64(controllerAxis)) > 1500 && paddle.y < windowHeight && paddle.y > 0 {
		pct := float32(controllerAxis) / 32767.0
		paddle.y += paddle.speed * pct * elapsedTime
	}
}

func (paddle *paddle) aiUpdate(ball *ball, elapsedTime float32) {
	if (paddle.x - ball.x) < (float32(windowWidth) * 3 / 4) { // ball is close enough to be "seen"
		// paddlePixelsRangeForCalculation is used so that paddle is not moved with pixel precision
		if (paddle.y + paddlePixelsRangeForCalculation) < ball.y { // ball is above, paddle moves up
			paddle.y += paddle.speed * elapsedTime
			paddle.yv = paddle.speed
		} else if (paddle.y - paddlePixelsRangeForCalculation) > ball.y { // ball is below, paddle moves down
			paddle.y -= paddle.speed * elapsedTime
			paddle.yv = -paddle.speed
		} else {
			paddle.yv = 0
		}
	}
}
