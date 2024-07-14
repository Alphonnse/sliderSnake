package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	e "github.com/hajimehoshi/ebiten/v2"
	eu "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Snake struct {
	image *e.Image

	x int
	y int

	length int

	head bool
	tail bool

	direction string
}

type Eat struct {
	image *e.Image

	x int
	y int

	length int
}

type Gameover struct {
	image *e.Image

	isover bool
}

type Window struct {
	Height int
	Width  int
}

type Game struct {
	window Window
	snake  []Snake
	eat    Eat

	lastUpdate     time.Time
	updateInterval time.Duration

	gameover Gameover
}

func (g *Game) Update() error {

	if g.gameover.isover {
		return nil
	}

	now := time.Now()
	if now.Sub(g.lastUpdate) < g.updateInterval {
		return nil
	}
	g.lastUpdate = now

	// Borders
	if g.snake[0].y <= 0 || g.snake[0].x <= 0 ||
		g.snake[0].y >= g.window.Height-g.snake[0].length ||
		g.snake[0].x >= g.window.Width-g.snake[0].length {
		g.gameover.isover = true
	}

	if e.IsKeyPressed(e.KeyRight) && g.snake[0].direction != "left" && g.snake[0].direction != "right" {
		g.snake[0].direction = "right"
	}
	if e.IsKeyPressed(e.KeyUp) && g.snake[0].direction != "down" && g.snake[0].direction != "up" {
		g.snake[0].direction = "up"
	}
	if e.IsKeyPressed(e.KeyDown) && g.snake[0].direction != "up" && g.snake[0].direction != "down" {
		g.snake[0].direction = "down"
	}
	if e.IsKeyPressed(e.KeyLeft) && g.snake[0].direction != "right" && g.snake[0].direction != "left" {
		g.snake[0].direction = "left"
	}

	// Update positions of snake segments
	for i := len(g.snake) - 1; i > 0; i-- {
		g.snake[i].x = g.snake[i-1].x
		g.snake[i].y = g.snake[i-1].y
	}

	switch g.snake[0].direction {
	case "right":
		g.snake[0].x += 10
	case "up":
		g.snake[0].y -= 10
	case "down":
		g.snake[0].y += 10
	case "left":
		g.snake[0].x -= 10
	}

	if checkIntersection(g.eat.x, g.eat.y, g.eat.length, g.eat.length, g.snake[0].x, g.snake[0].y, g.snake[0].length, g.snake[0].length) {
		g.eat.x, g.eat.y = g.genEat()
		g.grow()
	}

	return nil
}

func (g *Game) genEat() (int, int) {
	needed := false
	var eatX, eatY int
	for !needed {
		eatX = rand.Intn(g.window.Width)
		eatY = rand.Intn(g.window.Height)
		needed = true
		for i := range g.snake {
			if eatX == g.snake[i].x && eatY == g.snake[i].y {
				needed = false
				break
			}
		}
	}
	return eatX, eatY
}

func (g *Game) grow() {
	snake := Snake{}
	snake.length = 20
	snake.head = false
	snake.tail = true

	for i := range g.snake {
		if g.snake[i].tail {
			switch g.snake[i].direction {
			case "right":
				snake.x = g.snake[i].x - snake.length
				snake.y = g.snake[i].y
			case "left":
				snake.x = g.snake[i].x + snake.length
				snake.y = g.snake[i].y
			case "up":
				snake.x = g.snake[i].x
				snake.y = g.snake[i].y + snake.length
			case "down":
				snake.x = g.snake[i].x
				snake.y = g.snake[i].y - snake.length
			}
			g.snake[i].tail = false
			snake.direction = g.snake[i].direction
			break
		}
	}

	snakeIMG := e.NewImage(snake.length, snake.length)
	snakeIMG.Fill(color.RGBA{255, 10, 10, 255})
	snake.image = snakeIMG

	g.snake = append(g.snake, snake)
}

func checkIntersection(x1, y1, w1, h1, x2, y2, w2, h2 int) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}

func (g *Game) Draw(screen *e.Image) {
	if g.gameover.isover {
		op := &e.DrawImageOptions{}
		op.GeoM.Scale(0.7, 0.75)
		screen.DrawImage(g.gameover.image, op)
		return
	}

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

func (game *Game) init() {
	game.window.Width = 1400
	game.window.Height = 900
	game.updateInterval = 15 * time.Millisecond

	snake := Snake{}
	snake.length = 20
	snake.head = true
	snake.tail = true
	snake.x = 60
	snake.y = 40
	snake.direction = "right"

	eat := Eat{}
	eat.length = 20
	eat.x, eat.y = game.genEat()

	snakeIMG := e.NewImage(snake.length, snake.length)
	snakeIMG.Fill(color.RGBA{255, 10, 10, 255})
	snake.image = snakeIMG

	eatIMG := e.NewImage(eat.length, eat.length)
	eatIMG.Fill(color.RGBA{67, 39, 245, 255})
	eat.image = eatIMG

	var err error
	game.gameover.image, _, err = eu.NewImageFromFile("./gameover.png")
	if err != nil {
		log.Fatalf("cant load gameover iamge: %s", err.Error())
	}

	game.eat = eat
	game.snake = append(game.snake, snake)
}

func main() {
	game := &Game{}

	game.init()

	e.SetWindowSize(game.window.Width, game.window.Height)
	e.SetWindowTitle("Your game's title")

	if err := e.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
