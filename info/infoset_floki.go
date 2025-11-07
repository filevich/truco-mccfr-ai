package info

import (
	"encoding/hex"
	"encoding/json"
	"hash"
	"math/rand"
	"slices"

	"github.com/filevich/truco-mccfr-ai/abs"
	"github.com/filevich/truco-mccfr-ai/utils"
	"github.com/truquito/gotruco/enco"
	"github.com/truquito/gotruco/pdt"
)

type InfosetFloki struct {
	Vision string

	// floki fields
	we_are_starters                 bool
	muestra_valor                   int
	cards_i_still_own               int
	available_teammate_cards        int
	all_our_cards                   int
	num_our_flores                  int
	num_opp_flores                  int
	envite_estado                   int
	envite_cantado_por_our_team     bool
	truco_estado                    int
	truco_cantado_por_our_team      bool
	mano_en_juego                   int
	resultado_mano                  int
	num_opps_after_me_that_can_play int

	max_poder  []int
	is_from_us []bool

	Chi []int
}

func (info *InfosetFloki) setWeAreStarters(p *pdt.Partida, m *pdt.Manojo) {
	starter_team := p.Ronda.Manojos[p.Ronda.ElMano].Jugador.Equipo
	info.we_are_starters = starter_team == m.Jugador.Equipo
}

func (info *InfosetFloki) setMuestraValor(p *pdt.Partida) {
	info.muestra_valor = p.Ronda.Muestra.Valor
}

func (info *InfosetFloki) setCardsIStillOwn(p *pdt.Partida, m *pdt.Manojo, abs abs.IAbstraction) {
	info.cards_i_still_own = 1
	for i, c := range m.Cartas {
		if !m.Tiradas[i] {
			a := abs.Abstract(c, &p.Ronda.Muestra)
			info.cards_i_still_own *= utils.AllPrimes[a]
		}
	}
}

func (info *InfosetFloki) setTeammatesAfterMeCardsStillOwned(p *pdt.Partida, m *pdt.Manojo, abs abs.IAbstraction) {
	info.available_teammate_cards = 1

	num_players := len(p.Ronda.Manojos)
	our_team := m.Jugador.Equipo

	ix := (p.Ronda.MIXS[m.Jugador.ID] + 1) % num_players
	end := (int(p.Ronda.ElMano) + num_players - 1) % num_players
	for count := 0; count <= num_players; count++ {
		is_teammate := p.Ronda.Manojos[ix].Jugador.Equipo == our_team
		if !is_teammate {
			continue
		}
		has_folded := p.Ronda.Manojos[ix].SeFueAlMazo
		for i := 0; i < 3; i++ {
			tirada := p.Ronda.Manojos[ix].Tiradas[i]
			if !has_folded && !tirada {
				a := abs.Abstract(p.Ronda.Manojos[ix].Cartas[i], &p.Ronda.Muestra)
				info.available_teammate_cards *= utils.AllPrimes[a]
			}
		}
		if ix == end {
			break
		}
		ix = (ix + 1) % num_players
	}
}

func (info *InfosetFloki) setAllOurCards(p *pdt.Partida, m *pdt.Manojo, abs abs.IAbstraction) {
	info.all_our_cards = 1
	num_players := len(p.Ronda.Manojos)
	our_team := m.Jugador.Equipo

	for ix := 0; ix < num_players; ix++ {
		is_teammate := p.Ronda.Manojos[ix].Jugador.Equipo == our_team
		if !is_teammate {
			continue
		}
		for i := 0; i < 3; i++ {
			a := abs.Abstract(p.Ronda.Manojos[ix].Cartas[i], &p.Ronda.Muestra)
			info.all_our_cards *= utils.AllPrimes[a]
		}
	}
}

func (info *InfosetFloki) setFlores(p *pdt.Partida, m *pdt.Manojo) {
	info.num_our_flores = 0
	info.num_opp_flores = 0

	our_team := m.Jugador.Equipo

	for _, con_flor := range p.Ronda.Envite.JugadoresConFlor {
		is_teammate := con_flor.Jugador.Equipo == our_team
		if is_teammate {
			// If it's one of us, I count it independently of whether it has
			// been made public information or not.
			info.num_our_flores++
		} else {
			// It's an opponent; then I can count him ONLY if he has made it
			// public.
			has_not_sung_yet := slices.Index(p.Ronda.Envite.SinCantar, con_flor.Jugador.ID) != -1
			has_made_it_public := !has_not_sung_yet
			// Include if: teammate OR (opponent AND has made it public)
			if has_made_it_public {
				info.num_opp_flores++
			}
		}
	}
}

