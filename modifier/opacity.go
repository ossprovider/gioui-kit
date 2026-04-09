package modifier

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
)

// OpacityMod applies an opacity modifier.
type OpacityMod struct {
	Opacity float32 // 0.0 - 1.0
}

func (o OpacityMod) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {
	defer paint.PushOpacity(gtx.Ops, o.Opacity).Pop()
	return w(gtx)
}
