package component

import (
	"fmt"
	"strings"
	"termfedi/layer"
	"termfedi/utils"

	"github.com/gdamore/tcell/v2"
	"golang.org/x/net/html"
)

func DrawNotificationComponent(x int, y int, noti layer.Notification, ctx tcell.Screen, style tcell.Style, maxheight int) {
	switch noti.Type {
	case layer.NOTI_MENTION:
		mention_from := fmt.Sprintf("Got a mention from (@%s)", noti.ReactedUser.User_finger)
		utils.WriteTo(ctx, x+1, y, mention_from, style)
		drawOthers(x, y, noti, ctx, style, maxheight)

	case layer.NOTI_FAVOURITE:
		mention_from := fmt.Sprintf("Got a favourite from (@%s)", noti.ReactedUser.User_finger)
		utils.WriteTo(ctx, x+1, y, mention_from, style)
		drawOthers(x, y, noti, ctx, style, maxheight)

	case layer.NOTI_RENOTE:
		mention_from := fmt.Sprintf("Got a renote from (@%s)", noti.ReactedUser.User_finger)
		utils.WriteTo(ctx, x+1, y, mention_from, style)
		drawOthers(x, y, noti, ctx, style, maxheight)

	case layer.NOTI_FOLLOW:
		mention_from := fmt.Sprintf("Got a follow from (@%s)", noti.ReactedUser.User_finger)
		utils.WriteTo(ctx, x+1, y, mention_from, style)
	}

}

func drawOthers(x int, y int, noti layer.Notification, ctx tcell.Screen, style tcell.Style, maxheight int) {
	utils.WriteTo(ctx, x+1, y+1, "Note RN: ", style)

	var render_targets []string = make([]string, 0)

	var result strings.Builder
	htmls := html.NewTokenizer(strings.NewReader(noti.Content))
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
		utils.WriteTo(ctx, x+1, y+2+i, ntx, style)
		if i >= maxheight-3 {
			break
		}
	}
}