func (info *InfosetFloki) setEnvite(p *pdt.Partida) {
	info.envite_estado = int(p.Ronda.Envite.Estado)
}

func (info *InfosetFloki) setEnviteCantadoPorOurTeam(p *pdt.Partida, m *pdt.Manojo) {
	our_team := m.Jugador.Equipo
	info.envite_cantado_por_our_team = false

	if p.Ronda.Envite.Estado != pdt.DESHABILITADO &&
		p.Ronda.Envite.Estado != pdt.NOCANTADOAUN &&
		p.Ronda.Envite.CantadoPor != "" {
		info.envite_cantado_por_our_team = p.Ronda.Manojo(p.Ronda.Envite.CantadoPor).Jugador.Equipo == our_team
	}
}

func (info *InfosetFloki) setTruco(p *pdt.Partida) {
	info.truco_estado = int(p.Ronda.Truco.Estado)
}

func (info *InfosetFloki) setTrucoCantadoPorOurTeam(p *pdt.Partida, m *pdt.Manojo) {
	our_team := m.Jugador.Equipo
	info.truco_cantado_por_our_team = false

	if p.Ronda.Truco.Estado != pdt.NOGRITADOAUN &&
		p.Ronda.Truco.CantadoPor != "" {
		info.truco_cantado_por_our_team = p.Ronda.Manojo(p.Ronda.Truco.CantadoPor).Jugador.Equipo == our_team
	}
}

func (info *InfosetFloki) setManoEnJuego(p *pdt.Partida) {
	info.mano_en_juego = int(p.Ronda.ManoEnJuego)
}

func (info *InfosetFloki) setManosInfo(p *pdt.Partida, m *pdt.Manojo) {
	our_team := m.Jugador.Equipo

	info.max_poder = make([]int, 0, len(p.Ronda.Manos))
	info.is_from_us = make([]bool, 0, len(p.Ronda.Manos))

	for mano_idx := 0; mano_idx < len(p.Ronda.Manos); mano_idx++ {
		mano := p.Ronda.Manos[mano_idx]

		// Resultado
		info.resultado_mano = int(mano.Resultado)

		// Find highest card
		max_poder := -1
		is_from_us := false // if the highest is from us, that means we are leading

		for _, t := range mano.CartasTiradas {
			// Get the manojo and the actual card using the index
			poder := t.Carta.CalcPoder(p.Ronda.Muestra)
			if poder > max_poder {
				max_poder = poder
				is_from_us = p.Ronda.Manojo(t.Jugador).Jugador.Equipo == our_team
			}
		}

		info.max_poder = append(info.max_poder, max_poder)
		info.is_from_us = append(info.is_from_us, is_from_us)
	}
}

func (info *InfosetFloki) setHowManyAfterMeHaventFoldedYet(p *pdt.Partida, m *pdt.Manojo) {
	num_players := len(p.Ronda.Manojos)
	our_team := m.Jugador.Equipo
	info.num_opps_after_me_that_can_play = 0
	ix := p.Ronda.MIXS[m.Jugador.ID]
	end := (int(p.Ronda.ElMano) + num_players - 1) % num_players
	for count := 0; count < num_players; count++ {
		if p.Ronda.Manojos[ix].Jugador.Equipo != our_team &&
			!p.Ronda.Manojos[ix].SeFueAlMazo {
			info.num_opps_after_me_that_can_play++
		}

		if ix == end {
			break
		}
		ix = (ix + 1) % num_players
	}
}

func (info *InfosetFloki) setChi(
	p *pdt.Partida,
	manojo *pdt.Manojo,
	chi_i pdt.A,
	abs abs.IAbstraction,
) {
	n := abs.Len()
	counter := make([]int, n) // tamano fijo (num de buckets)

	for i := 0; i < 3; i++ {
		if cartaHabilitada := chi_i[i]; cartaHabilitada {
			c := manojo.Cartas[i]
			bucket := abs.Abstract(c, &p.Ronda.Muestra)
			counter[bucket]++
		}
	}

	resto := make([]int, len(chi_i[3:])) // tamano fijo
	for ix, v := range chi_i[3:] {
		if v {
			resto[ix] = 1
		} else {
			resto[ix] = 0
		}
	}

	info.Chi = append(counter, resto...)
}

