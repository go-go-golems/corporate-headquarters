package main

import (
	"fmt"
	"strings"

	svg "github.com/ajstarks/svgo"
	"gopkg.in/yaml.v2"
)

// Element is the interface that all SVG elements implement.
type Element interface {
	Render(canvas *svg.SVG)
}

// SVGDSL represents the root of the YAML DSL.
type SVGDSL struct {
	SVG Canvas `yaml:"svg"`
}

// Canvas represents the SVG canvas configuration.
type Canvas struct {
	Width      int              `yaml:"width"`
	Height     int              `yaml:"height"`
	Background Background       `yaml:"background"`
	Elements   []ElementWrapper `yaml:"elements"`
}

// GetElements converts []ElementWrapper to []Element.
func (c *Canvas) GetElements() []Element {
	elements := make([]Element, len(c.Elements))
	for i, ew := range c.Elements {
		elements[i] = ew.Element
	}
	return elements
}

// Background represents the canvas background.
type Background struct {
	Color string `yaml:"color,omitempty"`
	Image string `yaml:"image,omitempty"`
}

// ElementWrapper is a proxy for unmarshaling different Element types.
type ElementWrapper struct {
	Element
}

// UnmarshalYAML implements custom unmarshaling for the Element interface.
func (ew *ElementWrapper) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Temporary map to extract the "type" field
	var raw map[string]interface{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	typ, ok := raw["type"].(string)
	if !ok {
		return fmt.Errorf("missing or invalid type field")
	}

	// Based on the type, unmarshal into the appropriate struct
	switch typ {
	case "rectangle":
		var rect Rectangle
		if err := mapToStruct(raw, &rect); err != nil {
			return err
		}
		ew.Element = &rect
	case "line":
		var line Line
		if err := mapToStruct(raw, &line); err != nil {
			return err
		}
		ew.Element = &line
	case "image":
		var img Image
		if err := mapToStruct(raw, &img); err != nil {
			return err
		}
		ew.Element = &img
	case "text":
		var txt Text
		if err := mapToStruct(raw, &txt); err != nil {
			return err
		}
		ew.Element = &txt
	case "group":
		var grp Group
		if err := mapToStruct(raw, &grp); err != nil {
			return err
		}
		ew.Element = &grp
	default:
		return fmt.Errorf("unsupported element type: %s", typ)
	}

	return nil
}

// mapToStruct helps in converting a map to a struct using YAML marshalling.
func mapToStruct(m map[string]interface{}, out interface{}) error {
	data, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, out)
}

// Rectangle represents an SVG rectangle.
type Rectangle struct {
	Type        string     `yaml:"type"`
	ID          string     `yaml:"id,omitempty"`
	X           int        `yaml:"x"`
	Y           int        `yaml:"y"`
	Width       int        `yaml:"width"`
	Height      int        `yaml:"height"`
	Fill        string     `yaml:"fill,omitempty"`
	Stroke      string     `yaml:"stroke,omitempty"`
	StrokeWidth int        `yaml:"stroke_width,omitempty"`
	Transform   *Transform `yaml:"transform,omitempty"`
}

// Render renders the rectangle onto the SVG canvas.
func (r *Rectangle) Render(canvas *svg.SVG) {
	styles := buildStyles(r.Fill, r.Stroke, r.StrokeWidth)
	if r.Transform != nil {
		canvas.Gtransform(buildTransform(r.Transform))
	}
	canvas.Rect(r.X, r.Y, r.Width, r.Height, styles)
	if r.Transform != nil {
		canvas.Gend()
	}
}

// Line represents an SVG line.
type Line struct {
	Type        string     `yaml:"type"`
	ID          string     `yaml:"id,omitempty"`
	X1          int        `yaml:"x1"`
	Y1          int        `yaml:"y1"`
	X2          int        `yaml:"x2"`
	Y2          int        `yaml:"y2"`
	Stroke      string     `yaml:"stroke,omitempty"`
	StrokeWidth int        `yaml:"stroke_width,omitempty"`
	Transform   *Transform `yaml:"transform,omitempty"`
}

// Render renders the line onto the SVG canvas.
func (l *Line) Render(canvas *svg.SVG) {
	styles := buildStyles("", l.Stroke, l.StrokeWidth)
	if l.Transform != nil {
		canvas.Gtransform(buildTransform(l.Transform))
	}
	canvas.Line(l.X1, l.Y1, l.X2, l.Y2, styles)
	if l.Transform != nil {
		canvas.Gend()
	}
}

