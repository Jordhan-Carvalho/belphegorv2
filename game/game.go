package game

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jordhan-carvalho/belphegorv2/interfaces"
	"github.com/jordhan-carvalho/belphegorv2/sound"
)

var gameTime = 0

func StartListeningToGame(gEventC chan interfaces.GameEvents, vc *discordgo.VoiceConnection, gDone chan bool) {
	fmt.Println("Start Listening to the game input")
	buyWardsLastCall := 0
	for {
		select {
		case <-gDone:
			return
		case event := <-gEventC:
			// since the throttle on the post request is 1, we get the same ClockTime sometimes
			fmt.Println("Inside event receiver on StartListeningToGame")
			if gameTime != event.Map.ClockTime && event.Map.GameState == "DOTA_GAMERULES_STATE_GAME_IN_PROGRESS" {
				gameTime = event.Map.ClockTime
				wardsPurchaseCd := event.Map.WardPurchaseCooldown
				fmt.Println("Game clockTime", event.Map.ClockTime)
				fmt.Println("Game ward cooldown", event.Map.WardPurchaseCooldown)

				checkBountyRunes(vc)
				checkMidRunes(vc)
				checkForStack(vc)
				checkForSmokeInShop(vc)
				checkForWards(vc, wardsPurchaseCd, &buyWardsLastCall)
			}
		}
	}
}

func StartRoshanAndAegisTimers(gEventC chan interfaces.GameEvents, killedTime int, vc *discordgo.VoiceConnection, isTimeRunning *bool, eventReceivers *int) {
	roshanMinSpawnWarningTime := 470
	// roshanMinSpawnDelay := 480
	roshanMaxSpawnWarningTime := 659
	// roshanMaxSpawnDelay := 660
	aegis2minWarnTime := 180
	aegies30sWarnTime := 270
	// aegisDelay := 300
	fmt.Println("Roshan timer started")

myLoop:
	for {
		select {
		case event := <-gEventC:
			fmt.Println("Inside Roshan event receiver")
			gameTime = event.Map.ClockTime

			if killedTime+aegis2minWarnTime == gameTime {
				go sound.PlaySpecificSound(vc, "aegis-2min.dca")
			}

			if killedTime+aegies30sWarnTime == gameTime {
				go sound.PlaySpecificSound(vc, "aegis-30s.dca")
			}

			if killedTime+roshanMinSpawnWarningTime == gameTime {
				go sound.PlaySpecificSound(vc, "roshan-min.dca")
			}

			if killedTime+roshanMaxSpawnWarningTime <= gameTime {
				fmt.Println("Entrou killedtime, roshanMax, gametime", killedTime, roshanMaxSpawnWarningTime, gameTime)
				go sound.PlaySpecificSound(vc, "roshan-max.dca")
				*eventReceivers--
				*isTimeRunning = false
				break myLoop
			}
		}
	}
	return

}

// smoke logic, will check if any on invertory if its any it will start a 7 min count
// smoke at every 7 minutes... it starts at 2... max stack is 3... after 7 min you will have max stack
func checkForSmokeInShop(vc *discordgo.VoiceConnection) {
	smokeWarnTime := 415
	smokeDelay := 420

	if (gameTime-smokeWarnTime)%smokeDelay == 0 {
		go sound.PlaySpecificSound(vc, "smoke.dca")
	}
}

func checkForStack(vc *discordgo.VoiceConnection) {
	stackGameTime := 48 // ingame time to stack
	stackDelay := 60    // interval between stack

	if (gameTime-stackGameTime)%stackDelay == 0 {
		go sound.PlaySpecificSound(vc, "stack.dca")
	}
}

func checkBountyRunes(vc *discordgo.VoiceConnection) {
	bountyRunesGameTime := 173
	bountyRunesDelay := 180

	if (gameTime-bountyRunesGameTime)%bountyRunesDelay == 0 {
		go sound.PlaySpecificSound(vc, "bounty-rune.dca")
	}
}

func checkMidRunes(vc *discordgo.VoiceConnection) {
	midRunesGameTime := 112
	midRunesDelay := 120

	if (gameTime-midRunesGameTime)%midRunesDelay == 0 {
		go sound.PlaySpecificSound(vc, "mid-rune.dca")
	}
}

func checkForWards(vc *discordgo.VoiceConnection, wardCd int, buyWardsLastCall *int) {
  timeBetweenCalls := 30
	if wardCd == 0 && (*buyWardsLastCall+timeBetweenCalls) <= gameTime {
		go sound.PlaySpecificSound(vc, "wards.dca")
    *buyWardsLastCall = gameTime
	}
}

// TODO tower in deny range
