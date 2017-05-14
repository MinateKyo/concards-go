package termboxgui

import (
	"github.com/alanxoc3/concards-go/algs"
	"github.com/alanxoc3/concards-go/deck"
	"github.com/alanxoc3/concards-go/termhelp"
	termbox "github.com/nsf/termbox-go"
)

func TermBoxRun(d deck.Deck, cfg *termhelp.Config, ctrl []*deck.DeckControl) error {
	err := termbox.Init()
	if err != nil {
		return err
	}

	defer termbox.Close()
	defer termbox.Sync()
	termbox.SetInputMode(termbox.InputAlt)
	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)
	flush()

	data = make([]byte, 0, 64)
	const coldef = termbox.ColorDefault

	card_shown := false
	help_mode := false
	quit_mode := false
	finished_editing := false

	for len(d) > 0 {
		draw_screen(d, help_mode, card_shown, finished_editing)
		finished_editing = false

		inp, err := update_input()
		if err != nil {
			return err
		}

		if !quit_mode {
			if inp == "q" {
				update_stat_msg("Quit concards? [y/N]", termbox.ColorYellow)
				quit_mode = true
			} else if inp == "h" {
				help_mode = !help_mode
			} else if inp == "w" {
				err = deck.WriteDeckControls(ctrl)
				if err != nil {
					update_stat_msg(err.Error(), termbox.ColorRed)
				} else {
					update_stat_msg("Cards were written.", termbox.ColorYellow)
				}
			} else if !help_mode {
				if inp == "1" {
					update_stat_msg_and_card(d.Top(), algs.NO)
					d = append(d[1:], d[0]) // top to bottom
					card_shown = false
				} else if inp == "2" {
					update_stat_msg_and_card(d.Top(), algs.IDK)
					d = append(d[1:], d[0]) // top to bottom
					card_shown = false
				} else if inp == "3" {
					update_stat_msg_and_card(d.Top(), algs.YES)
					d = d[1:]
					card_shown = false
				} else if inp == "e" {
					err := deck.EditCard(cfg.Editor, d.Top(), "You may ONLY EDIT the cards here.\nREARRANGING, DELETING, or ADDING cards WILL CORRUPT your files.")

					if err != nil {
						update_stat_msg(err.Error(), termbox.ColorRed)
					} else {
						update_stat_msg("Card was successfully edited.", termbox.ColorYellow)
					}
					finished_editing = true
				} else if inp == " " || inp == "\r" {
					if d.Top().HasAnswer() {
						card_shown = !card_shown
					} else {
						update_stat_msg("This card has no answer.", termbox.ColorRed)
					}
				}
			}
		} else {
			if inp == "y" {
				return nil
			} else {
				quit_mode = false
				update_stat_msg("", coldef)
			}
		}

	}

	return nil
}

func draw_screen(d deck.Deck, help_mode, card_shown, finished_editing bool) {
	if finished_editing {
		termbox.Sync()
	}

	termbox.Clear(termbox.ColorBlack, termbox.ColorBlack)

	if help_mode {
		display_help_mode(termbox.ColorCyan)
	} else {
		display_card_mode(d.Top(), card_shown)
	}

	tbprint_statusbar(d)
	tbprint_stat_msg()
	flush()
}
