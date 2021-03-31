package util

import (
	"encoding/json"
	"encoding/xml"
	"regexp"
	"strconv"

	"github.com/google/uuid"
)

type Node struct {
	ID            string  `json:"_id"`
	Index         int     `json:"index" xml:"index,attr"`
	Text          string  `json:"text" xml:"text,attr"`
	ResourceId    string  `json:"resourceId" xml:"resource-id,attr"`
	Type          string  `json:"_type" xml:"class,attr"`
	Package       string  `json:"package" xml:"package,attr"`
	Description   string  `json:"description" xml:"content-desc,attr"`
	Checkable     bool    `json:"checkable" xml:"checkable,attr"`
	Checked       bool    `json:"checked" xml:"checked,attr"`
	Clickable     bool    `json:"clickable" xml:"clickable,attr"`
	Enable        bool    `json:"enabled" xml:"enabled,attr"`
	Focusable     bool    `json:"focusable" xml:"focusable,attr"`
	Focused       bool    `json:"focused" xml:"focused,attr"`
	Scrollable    bool    `json:"scrollable" xml:"scrollable,attr"`
	LongClickable bool    `json:"longClickable" xml:"long-clickable,attr"`
	Password      bool    `json:"password" xml:"password,attr"`
	Selected      bool    `json:"selected" xml:"selected,attr"`
	Bounds        string  `json:"-" xml:"bounds,attr"`
	Nodes         []Child `json:"children" xml:"node"`
}

type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Child struct {
	Node
	Rect `json:"rect"`
}

func (n *Child) MarshalJSON() ([]byte, error) {
	n.Node.ID = uuid.NewString()
	if n.Bounds == "" {
		child := struct {
			Node
			Rect `json:"rect"`
		}{
			Node: n.Node,
			Rect: Rect{},
		}
		return json.Marshal(child)
	}
	re, err := regexp.Compile(`\[(\d+),(\d+)\]\[(\d+),(\d+)\]`)
	if err == nil {
		arr := re.FindStringSubmatch(n.Bounds)
		if len(arr) == 5 {
			x, _ := strconv.Atoi(arr[1])
			y, _ := strconv.Atoi(arr[2])
			w, _ := strconv.Atoi(arr[3])
			h, _ := strconv.Atoi(arr[4])
			n.Rect.X = x
			n.Rect.Y = y
			n.Rect.Width = w - x
			n.Rect.Height = h - y
		}
	}
	child := struct {
		Node
		Rect `json:"rect"`
	}{
		Node: n.Node,
		Rect: n.Rect,
	}
	return json.Marshal(child)
}

type Hierarchy struct {
	ID       string  `json:"_id"`
	Rotation string  `json:"rotation" xml:"rotation,attr"`
	Nodes    []Child `json:"children" xml:"node"`
}
type XML struct {
	Hierarchy `json:"hierarchy" xml:"hierarchy"`
}

func Convert(data []byte) (x XML, err error) {
	x = XML{}
	err = xml.Unmarshal(data, &x)
	return
}
