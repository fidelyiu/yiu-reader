package LayoutUtil

import (
	"yiu/yiu-reader/model/entity"
)

type point struct {
	x int
	y int
}

// 点的四个方向输出的盒子
// (x-w, y-h)|(x, y-h)  |  (x, y-h)|(x+w, y-h)
// (x-w, y  )|(x, y  )  |  (x, y  )|(x+w, y  )
// -------------------(x,y)-------------------
// (x-w, y  )|(x, y  )  |  (x, y  )|(x+w, y  )
// (x-w, y+h)|(x, y+h)  |  (x, y+h)|(x+w, y+h)

func (p *point) outBoxLeftTop(width int, height int) box {
	return box{
		pointA: point{p.x - width, p.y - height},
		pointB: point{p.x, p.y - height},
		pointC: point{p.x - width, p.y},
		pointD: point{p.x, p.y},
	}
}
func (p *point) outBoxRightTop(width int, height int) box {
	return box{
		pointA: point{p.x, p.y - height},
		pointB: point{p.x + width, p.y - height},
		pointC: point{p.x, p.y},
		pointD: point{p.x + width, p.y},
	}
}
func (p *point) outBoxLeftBottom(width int, height int) box {
	return box{
		pointA: point{p.x - width, p.y},
		pointB: point{p.x, p.y},
		pointC: point{p.x - width, p.y + height},
		pointD: point{p.x, p.y + height},
	}
}
func (p *point) outBoxRightBottom(width int, height int) box {
	return box{
		pointA: point{p.x, p.y},
		pointB: point{p.x + width, p.y},
		pointC: point{p.x, p.y + height},
		pointD: point{p.x + width, p.y + height},
	}
}

type box struct {
	// A B
	// C D
	pointA point
	pointB point
	pointC point
	pointD point
}

// 两个盒子是否相交
func (b *box) hasIntersect(target box) bool {
	return b.pointInBox(target.pointA) ||
		b.pointInBox(target.pointB) ||
		b.pointInBox(target.pointC) ||
		b.pointInBox(target.pointD)
}

// pointInBox 点是否在盒子中
func (b *box) pointInBox(p point) bool {
	return b.getMinX() < p.x && p.x < b.getMaxX() &&
		b.getMinY() < p.y && p.y < b.getMaxY()
}

// getMinX 获取值x轴最小值
func (b *box) getMinX() int {
	return b.pointA.x
}

// getMinY 获取值y轴最小值
func (b *box) getMinY() int {
	return b.pointA.y
}

// getMaxX 获取值x轴最大值
func (b *box) getMaxX() int {
	return b.pointD.x
}

// getMaxY 获取值y轴最大值
func (b *box) getMaxY() int {
	return b.pointD.y
}

func (b *box) getWidth() int {
	return b.getMaxX() - b.getMinX()
}
func (b *box) getHeight() int {
	return b.getMaxY() - b.getMinY()
}

func (b *box) toLayout() entity.Layout {
	return entity.Layout{
		Width:  b.getMaxX() - b.getMinX() - 16,
		Height: b.getMaxY() - b.getMinY() - 16,
		Left:   b.getMinX() + 8,
		Top:    b.getMinY() + 8,
	}
}

func (b *box) getBorderBoxListByWithAndHeight(width int, height int) [12]box {
	// 输出的盒子顺序
	// 1 | 4 5 | 7
	// --A-----B---
	// 2 |     | 10
	// 3 |     | 11
	// --C-----D---
	// 6 | 8 9 | 12

	return [12]box{
		b.pointA.outBoxLeftTop(width, height),     // 1
		b.pointA.outBoxLeftBottom(width, height),  // 2
		b.pointC.outBoxLeftTop(width, height),     // 3
		b.pointA.outBoxRightTop(width, height),    // 4
		b.pointB.outBoxLeftTop(width, height),     // 5
		b.pointC.outBoxLeftBottom(width, height),  // 6
		b.pointB.outBoxRightTop(width, height),    // 7
		b.pointC.outBoxRightBottom(width, height), // 8
		b.pointD.outBoxLeftBottom(width, height),  // 9
		b.pointB.outBoxRightBottom(width, height), // 10
		b.pointD.outBoxRightTop(width, height),    // 11
		b.pointD.outBoxRightBottom(width, height), // 12
	}
}

