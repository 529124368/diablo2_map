package main

import (
	"fmt"
	_ "image/png"
	"os"
	"runtime"
	"test/mapCreator/dat"
	"test/mapCreator/ds1"
	"test/mapCreator/dt1"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2enum"
	"github.com/hajimehoshi/ebiten/v2"
)

var img [][]*ebiten.Image

var img2 [][]imgWall

var img3 [][]imgWall

//var flg bool = false
var (
	offsetX float64 = 0
	offsetY float64 = 0
	Scale   float64 = 1
	disPlay bool    = true
)

type imgWall struct {
	img *ebiten.Image
	h   int
}

func init() {
	//加载地块dt1素材
	re, _ := os.ReadFile("mapsucai/floor.dt1")
	ss, _ := dt1.LoadDT1(re)
	re, _ = os.ReadFile("mapsucai/objects.dt1")
	ss1, _ := dt1.LoadDT1(re)

	re, _ = os.ReadFile("mapsucai/outdoor/objects.dt1")
	ss2, _ := dt1.LoadDT1(re)
	re, _ = os.ReadFile("mapsucai/outdoor/treegroups.dt1")
	ss3, _ := dt1.LoadDT1(re)
	re, _ = os.ReadFile("mapsucai/fence.dt1")
	ss4, _ := dt1.LoadDT1(re)
	re, _ = os.ReadFile("mapsucai/outdoor/bridge.dt1")
	ss5, _ := dt1.LoadDT1(re)
	re, _ = os.ReadFile("mapsucai/outdoor/stonewall.dt1")
	ss6, _ := dt1.LoadDT1(re)
	re, _ = os.ReadFile("mapsucai/outdoor/river.dt1")
	ss7, err := dt1.LoadDT1(re)

	//wall
	ss2.Tiles = append(ss2.Tiles, ss1.Tiles...)
	ss2.Tiles = append(ss2.Tiles, ss3.Tiles...)
	ss2.Tiles = append(ss2.Tiles, ss4.Tiles...)
	ss2.Tiles = append(ss2.Tiles, ss5.Tiles...)
	ss2.Tiles = append(ss2.Tiles, ss6.Tiles...)
	ss2.Tiles = append(ss2.Tiles, ss7.Tiles...)

	//floor
	ss.Tiles = append(ss.Tiles, ss2.Tiles...)

	if err != nil {
		fmt.Println(err)
	}
	//读取DS1文件
	dd, _ := os.ReadFile("mapsucai/townE1.ds1")
	d, _ := ds1.Unmarshal(dd)

	//加载素材信息提取
	// for i := 0; i < len(d.Files); i++ {
	// 	fmt.Println(strings.ReplaceAll(d.Files[i], "tg1", "dt1"))
	// }

	//floor
	w, h := d.Floors[0].Size()
	img = make([][]*ebiten.Image, h)
	for i := 0; i < h; i++ {
		img[i] = make([]*ebiten.Image, w)
		for j := 0; j < w; j++ {
			ds1Tile := d.Floors[0].Tile(j, i)
			if !ds1Tile.Hidden() && ds1Tile.Prop1 != 0 {
				ds := GetTiles(int(ds1Tile.Style), int(ds1Tile.Sequence), 0, ss.Tiles)
				if ds != nil {
					img[i][j] = getTitleImage(ds[ds1Tile.RandomIndex])
				}
			}
		}
	}

	//wall
	img2 = make([][]imgWall, h)
	for i := 0; i < h; i++ {
		img2[i] = make([]imgWall, w)
		for j := 0; j < w; j++ {
			ds1Tile := d.Walls[0].Tile(j, i)
			if !ds1Tile.Hidden() && ds1Tile.Prop1 != 0 {
				ds := GetTiles(int(ds1Tile.Style), int(ds1Tile.Sequence), ds1Tile.Type, ss2.Tiles)
				if ds != nil {
					if ds1Tile.Type == d2enum.TileRightPartOfNorthCornerWall {
						dss := GetTiles(int(ds1Tile.Style), int(ds1Tile.Sequence), d2enum.TileLeftPartOfNorthCornerWall, ss2.Tiles)
						if dss != nil && dss[ds1Tile.RandomIndex].Height < ds[ds1Tile.RandomIndex].Height {
							m, h := getWallTitleImage(dss[ds1Tile.RandomIndex], ds1Tile)
							img2[i][j].img = m
							img2[i][j].h = h
						} else {
							m, h := getWallTitleImage(ds[ds1Tile.RandomIndex], ds1Tile)
							img2[i][j].img = m
							img2[i][j].h = h
						}
					} else {
						m, h := getWallTitleImage(ds[ds1Tile.RandomIndex], ds1Tile)
						img2[i][j].img = m
						img2[i][j].h = h
					}
				}
			}
		}
	}

	//地图显示补图  hardcode
	img3 = make([][]imgWall, h)
	for i := 0; i < h; i++ {
		img3[i] = make([]imgWall, w)
		for j := 0; j < w; j++ {
			if j == 18 && i == 6 {
				img3[i][j].img = getTitleImage(ss4.Tiles[7])
				img3[i][j].h = -170
			}
			if j == 30 && i == 11 {
				img3[i][j].img = getTitleImage(ss1.Tiles[22])
				img3[i][j].h = -80
			}
			if j == 31 && i == 11 {
				img3[i][j].img = getTitleImage(ss1.Tiles[23])
				img3[i][j].h = -110
			}

			if j == 32 && i == 11 {
				img3[i][j].img = getTitleImage(ss1.Tiles[24])
				img3[i][j].h = -110
			}
			if j == 22 && i == 22 {
				img3[i][j].img = getTitleImage(ss1.Tiles[24])
				img3[i][j].h = -110
			}
			if j == 33 && i == 11 {
				img3[i][j].img = getTitleImage(ss1.Tiles[21])
				img3[i][j].h = -20
			}
			if j == 33 && i == 10 {
				img3[i][j].img = getTitleImage(ss1.Tiles[25])
				img3[i][j].h = -80
			}
			if j == 33 && i == 9 {
				img3[i][j].img = getTitleImage(ss1.Tiles[26])
				img3[i][j].h = -80
			}
			if j == 33 && i == 8 {
				img3[i][j].img = getTitleImage(ss1.Tiles[27])
				img3[i][j].h = -40
			}
			if j == 33 && i == 7 {
				img3[i][j].img = getTitleImage(ss1.Tiles[28])
				img3[i][j].h = -40
			}
			if j == 22 && i == 22 {
				img3[i][j].img = getTitleImage(ss1.Tiles[19])
				img3[i][j].h = -80
			}
			if j == 21 && i == 24 {
				img3[i][j].img = getTitleImage(ss1.Tiles[16])
				img3[i][j].h = -80
			}

		}
	}
	go func() {
		runtime.GC()
	}()
}

