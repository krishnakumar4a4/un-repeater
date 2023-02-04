package menu

import (
	"log"

	"github.com/caseymrm/menuet"
	"github.com/krishnakumar4a4/un-repeater/worker"
)

func GetMenuItems(lister *worker.ScriptLister) []menuet.MenuItem {
	log.Printf("re-evaluating Get menuitems \n")

	_, startScripts := lister.ListStartScripts()
	_, stopScripts := lister.ListStopScripts()

	startScriptMenuItems := createMenuItems(startScripts, lister)
	stopScriptMenuItems := createMenuItems(stopScripts, lister)

	menuItems := make([]menuet.MenuItem, 0, len(startScripts)+len(stopScripts))
	menuItems = append(menuItems, menuet.MenuItem{
		Text: "Toggle Start Scripts",
		Children: func() []menuet.MenuItem {
			return startScriptMenuItems
		},
	})
	menuItems = append(menuItems, menuet.MenuItem{
		Text: "Toggle Stop Scripts",
		Children: func() []menuet.MenuItem {
			return stopScriptMenuItems
		},
	})
	return menuItems
}

func createMenuItems(scripts []string, lister *worker.ScriptLister) []menuet.MenuItem {
	menuItems := make([]menuet.MenuItem, 0, len(scripts))
	for _, script := range scripts {
		scriptName := script
		enabledState := true
		if v, ok := lister.GetToggleScripts()[scriptName]; ok {
			enabledState = v
		}
		menuItems = append(menuItems, menuet.MenuItem{
			Text:  scriptName,
			State: enabledState,
			Clicked: func() {
				log.Printf("toggling script %s to %v \n", scriptName, enabledState)
				lister.ToggleScript(scriptName)
			},
		})
	}
	return menuItems
}