func (info *InfosetFloki) ChiLen() int {
	chi_len := 0
	for _, a := range info.Chi {
		if a > 0 {
			chi_len += 1
		}
	}
	return chi_len
}

func (info *InfosetFloki) Iterable(
	p *pdt.Partida,
	m *pdt.Manojo,
	aixs pdt.A, // array de 15 acciones (bool): 3 cartas + 12 "jugadas"
	abs abs.IAbstraction,
) []pdt.IJugada {

	res := make([]pdt.IJugada, 0, 3+12)

	n := abs.Len()

	counter := make([][]int, n)
	for i := 0; i < 3; i++ {
		noLaTiro := !m.Tiradas[i]
		laPuedeTirar := aixs[i]
		cartaHabilitada := noLaTiro && laPuedeTirar
		if cartaHabilitada {
			c := m.Cartas[i]
			bucket := abs.Abstract(c, &p.Ronda.Muestra)
			if counter[bucket] == nil {
				counter[bucket] = []int{i}
			} else {
				counter[bucket] = append(counter[bucket], i)
			}
		}
	}

	for _, bucketCount := range counter {
		if len(bucketCount) > 0 {

			randomIx := rand.Intn(len(bucketCount))
			cartaIx := bucketCount[randomIx]
			res = append(
				res,
				&pdt.TirarCarta{
					JID:   m.Jugador.ID,
					Carta: *m.Cartas[cartaIx],
				},
			)
		}
	}

	restoDeAcciones := info.Chi[n:]
	for ix, a := range restoDeAcciones {
		if a > 0 {
			canonicalAix := ix + 3

			res = append(
				res,
				pdt.ToJugada(
					p,
					p.Ronda.JIX(m.Jugador.ID),
					canonicalAix))
		}
	}

	return res
}

func (info *InfosetFloki) HashBytes(h hash.Hash) []byte {
	h.Reset()
	hsep := []byte(sep)

	{
		bs, _ := json.Marshal(info.we_are_starters)
		h.Write(bs)
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.muestra_valor)
		h.Write(bs)
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.cards_i_still_own)
		h.Write(bs)
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.available_teammate_cards)
		h.Write(bs)
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.all_our_cards)
		h.Write(bs)
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.num_our_flores)
		h.Write(bs)
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.num_opp_flores)
		h.Write(bs)
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.envite_estado)
		h.Write(bs)
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.envite_cantado_por_our_team)
		h.Write(bs)
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.truco_estado)
		h.Write(bs)
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.truco_cantado_por_our_team)
		h.Write(bs)
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.mano_en_juego)
		h.Write(bs)
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.resultado_mano)
		h.Write(bs)
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.num_opps_after_me_that_can_play)
		h.Write(bs)
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.max_poder)
		h.Write(bs)
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.is_from_us)
		h.Write(bs)
		h.Write(hsep)
	}
	{
		bs, _ := json.Marshal(info.Chi)
		h.Write(bs)
	}

	return h.Sum(nil)
}

func (info *InfosetFloki) Hash(h hash.Hash) string {
	return hex.EncodeToString(info.HashBytes(h))
}

func (info *InfosetFloki) Dump(indent bool) string {
	var bs []byte = nil
	if indent {
		bs, _ = json.MarshalIndent(info, "", "\t")
	} else {
		bs, _ = json.Marshal(info)
	}
	return string(bs)
}

func infosetFlokiFactory(a abs.IAbstraction) InfosetBuilder {
	return func(
		p *pdt.Partida,
		m *pdt.Manojo,
		msgs []enco.IMessage,
	) Infoset {
		info := &InfosetFloki{
			Vision: m.Jugador.ID, // <- tiene motivos solo depurativos
		}
		chi_i := pdt.GetA(p, m)

		info.setWeAreStarters(p, m)
		info.setMuestraValor(p)
		info.setCardsIStillOwn(p, m, a)
		info.setTeammatesAfterMeCardsStillOwned(p, m, a)
		info.setAllOurCards(p, m, a)
		info.setFlores(p, m)
		info.setEnvite(p)
		info.setEnviteCantadoPorOurTeam(p, m)
		info.setTruco(p)
		info.setTrucoCantadoPorOurTeam(p, m)
		info.setManoEnJuego(p)
		info.setManosInfo(p, m)
		info.setHowManyAfterMeHaventFoldedYet(p, m)

		info.setChi(p, m, chi_i, a)
		return info
	}
}