func getTitleImage(tileData dt1.Tile) *ebiten.Image {
	tileYMinimum := int32(0)
	for _, block := range tileData.Blocks {
		tileYMinimum = dt1.MinInt32(tileYMinimum, int32(block.Y))
	}
	tileYOffset := dt1.AbsInt32(tileYMinimum)
	tileHeight := dt1.AbsInt32(tileData.Height)
	indexData := make([]byte, tileData.Width*int32(tileHeight))
	dt1.DecodeTileGfxData(tileData.Blocks, &indexData, tileYOffset, tileData.Width)
	//加载调色板
	re, _ := os.ReadFile("mapsucai/pal.dat")
	w, _ := dat.Load(re)
	pixels := dt1.ImgIndexToRGBA(indexData, w)
	imgss := ebiten.NewImage(int(tileData.Width), int(tileHeight))
	imgss.ReplacePixels(pixels)
	return imgss
}

func getWallTitleImage(tileData dt1.Tile, tile *ds1.Tile) (*ebiten.Image, int) {

	tileMinY := int32(0)
	tileMaxY := int32(0)
	for _, block := range tileData.Blocks {

		tileMinY = dt1.MinInt32(tileMinY, int32(block.Y))
		tileMaxY = dt1.MaxInt32(tileMaxY, int32(block.Y+32))
	}

	realHeight := dt1.MaxInt32(dt1.AbsInt32(tileData.Height), tileMaxY-tileMinY)
	tileYOffset := -tileMinY

	if tile.Type == d2enum.TileRoof {
		tile.YAdjust = -int(tileData.RoofHeight)
	} else {
		tile.YAdjust = int(tileMinY) + 80
	}

	indexData := make([]byte, 160*realHeight)
	dt1.DecodeTileGfxData(tileData.Blocks, &indexData, tileYOffset, 160)
	//加载调色板
	re, _ := os.ReadFile("mapsucai/pal.dat")
	w, _ := dat.Load(re)
	pixels := dt1.ImgIndexToRGBA(indexData, w)
	imgss := ebiten.NewImage(160, int(realHeight))
	imgss.ReplacePixels(pixels)
	return imgss, tile.YAdjust
}

