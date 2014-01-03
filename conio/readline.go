package conio

import "bytes"
import "fmt"
import "io"
import "strings"

import "github.com/mattn/go-runewidth"

func getCharWidth(n rune)int{
	return runewidth.RuneWidth(n)
	// if n > 0xFF {
	//	return 2;
	//}else{
	//	return 1;
	//}
}

func putRep(ch rune,n int){
	for i := 0 ; i < n ; i++ {
		fmt.Printf("%c",ch)
	}
}

type ReadLineBuffer struct{
	buffer[]rune
	length int
	cursor int
	unicode rune
	keycode uint16
	viewstart int
	viewwidth int
}

func (this*ReadLineBuffer)Insert(pos int,c[]rune)bool{
	n := len(c)
	for this.length + n >= len(this.buffer) {
        tmp := make([]rune,len(this.buffer)*2)
        copy(tmp,this.buffer)
        this.buffer = tmp
	}
	for i := this.length ; i >= pos ; i-- {
		this.buffer[i+n] = this.buffer[i]
	}
	for i := 0 ; i < n ; i++ {
		this.buffer[pos+i] = c[i]
	}
	this.length += n
	return true
}

func (this*ReadLineBuffer)InsertString(pos int,s string)int{
	list := []rune{}
	reader := strings.NewReader(s)
	for {
		r,_,err := reader.ReadRune()
		if err == io.EOF {
			break
		}
		list = append(list,r)
	}
	if this.Insert(pos,list) {
		return len(list)
	}else{
		return -1
	}
}

func (this*ReadLineBuffer)Delete(pos int,n int)int{
	if this.length < pos+n {
		return 0
	}
	delw := 0
	for i := pos ; i < pos+n ; i++{
		delw += getCharWidth(this.buffer[i])
	}
	for i := pos ; i < this.length-n ; i++ {
		this.buffer[i] = this.buffer[i+n]
	}
	this.length -= n
	return delw
}

func (this*ReadLineBuffer)getWidthBetween(from int,to int)int{
	width := 0
	for i := from ; i < to ; i++ {
		width += getCharWidth(this.buffer[i])
	}
	return width
}

func (this*ReadLineBuffer)Repaint(pos int,del int){
	bs := 0
	vp := this.getWidthBetween( this.viewstart , pos )

	for i:=pos ; i < this.length ; i++ {
		w1 := getCharWidth(this.buffer[i])
		vp += w1
		if vp >= this.viewwidth {
			break
		}
		putRep(this.buffer[i],1)
		bs += w1
	}
	putRep(' ',del)
	putRep('\b',bs+del)
}

func (this*ReadLineBuffer)String() string {
	var result bytes.Buffer
	for i:=0 ; i< this.length ; i++ {
		result.WriteRune(this.buffer[i])
	}
	return result.String()
}

type KeyFuncResult int

const (
	CONTINUE  KeyFuncResult = iota
	ENTER     KeyFuncResult = iota
	ABORT     KeyFuncResult = iota
)

func KeyFuncPass(this*ReadLineBuffer) KeyFuncResult {
	return CONTINUE
}

func KeyFuncEnter(this*ReadLineBuffer) KeyFuncResult { // Ctrl-M
	return ENTER
}

func KeyFuncHead(this*ReadLineBuffer) KeyFuncResult { // Ctrl-A
	putRep('\b',this.getWidthBetween(this.viewstart,this.cursor))
	this.cursor = 0
	this.viewstart = 0
	this.Repaint(0,1)
	return CONTINUE
}

func KeyFuncBackword(this*ReadLineBuffer) KeyFuncResult { // Ctrl-B
	if this.cursor <= 0 {
		return CONTINUE
	}
	this.cursor--
	if this.cursor < this.viewstart {
		this.viewstart--
		this.Repaint(this.cursor,1)
	}else{
		putRep('\b',getCharWidth(this.buffer[this.cursor]))
	}
	return CONTINUE
}

func KeyFuncTail(this*ReadLineBuffer) KeyFuncResult {// Ctrl-E
	allength := this.getWidthBetween(this.viewstart,this.length)
	if allength < this.viewwidth {
		for ; this.cursor < this.length ; this.cursor++ {
			putRep(this.buffer[this.cursor],1)
		}
	}else{
		putRep('\a',1)
		putRep('\b',this.getWidthBetween(this.viewstart,this.cursor))
		this.viewstart = this.length-1
		w := getCharWidth(this.buffer[this.viewstart])
		for{
			if this.viewstart <= 0 {
				break
			}
			w_ := w + getCharWidth( this.buffer[this.viewstart-1] )
			if w_ >= this.viewwidth {
				break
			}
			w = w_
			this.viewstart--
		}
		for this.cursor = this.viewstart ; this.cursor < this.length ; this.cursor++ {
			putRep( this.buffer[this.cursor] , 1 )
		}
	}
	return CONTINUE
}

