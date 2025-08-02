package component

import (
	"fmt"
	"strings"
	"termfedi/layer"
	utils "termfedi/utils"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/net/html"
)

/*
	TODO:
	- Support for CW
	- Image/Renote alias
	- Renote/Favoriute count
	- Visiblity Range
*/

func DrawNoteComponent(x int, y int, note layer.Note, ctx tcell.Screen, style tcell.Style, maxheight int, viewcw bool) {
	w, _ := ctx.Size()

	name := fmt.Sprintf("%s (@%s)", note.Author_name, note.Author_finger)
	if len(name) >= w {
		name = fmt.Sprintf("%s (...)", note.Author_name)
	}

	utils.WriteTo(ctx, x+1, y, name, style)

	var content string

	if len(note.Spoiler) > 0 {
		if viewcw {
			content = note.Content
		} else {
			content = fmt.Sprintf("CW: %s", note.Spoiler)
		}
	} else {
		content = note.Content
	}

	if note.IsRenote {
		content = note.Renote
	}

	content = strings.ReplaceAll(content, "\n", "<br>")

	var render_targets []string = make([]string, 0)

	if note.HasMedia {
		for _, media := range note.Medias {
			render_targets = append(render_targets, fmt.Sprintf("[Image: %s]", media))
		}
	}

	var result strings.Builder
	htmls := html.NewTokenizer(strings.NewReader(content))
	loops := true

	for loops {
		tokenType := htmls.Next()
		switch tokenType {
		case html.TextToken:
			h_text := htmls.Token().Data
			width, _ := ctx.Size()
			if len(h_text) >= width {
				result.WriteString(h_text[:width])
				render_targets = append(render_targets, result.String())
				result.Reset()
				result.WriteString(h_text[width:])
			} else {
				result.WriteString(h_text)
			}

		case html.SelfClosingTagToken, html.StartTagToken:
			tname, _ := htmls.TagName()
			if string(tname) == "br" {
				result.WriteString(htmls.Token().Data)
				render_targets = append(render_targets, result.String())
				result.Reset()
			} else {
				result.WriteString(htmls.Token().Data)
			}

		case html.ErrorToken:
			render_targets = append(render_targets, result.String())
			loops = false
		}
	}

	for i, ntx := range render_targets {
		utils.WriteTo(ctx, x+1, y+1+i, ntx, style)
		if i >= maxheight-4 {
			break
		}
	}

	status := fmt.Sprintf("(%s Note) RENOTE: %d | FAVOURITES: %d", layer.VisiblityToText(note.Visiblity), note.RenoteCount, note.ReactionCount)
	utils.WriteTo(ctx, x+1, y+maxheight-2, status, style)
}