//根据ds1 获取对应dt1
func GetTiles(style, sequence int, tileType d2enum.TileType, m []dt1.Tile) []dt1.Tile {
	tiles := make([]dt1.Tile, 0)

	for idx := range m {
		if m[idx].Style != int32(style) || m[idx].Sequence != int32(sequence) ||
			m[idx].Type != int32(tileType) {
			continue
		}
		tiles = append(tiles, m[idx])
	}
	if len(tiles) == 0 {
		return nil
	}
	return tiles
}

type ATest struct {
}

func (a *ATest) Update() (err error) {
	//this is where the code would be
	// if !flg {
	// 	flg = true
	// 	go func() {
	// 		//newImage := Screenshot(img)
	// 		outFile, err := os.Create("changed.png")
	// 		if err != nil {

	// 			log.Fatal(err)

	// 		}

	// 		defer outFile.Close()

	// 		errs := png.Encode(outFile, img)
	// 		if errs != nil {
	// 			fmt.Println(err)
	// 		}
	// 	}()
	// }
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		offsetX -= 50
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		offsetX += 50
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		offsetY -= 50
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		offsetY += 50
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		disPlay = !disPlay
	}

	return nil
}
func (a *ATest) Draw(screen *ebiten.Image) {
	if true {
		//floor
		sumX := 0
		startY := 0
		if disPlay {
			for i := 0; i < 41; i++ {
				startY += 40
				sumX = 0
				for j := 0; j < 57; j++ {
					s := img[i][j]
					sumX += 80
					if s != nil {
						op := &ebiten.DrawImageOptions{}
						op.GeoM.Translate(float64(i)*(-80)+float64(sumX)+offsetX, float64(startY)+float64(j)*40+offsetY)
						op.GeoM.Scale(Scale, Scale)
						screen.DrawImage(s, op)
					}
				}
			}
		}

		//补图
		sumX = 0
		startY = 0
		for i := 0; i < 41; i++ {
			startY += 40
			sumX = 0
			for j := 0; j < 57; j++ {
				s := img3[i][j].img
				sumX += 80
				if s != nil {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(i)*(-80)+float64(sumX)+offsetX, float64(startY)+float64(j)*40+offsetY+float64(img3[i][j].h))
					op.GeoM.Scale(Scale, Scale)
					screen.DrawImage(s, op)
				}
			}
		}
		//wall
		sumX = 0
		startY = 0
		for i := 0; i < 41; i++ {
			startY += 40
			sumX = 0
			for j := 0; j < 57; j++ {
				s := img2[i][j].img
				sumX += 80
				if s != nil {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(float64(i)*(-80)+float64(sumX)+offsetX, float64(startY)+float64(j)*40+offsetY+float64(img2[i][j].h))
					op.GeoM.Scale(Scale, Scale)
					screen.DrawImage(s, op)
				}
			}
		}

	}

}
func (a *ATest) Layout(ow, oh int) (int, int) {
	return 1920, 1080
}
func main() {
	ebiten.SetWindowSize(1200, 780)
	ebiten.SetWindowTitle("diablo Map title dump")
	ebiten.RunGame(&ATest{})
}
