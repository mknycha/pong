package main

type ball struct {
	// pos    pos //this is composition, the x would be referred to as ball.pos.x
	pos    // this birngs one struct into another, this allows us to refer to ball.x. It copies all the functions too!
	radius float32
	xv     float32
	yv     float32
	tex    texture
}

func (ball *ball) draw(pixels []byte) {
	ball.tex.drawAlpha(ball.pos, pixels)
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
			ball.xv = -ball.xv * velocityAfterBounceMultiplier
			ball.x = leftPaddle.x + leftPaddle.w/2.0 + ball.radius
			// handle bouncing angles differently closer the paddle edges (to the outside)
			ball.yv += (ball.y - leftPaddle.y) * elapsedTime * paddleConvexityEffectMultiplier
			// pass velocity of the moving paddle to the bar
			ball.yv += leftPaddle.yv * elapsedTime * paddleVelocityEffectMultiplier
		}
	}
	if ball.x+ball.radius > rightPaddle.x-rightPaddle.w/2 {
		if ball.y < rightPaddle.y+rightPaddle.h/2 && ball.y > rightPaddle.y-rightPaddle.h/2 {
			ball.xv = -ball.xv * velocityAfterBounceMultiplier
			ball.x = rightPaddle.x - rightPaddle.w/2.0 - ball.radius
			// handle bouncing angles differently closer the paddle edges (to the outside)
			ball.yv += (ball.y - rightPaddle.y) * elapsedTime * paddleConvexityEffectMultiplier
			// pass velocity of the moving paddle to the bar
			ball.yv += rightPaddle.yv * elapsedTime * paddleVelocityEffectMultiplier
		}
	}
}
