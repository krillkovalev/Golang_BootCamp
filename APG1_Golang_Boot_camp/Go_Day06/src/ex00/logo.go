package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"golang.org/x/image/font"
    "golang.org/x/image/font/basicfont"
    "golang.org/x/image/math/fixed"
)


func main() {

	const width = 130
	const height = 50

	img := image.NewRGBA((image.Rect(0, 0, 300, 300)))
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	Triangle(150, 100, 100, 200, 200, 200, img)
	FillTriangleWithGradient(150, 100, 100, 200, 200, 200, img)
	addLabel(img, 90, 40, "My wisdom is here!")


	file, _ := os.Create("amazing_logo.png")

	defer file.Close()
	png.Encode(file, img)

}

// Рисуем линию
func Line(x1, y1, x2, y2 int, img *image.RGBA) {
	dx := absInt(x2 - x1)
	dy := absInt(y2 - y1)
	sx := -1
	if x1 < x2 {
		sx = 1
	}
	sy := -1
	if y1 < y2 {
		sy = 1
	}
	err := dx - dy

	for {
		img.Set(x1, y1, color.Black)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func Triangle(x1, y1, x2, y2, x3, y3 int, img *image.RGBA) {
	Line(x1, y1, x2, y2, img)
	Line(x2, y2, x3, y3, img)
	Line(x3, y3, x1, y1, img)
}

func absInt(x int) int {
	return absDiffInt(x, 0)
}

func absDiffInt(x, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func FillTriangleWithGradient(x1, y1, x2, y2, x3, y3 int, img *image.RGBA) {
    // Определение границ треугольника
    minX := min(x1, min(x2, x3))
    maxX := max(x1, max(x2, x3))
    minY := min(y1, min(y2, y3))
    maxY := max(y1, max(y2, y3))

    // Проход по каждому пикселю внутри предполагаемой области треугольника
    for x := minX; x <= maxX; x++ {
        for y := minY; y <= maxY; y++ {
            // Проверка, находится ли пиксель внутри треугольника
            if isInsideTriangle(x, y, x1, y1, x2, y2, x3, y3) {
                // Генерация градиента
                r := uint8((x - minX) * 255 / (maxX - minX)) // Градиент по X
                g := uint8((y - minY) * 255 / (maxY - minY)) // Градиент по Y
                b := uint8(128)                              // Постоянное значение для синего цвета

                // Установка цвета пикселя
                img.Set(x, y, color.RGBA{r, g, b, 255})
            }
        }
    }
}

// Функция для проверки, находится ли точка (px, py) внутри треугольника
func isInsideTriangle(px, py, x1, y1, x2, y2, x3, y3 int) bool {
    // Вычисляем барицентрические координаты
    area := float64((x2-x1)*(y3-y1) - (x3-x1)*(y2-y1)) // Площадь треугольника
    w1 := float64((px-x1)*(y3-y1) - (x3-x1)*(py-y1)) / area
    w2 := float64((x2-x1)*(py-y1) - (px-x1)*(y2-y1)) / area
    w3 := 1 - w1 - w2

    // Точка внутри треугольника, если все веса в диапазоне [0, 1]
    return w1 >= 0 && w2 >= 0 && w3 >= 0
}

func addLabel(img *image.RGBA, x, y int, label string) {
    col := color.Black
    point := fixed.Point26_6{fixed.I(x), fixed.I(y)}

    d := &font.Drawer{
        Dst:  img,
        Src:  image.NewUniform(col),
        Face: basicfont.Face7x13,
        Dot:  point,
    }
    d.DrawString(label)
}
