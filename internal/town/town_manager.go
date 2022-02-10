package town

import (
	"github.com/hectorgimenez/koolo/internal/action"
	"github.com/hectorgimenez/koolo/internal/config"
	"github.com/hectorgimenez/koolo/internal/game/data"
	"github.com/hectorgimenez/koolo/internal/helper"
	"github.com/hectorgimenez/koolo/internal/hid"
	"time"
)

type Manager struct {
	towns       map[data.Area]Town
	cfg         config.Config
	pf          helper.PathFinder
	shopManager ShopManager
}

func NewTownManager(cfg config.Config, pf helper.PathFinder, shopManager ShopManager) Manager {
	return Manager{
		towns: map[data.Area]Town{
			data.AreaHarrogath: A5{},
		},
		cfg:         cfg,
		pf:          pf,
		shopManager: shopManager,
	}
}

func (tm Manager) BuyPotionsAndTPs(area data.Area, buyTPs bool) {
	t := tm.getTownByArea(area)
	tm.pf.InteractToNPC(t.RefillNPC())
	tm.openTradeMenu()
	tm.shopManager.buyPotsAndTPs(buyTPs)
}

func (tm Manager) Repair(area data.Area) {
	t := tm.getTownByArea(area)
	tm.pf.InteractToNPC(t.RepairNPC())
	tm.openTradeMenu()
	action.Run(
		action.NewMouseDisplacement(int(float32(hid.GameAreaSizeX)/3.52), int(float32(hid.GameAreaSizeY)/1.37), time.Millisecond*850),
		action.NewMouseClick(hid.LeftButton, time.Millisecond*1300),
	)
}

func (tm Manager) ReviveMerc(area data.Area) {
	t := tm.getTownByArea(area)
	tm.pf.InteractToNPC(t.MercContractorNPC())
	tm.openTradeMenu()
}

func (tm Manager) Stash() {
	for _, o := range data.Status().Objects {
		if o.Name == "Bank" {
			tm.pf.InteractToObject(o)
			tm.stashAllItems()
		}
	}

}

func (tm Manager) WPTo(act int, area int) {
	for _, o := range data.Status().Objects {
		if o.IsWaypoint() {
			tm.pf.InteractToObject(o)
			return
		}
	}
}

func (tm Manager) openTradeMenu() {
	if data.Status().OpenMenus.NPCInteract {
		action.Run(action.NewKeyPress("down", time.Millisecond*150), action.NewKeyPress("enter", time.Millisecond*500))
	}
}

func (tm Manager) getTownByArea(area data.Area) Town {
	return tm.towns[area]
}