func KeyFuncForward(this*ReadLineBuffer) KeyFuncResult {// Ctrl-F
	if this.cursor >= this.length {
		return CONTINUE
	}
	w := this.getWidthBetween(this.viewstart,this.cursor+1)
	if w < this.viewwidth {
		// No Scroll
		putRep(this.buffer[this.cursor],1)
	}else{
		// Right Scroll
		putRep('\b',this.getWidthBetween(this.viewstart,this.cursor))
		if getCharWidth(this.buffer[this.cursor]) > getCharWidth(this.buffer[this.viewstart]) {
			this.viewstart++
		}
		this.viewstart++
		for i:=this.viewstart ; i <= this.cursor ; i++ {
			putRep(this.buffer[i],1)
		}
		putRep(' ',1)
		putRep('\b',1)
	}
	this.cursor++
	return CONTINUE
}

func KeyFuncBackSpace(this*ReadLineBuffer) KeyFuncResult {// Backspace
	if this.cursor > 0 {
		this.cursor--
		delw := this.Delete( this.cursor , 1 )
		if this.cursor >= this.viewstart {
			putRep('\b',delw)
		}else{
			this.viewstart = this.cursor
		}
		this.Repaint(this.cursor,delw)
	}
	return CONTINUE
}

func KeyFuncDelete(this*ReadLineBuffer) KeyFuncResult { // Ctrl-D
	delw := this.Delete( this.cursor , 1 )
	this.Repaint(this.cursor,delw)
	return CONTINUE
}

func KeyFuncInsertSelf(this*ReadLineBuffer) KeyFuncResult {
	ch := this.unicode
	if ch < 0x20 || ! this.Insert( this.cursor , []rune{ch} ) {
		return CONTINUE
	}
	w := this.getWidthBetween( this.viewstart , this.cursor )
	w1 := getCharWidth(ch)
	if w+w1 >= this.viewwidth {
		// scroll left
		putRep('\b',w)
		if getCharWidth(this.buffer[this.viewstart]) < w1 {
			this.viewstart++
		}
		this.viewstart++
		for i:= this.viewstart ; i <= this.cursor ; i++ {
			putRep(this.buffer[i],1)
		}
		putRep(' ',1)
		putRep('\b',1)
	}else{
		this.Repaint( this.cursor , -w1 )
	}
	this.cursor++
	return CONTINUE
}

func KeyFuncInsertReport(this*ReadLineBuffer) KeyFuncResult {
	L := this.InsertString( this.cursor , fmt.Sprintf("[%02X]",this.keycode) )
	if L >= 0 {
		this.Repaint( this.cursor , -L )
		this.cursor += L
	}
	return CONTINUE
}

var KeyMap = map[rune]func(*ReadLineBuffer)KeyFuncResult {
	'\r'   : KeyFuncEnter ,
	'\x01' : KeyFuncHead ,
	'\x02' : KeyFuncBackword ,
	'\x05' : KeyFuncTail ,
	'\x06' : KeyFuncForward ,
	'\b'   : KeyFuncBackSpace ,
	'\x04' : KeyFuncDelete ,
	'\x7F' : KeyFuncDelete ,
}

const (
	K_LEFT  = 0x25
	K_RIGHT = 0x27
	K_DEL   = 0x2E
	K_HOME  = 0x24
	K_END   = 0x23
	K_CTRL  = 0x11
	K_SHIFT = 0x10
)

var ZeroMap = map[uint16]func(*ReadLineBuffer)KeyFuncResult {
	K_LEFT  : KeyFuncBackword ,
	K_RIGHT : KeyFuncForward ,
	K_DEL   : KeyFuncDelete ,
	K_HOME  : KeyFuncHead ,
	K_END   : KeyFuncTail ,
	K_CTRL  : KeyFuncPass ,
	K_SHIFT : KeyFuncPass ,
}

func ReadLine() string {
	var this ReadLineBuffer
    this.buffer = make([]rune,20)
	this.length=0
	this.cursor=0
	this.viewstart = 0
	this.viewwidth = 60
	for{
		this.unicode , this.keycode = GetKey()
		var f func(*ReadLineBuffer)KeyFuncResult
		if this.unicode != 0 {
			f = KeyMap[this.unicode]
			if f == nil {
				f = KeyFuncInsertSelf
			}
		}else{
			f = ZeroMap[this.keycode]
			if f == nil {
				f = KeyFuncPass
			}
		}
		rc := f(&this)
		if rc == ENTER {
			fmt.Print("\n")
			return this.String()
		}
	}
}
// vim:set ts=4 sw=4 fenc=utf8 :