func layoutToBox(layout entity.Layout) box {
	return box{
		pointA: point{layout.Left - 8, layout.Top - 8},
		pointB: point{layout.Left + layout.Width + 8, layout.Top - 8},
		pointC: point{layout.Left - 8, layout.Top + layout.Height + 8},
		pointD: point{layout.Left + layout.Width + 8, layout.Top + layout.Height + 8},
	}
}

type layoutStaff struct {
	boxList []box
	minX    int
	minY    int
	maxX    int
	maxY    int
}

// pointHasOutOfBounds 点是否越界
func (l *layoutStaff) pointInStaff(x int, y int) bool {
	return l.minX < x && x < l.maxX &&
		l.minY < y && y < l.maxY
}

func (l *layoutStaff) boxCanPush(b box) bool {
	for _, v := range l.boxList {
		if v.hasIntersect(b) {
			return false
		}
	}
	return true
}

func (l *layoutStaff) push(b box) bool {
	if l.boxCanPush(b) {
		l.boxList = append(l.boxList, b)
		l.chengMaxValue(b)
		return true
	}
	return false
}

func (l *layoutStaff) boxIsValid(b box) bool {
	if b.getMinX() < l.minX || b.getMinY() < l.minY &&
		b.getMaxX() <= l.maxX {
		return false
	}
	for _, v := range l.boxList {
		if v.hasIntersect(b) {
			return false
		}
	}
	return true
}

func (l *layoutStaff) chengMaxValue(b box) {
	if l.maxX > b.getMaxX() {
		l.maxX = b.getMaxX()
	}
	if l.maxY > b.getMaxY() {
		l.maxY = b.getMaxY()
	}
}

func (l *layoutStaff) pushWithChange(b box) box {
	for _, v := range l.boxList {
		testBoxList := v.getBorderBoxListByWithAndHeight(b.getWidth(), b.getHeight())
		for _, t := range testBoxList {
			if l.boxIsValid(t) {
				l.chengMaxValue(t)
				return t
			}
		}
	}
	// 都不符合，找l.MaxY，排第一个
	result := box{
		pointA: point{-8, l.maxY},
		pointB: point{-8 + b.getWidth(), l.maxY},
		pointC: point{-8, l.maxY + b.getHeight()},
		pointD: point{-8 + b.getWidth(), l.maxY + b.getHeight()},
	}
	l.maxY = l.maxY + b.getHeight()
	l.chengMaxValue(result)
	return result
}

func newLayoutStaff() layoutStaff {
	return layoutStaff{
		minX: -8,
		minY: -8,
	}
}

func FormatLayout(layoutList *[]entity.Layout) {
	staff := newLayoutStaff()
	var notPushLayout []*entity.Layout
	for _, v := range *layoutList {
		if !staff.push(layoutToBox(v)) {
			notPushLayout = append(notPushLayout, &v)
		}
	}
	for _, v := range notPushLayout {
		changeBox := staff.pushWithChange(layoutToBox(*v))
		changLayout := changeBox.toLayout()
		(*v).Left = changLayout.Left
		(*v).Height = changLayout.Height
	}
}

func GetDefaultLayout(layoutList []entity.Layout, layout *entity.Layout) {
	staff := newLayoutStaff()
	for _, v := range layoutList {
		staff.push(layoutToBox(v))
	}
	changeBox := staff.pushWithChange(layoutToBox(*layout))
	changLayout := changeBox.toLayout()
	(*layout).Left = changLayout.Left
	(*layout).Height = changLayout.Height
}
