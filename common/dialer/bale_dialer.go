package dialer

import (
	"context"

	"github.com/sagernet/sing-box/common/dialer/baleDialing"
	"github.com/sagernet/sing-box/option"
	N "github.com/sagernet/sing/common/network"
)

func NewBaleDialer(ctx context.Context, options option.DialerOptions, remoteIsDomain bool) (N.Dialer, error) {
	return NewBaleDialerWithOptions(Options{
		Context:        ctx,
		Options:        options,
		RemoteIsDomain: remoteIsDomain,
	})
}

func NewBaleDialerWithOptions(options Options) (N.Dialer, error) {
	dialOptions := options.Options

	// myEn := "ABCDEFGHIJKL"
	// eResult := myEn[1:6]
	// print(eResult)

	// eResult2 := baleDialing.SubString(myEn, 1, 6)
	// print(eResult2)

	// myValue := ""
	// temp := "سلام،"
	// for range 10 {
	// 	myValue += temp
	// }

	// myValue += "آخر"

	// result := baleDialing.MakeChunks(myValue, utf8.RuneCountInString(myValue), 5)
	// print(result)

	// dialer, err := NewDefault(options.Context, dialOptions)
	dialer, err := baleDialing.NewBaleDialerContainer(dialOptions)
	if err != nil {
		return nil, err
	}
	return dialer, nil
}
