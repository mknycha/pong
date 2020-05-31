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
	yv    float32
	color color
	tex   texture
}

func (paddle *paddle) draw(pixels []byte) {
	paddle.tex.drawAlpha(pos{paddle.x, paddle.y}, pixels)

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
