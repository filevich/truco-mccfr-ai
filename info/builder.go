package info

import (
	"hash"
	"strings"

	"github.com/filevich/truco-ai/abs"
	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/pdt"
)

type InfosetBuilder func(

	p *pdt.Partida,
	m *pdt.Manojo,
	msgs []enco.IMessage,

) Infoset

type BuilderData struct {
	Hash string
	Info string
}

type Builder struct {
	// data fields
	Data *BuilderData
	//
	Hash hash.Hash
	Info InfosetBuilder
	Abs  abs.IAbstraction
}

func BuilderFactory(hash, info, a string) *Builder {
	b := &Builder{
		Data: &BuilderData{
			Hash: hash,
			Info: info,
			// Abs:  a, // <- not needed because of `Builder.Abs.String()`
		},
	}

	b.Hash = ParseHashFn(hash)
	b.Info = nil
	b.Abs = abs.ParseAbs(a)

	if strings.EqualFold(info, "InfosetRondaBase") {
		b.Info = infosetRondaBaseFactory(b.Abs)
	} else if strings.EqualFold(info, "InfosetRondaLarge") {
		b.Info = infosetRondaLargeFactory(b.Abs)
	} else if strings.EqualFold(info, "InfosetRondaXLarge") {
		b.Info = infosetRondaXLargeFactory(b.Abs)
	} else if strings.EqualFold(info, "InfosetRondaXXLarge") {
		b.Info = infosetRondaXXLargeFactory(b.Abs)
	} else if strings.EqualFold(info, "InfosetPartidaXXLarge") {
		b.Info = infosetPartidaXXLargeFactory(b.Abs)
	} else {
		panic("either infoset impl. does not exists or there's an error with args")
	}

	return b
}