// Image represents an SVG image.
type Image struct {
	Type      string     `yaml:"type"`
	ID        string     `yaml:"id,omitempty"`
	Href      string     `yaml:"href"`
	X         int        `yaml:"x"`
	Y         int        `yaml:"y"`
	Width     int        `yaml:"width"`
	Height    int        `yaml:"height"`
	Transform *Transform `yaml:"transform,omitempty"`
}

// Render renders the image onto the SVG canvas.
func (img *Image) Render(canvas *svg.SVG) {
	styles := buildStyles("", "", 0) // Assuming no additional styles
	if img.Transform != nil {
		canvas.Gtransform(buildTransform(img.Transform))
	}
	canvas.Image(img.X, img.Y, img.Width, img.Height, img.Href, styles)
	if img.Transform != nil {
		canvas.Gend()
	}
}

// Text represents an SVG text element.
type Text struct {
	Type       string     `yaml:"type"`
	ID         string     `yaml:"id,omitempty"`
	X          int        `yaml:"x"`
	Y          int        `yaml:"y"`
	Content    string     `yaml:"content"`
	FontSize   string     `yaml:"font_size,omitempty"`
	FontFamily string     `yaml:"font_family,omitempty"`
	Fill       string     `yaml:"fill,omitempty"`
	TextAnchor string     `yaml:"text_anchor,omitempty"`
	Transform  *Transform `yaml:"transform,omitempty"`
}

// Render renders the text onto the SVG canvas.
func (t *Text) Render(canvas *svg.SVG) {
	styles := buildTextStyles(t.Fill, t.FontSize, t.FontFamily, t.TextAnchor)
	if t.Transform != nil {
		canvas.Gtransform(buildTransform(t.Transform))
	}
	canvas.Text(t.X, t.Y, t.Content, styles)
	if t.Transform != nil {
		canvas.Gend()
	}
}

// Group represents an SVG group, which can contain nested elements.
type Group struct {
	Type      string     `yaml:"type"`
	ID        string     `yaml:"id,omitempty"`
	Transform *Transform `yaml:"transform,omitempty"`
	Elements  []Element  `yaml:"elements"`
}

// Render renders the group and its nested elements onto the SVG canvas.
func (g *Group) Render(canvas *svg.SVG) {
	if g.Transform != nil {
		canvas.Gtransform(buildTransform(g.Transform))
	} else {
		canvas.G()
	}
	for _, elem := range g.Elements {
		elem.Render(canvas)
	}
	canvas.Gend()
}

// Transform represents transformations applied to SVG elements.
type Transform struct {
	Translate []int     `yaml:"translate,omitempty"` // [x, y]
	Rotate    float64   `yaml:"rotate,omitempty"`    // degrees
	Scale     []float64 `yaml:"scale,omitempty"`     // [x, y]
}

// buildStyles constructs the style string for fill, stroke, and stroke-width.
func buildStyles(fill, stroke string, strokeWidth int) string {
	styles := ""
	if fill != "" {
		styles += fmt.Sprintf("fill:%s;", fill)
	}
	if stroke != "" {
		styles += fmt.Sprintf("stroke:%s;", stroke)
	}
	if strokeWidth > 0 {
		styles += fmt.Sprintf("stroke-width:%d;", strokeWidth)
	}
	return styles
}

// buildTextStyles constructs the style string for text elements.
func buildTextStyles(fill, fontSize, fontFamily, textAnchor string) string {
	styles := ""
	if fill != "" {
		styles += fmt.Sprintf("fill:%s;", fill)
	}
	if fontSize != "" {
		styles += fmt.Sprintf("font-size:%s;", fontSize)
	}
	if fontFamily != "" {
		styles += fmt.Sprintf("font-family:%s;", fontFamily)
	}
	if textAnchor != "" {
		styles += fmt.Sprintf("text-anchor:%s;", textAnchor)
	}
	return styles
}

// buildTransform constructs the transformation string based on Translate, Rotate, and Scale.
func buildTransform(t *Transform) string {
	var transforms []string
	if len(t.Translate) == 2 {
		transforms = append(transforms, fmt.Sprintf("translate(%d,%d)", t.Translate[0], t.Translate[1]))
	}
	if t.Rotate != 0 {
		transforms = append(transforms, fmt.Sprintf("rotate(%f)", t.Rotate))
	}
	if len(t.Scale) == 2 {
		transforms = append(transforms, fmt.Sprintf("scale(%f,%f)", t.Scale[0], t.Scale[1]))
	}
	return strings.Join(transforms, " ")
}
