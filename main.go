package main

import (
	"image/color"
	"log"
	"math/rand"

	e "github.com/hajimehoshi/ebiten/v2"
)

type Snake struct {
	image *e.Image

	x int
	y int

	lenght int

	head bool
	tail bool

	direction string
}

type Eat struct {
	image *e.Image

	x int
	y int

	lenght int
}

type Window struct {
	// ширина
	Height int
	// высота
	Width int
}

type Game struct {
	window Window
	snake  []Snake
	eat    Eat

	rememberX int
	rememberY int
}

func (g *Game) Update() error {
	// borders
	if g.snake[0].y <= 0 {
		// g.snake[0].y = 0
		log.Fatal("bye")
	}

	if g.snake[0].x <= 0 {
		log.Fatal("bye")
		// g.snake[0].x = 0
	}

	if g.snake[0].y >= g.window.Height-g.snake[0].lenght {
		log.Fatal("bye")
		// g.snake[0].y = g.window.Height - g.snake[0].lenght
	}

	if g.snake[0].x >= g.window.Width-g.snake[0].lenght {
		log.Fatal("bye")
		// g.snake[0].x = g.window.Width - g.snake[0].lenght
	}
	// borders

	if e.IsKeyPressed(e.KeyRight) && g.snake[0].direction != "left" && g.snake[0].direction != "right" {
		g.snake[0].direction = "right"
		g.rememberX = g.snake[0].x
		g.rememberY = g.snake[0].y
	}
	if e.IsKeyPressed(e.KeyUp) && g.snake[0].direction != "down" && g.snake[0].direction != "up" {
		g.snake[0].direction = "up"
		g.rememberX = g.snake[0].x
		g.rememberY = g.snake[0].y
	}
	if e.IsKeyPressed(e.KeyDown) && g.snake[0].direction != "up" && g.snake[0].direction != "down" {
		g.snake[0].direction = "down"
		g.rememberX = g.snake[0].x
		g.rememberY = g.snake[0].y
	}
	if e.IsKeyPressed(e.KeyLeft) && g.snake[0].direction != "right" && g.snake[0].direction != "left" {
		g.snake[0].direction = "left"
		g.rememberX = g.snake[0].x
		g.rememberY = g.snake[0].y
	}

	for i := range g.snake {
		if g.snake[i].x == g.rememberX && g.snake[i].y == g.rememberY {
			g.snake[i].direction = g.snake[0].direction
		}
	}

	for i := range g.snake {
		switch g.snake[i].direction {
		case "right":
			g.snake[i].x += 10
		case "up":
			g.snake[i].y -= 10
		case "down":
			g.snake[i].y += 10
		case "left":
			g.snake[i].x -= 10
		}
	}

	if checkIntersection(g.eat.x, g.eat.y, g.eat.lenght, g.eat.lenght, g.snake[0].x, g.snake[0].y, g.snake[0].lenght, g.snake[0].lenght) {
		g.eat.x, g.eat.y = g.genEat()
		g.grow()
	}

	return nil
}

func (g *Game) genEat() (int, int) {
	needed := false
	var eatX int
	var eatY int
	for !needed {
		eatX = rand.Intn(g.window.Width)
		eatY = rand.Intn(g.window.Height)
		for i := range g.snake {
			if eatX == g.snake[i].x && eatY == g.snake[i].y {
				break
			}
		}
		needed = true
	}
	return eatX, eatY
}

func (g *Game) grow() {
	snake := Snake{}
	snake.lenght = 20
	snake.head = false
	snake.tail = true

	for i := range g.snake {
		if g.snake[i].tail {
			if g.snake[i].direction == "right" {
				snake.x = g.snake[i].x - snake.lenght
				snake.y = g.snake[i].y
			}
			if g.snake[i].direction == "left" {
				snake.x = g.snake[i].x + snake.lenght
				snake.y = g.snake[i].y
			}
			if g.snake[i].direction == "up" {
				snake.x = g.snake[i].x
				snake.y = g.snake[i].y + snake.lenght
			}
			if g.snake[i].direction == "down" {
				snake.x = g.snake[i].x
				snake.y = g.snake[i].y - snake.lenght
			}
			g.snake[i].tail = false
			snake.direction = g.snake[i].direction
		}
	}

	snakeIMG := e.NewImage(snake.lenght, snake.lenght)
	snakeIMG.Fill(color.RGBA{255, 10, 10, 1})
	snake.image = snakeIMG

	g.snake = append(g.snake, snake)
}

func checkIntersection(x1, y1, w1, h1, x2, y2, w2, h2 int) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

func (g *Game) Draw(screen *e.Image) {
	for _, block := range g.snake {
		snakeGeoM := e.GeoM{}
		snakeGeoM.Translate(float64(block.x), float64(block.y))
		screen.DrawImage(block.image, &e.DrawImageOptions{GeoM: snakeGeoM})
	}
	eatGeoM := e.GeoM{}
	eatGeoM.Translate(float64(g.eat.x), float64(g.eat.y))
	screen.DrawImage(g.eat.image, &e.DrawImageOptions{GeoM: eatGeoM})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.window.Width, g.window.Height
}

func main() {

	game := &Game{}
	game.window.Width = 1000
	game.window.Height = 1000

	snake := Snake{}
	snake.lenght = 20
	snake.head = true
	snake.tail = true
	snake.x = 60
	snake.y = 40
	snake.direction = "right"


	// snake2 := Snake{}
	// snake2.lenght = 20
	// snake2.head = false
	// snake2.tail = false
	// snake2.x = 40
	// snake2.y = 40
	// snake2.direction = "right"
	//
	// snake3 := Snake{}
	// snake3.lenght = 20
	// snake3.head = false
	// snake3.tail = true
	// snake3.x = 20
	// snake3.y = 40
	// snake3.direction = "right"
	//
	// snake4 := Snake{}
	// snake4.lenght = 20
	// snake4.head = false
	// snake4.tail = true
	// snake4.x = 20
	// snake4.y = 40
	// snake4.direction = "right"

	eat := Eat{}
	eat.lenght = 20
	eat.x, eat.y = game.genEat()

	snakeIMG := e.NewImage(snake.lenght, snake.lenght)
	snakeIMG.Fill(color.RGBA{255, 10, 10, 1})
	snake.image = snakeIMG
	// snake2.image = snakeIMG
	// snake3.image = snakeIMG
	// snake4.image = snakeIMG

	eatIMG := e.NewImage(eat.lenght, eat.lenght)
	eatIMG.Fill(color.RGBA{67, 39, 245, 1})
	eat.image = eatIMG

	// snake2.image = snakeIMG

	game.eat = eat
	game.snake = append(game.snake, snake)
	// game.snake = append(game.snake, snake2)
	// game.snake = append(game.snake, snake3)
	// game.snake = append(game.snake, snake4)

	e.SetWindowSize(game.window.Width, game.window.Height)
	e.SetWindowTitle("Your game's title")

	if err := e.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
