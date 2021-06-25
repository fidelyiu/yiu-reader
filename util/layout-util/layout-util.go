package LayoutUtil

import (
	YiuInt "github.com/fidelyiu/yiu-go/int"
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

func (p *point) isEqual(target point) bool {
	return p.x == target.x && p.y == target.y
}

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

func (b *box) isEqual(target box) bool {
	return b.pointA.isEqual(target.pointA) &&
		b.pointB.isEqual(target.pointB) &&
		b.pointC.isEqual(target.pointC) &&
		b.pointD.isEqual(target.pointD)
}

// 两个盒子是否相交
func (b *box) hasIntersect(target box) bool {
	return YiuInt.IsIntersect(b.getMinX(), b.getMaxX(), target.getMinX(), target.getMaxX(), false) &&
		YiuInt.IsIntersect(b.getMinY(), b.getMaxY(), target.getMinY(), target.getMaxY(), false)
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
		Width:  b.getMaxX() - b.getMinX(),
		Height: b.getMaxY() - b.getMinY(),
		Left:   b.getMinX(),
		Top:    b.getMinY(),
	}
}

// getBorderBoxListByWithAndHeight 获取盒子周边的12个盒子
func (b *box) getBorderBoxListByWithAndHeight(width int, height int) [12]box {
	// 输出的盒子顺序
	// 1 | 2  3  | 4
	// --A-------B---
	// 7 |       | 5
	// 8 |       | 6
	// --C-------D---
	// 9 | 10 11 | 12

	return [12]box{
		b.pointA.outBoxLeftTop(width, height),     // 1
		b.pointA.outBoxRightTop(width, height),    // 2
		b.pointB.outBoxLeftTop(width, height),     // 3
		b.pointB.outBoxRightTop(width, height),    // 4
		b.pointB.outBoxRightBottom(width, height), // 5
		b.pointD.outBoxRightTop(width, height),    // 6
		b.pointA.outBoxLeftBottom(width, height),  // 7
		b.pointC.outBoxLeftTop(width, height),     // 8
		b.pointC.outBoxLeftBottom(width, height),  // 9
		b.pointC.outBoxRightBottom(width, height), // 10
		b.pointD.outBoxLeftBottom(width, height),  // 11
		b.pointD.outBoxRightBottom(width, height), // 12
	}
}

func layoutToBox(layout entity.Layout) box {
	return box{
		pointA: point{layout.Left, layout.Top},
		pointB: point{layout.Left + layout.Width, layout.Top},
		pointC: point{layout.Left, layout.Top + layout.Height},
		pointD: point{layout.Left + layout.Width, layout.Top + layout.Height},
	}
}

// 面板
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

// boxCanPush 判断盒子是否可以插入面板中
// 如果盒子和任意一个已插入盒子相交了，那么就不能插入
func (l *layoutStaff) boxCanPush(b box) bool {
	for _, v := range l.boxList {
		if v.hasIntersect(b) {
			return false
		}
	}
	return true
}

// push 将一个盒子插入面板中，返回是否插入成功
func (l *layoutStaff) push(b box) bool {
	if l.boxCanPush(b) {
		// 如果能插入，就直接插入
		l.boxList = append(l.boxList, b)
		l.changeMaxValue(b)
		return true
	}
	return false
}

// boxIsValid 判断该盒子是否在该面板的有效范围内
// 大于最小值x、y
func (l *layoutStaff) boxIsValid(b box) bool {
	if l.maxX == 0 {
		l.maxX = 1080
	}
	if b.getMinX() < l.minX ||
		b.getMinY() < l.minY ||
		b.getMaxX() > l.maxX {
		return false
	}
	for _, v := range l.boxList {
		if v.hasIntersect(b) {
			return false
		}
	}
	return true
}

func (l *layoutStaff) changeMaxValue(b box) {
	if l.maxX < b.getMaxX() {
		l.maxX = b.getMaxX()
	}
	if l.maxY < b.getMaxY() {
		l.maxY = b.getMaxY()
	}
}

// pushWithChange 将一个盒子插入，
func (l *layoutStaff) pushWithChange(b box) box {
	// 先判断上面四个盒子
	for _, v := range l.boxList {
		testBoxList := v.getBorderBoxListByWithAndHeight(b.getWidth(), b.getHeight())
		for i, t := range testBoxList {
			if 1 <= i+1 && i+1 <= 4 {
				if l.boxIsValid(t) {
					l.changeMaxValue(t)
					return t
				}
			}
		}
	}
	// 再判断之间四个盒子
	for _, v := range l.boxList {
		testBoxList := v.getBorderBoxListByWithAndHeight(b.getWidth(), b.getHeight())
		for i, t := range testBoxList {
			if 5 <= i+1 && i+1 <= 8 {
				if l.boxIsValid(t) {
					l.changeMaxValue(t)
					return t
				}
			}
		}
	}
	// 最后判断底下四个盒子
	for _, v := range l.boxList {
		testBoxList := v.getBorderBoxListByWithAndHeight(b.getWidth(), b.getHeight())
		for i, t := range testBoxList {
			if 9 <= i+1 && i+1 <= 12 {
				if l.boxIsValid(t) {
					l.changeMaxValue(t)
					return t
				}
			}
		}
	}
	// 都不符合，找l.MaxY，排第一个
	result := box{
		pointA: point{0, l.maxY},
		pointB: point{b.getWidth(), l.maxY},
		pointC: point{0, l.maxY + b.getHeight()},
		pointD: point{b.getWidth(), l.maxY + b.getHeight()},
	}
	l.maxY = l.maxY + b.getHeight()
	l.changeMaxValue(result)
	return result
}

// newLayoutStaff 生成一个最小值xy为-8的面板
func newLayoutStaff(maxX int) layoutStaff {
	return layoutStaff{
		minX: 0,
		minY: 0,
		maxX: maxX,
	}
}

// FormatLayout 格式化所有传入的布局
func FormatLayout(layoutList *[]entity.Layout, maxX int) {
	staff := newLayoutStaff(maxX)
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

func GetDefaultLayout(layoutList []entity.Layout, layout *entity.Layout, maxX int) {
	// 生成一个默认面板
	staff := newLayoutStaff(maxX)
	// 将所有已有的盒子插入
	for _, v := range layoutList {
		staff.push(layoutToBox(v))
	}
	// 将待插入盒子尝试插入面板，插入后生成临时调整盒子
	changeBox := staff.pushWithChange(layoutToBox(*layout))
	// 将盒子转换成布局
	changLayout := changeBox.toLayout()
	// 修改待加入布局
	layout.Left = changLayout.Left
	layout.Top = changLayout.Top
}
