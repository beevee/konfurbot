package telegram

import (
	"sync"
	"time"

	"github.com/looplab/fsm"
	"github.com/tucnak/telebot"
	"gopkg.in/tomb.v2"

	"github.com/beevee/konfurbot"
)

// Bot handles interactions with Telegram users
type Bot struct {
	ScheduleStorage   konfurbot.ScheduleStorage
	TelegramToken     string
	Timezone          *time.Location
	Logger            konfurbot.Logger
	telebot           *telebot.Bot
	chatStateMachines map[int64]*fsm.FSM
	chatStateLock     sync.RWMutex
	tomb              tomb.Tomb
}

// Start initializes Telegram API connections
func (b *Bot) Start() error {
	b.chatStateMachines = make(map[int64]*fsm.FSM, 0)

	var err error
	b.telebot, err = telebot.NewBot(b.TelegramToken)
	if err != nil {
		return err
	}

	messages := make(chan telebot.Message)
	b.telebot.Listen(messages, 1*time.Second)

	b.tomb.Go(func() error {
		for {
			select {
			case message := <-messages:
				b.handleMessage(message)
			case <-b.tomb.Dying():
				return nil
			}
		}
	})

	return nil
}

// Stop gracefully stops Telegram API connections
func (b *Bot) Stop() error {
	b.tomb.Kill(nil)
	return b.tomb.Wait()
}

func (b *Bot) handleMessage(message telebot.Message) {
	b.Logger.Log("msg", "message received", "firstname", message.Sender.FirstName,
		"lastname", message.Sender.LastName, "username", message.Sender.Username,
		"chatid", message.Chat.ID, "command", message.Text)

	var stateMachine *fsm.FSM

	b.chatStateLock.RLock()
	stateMachine, ok := b.chatStateMachines[message.Chat.ID]
	b.chatStateLock.RUnlock()

	if !ok {
		b.Logger.Log("msg", "no state machine exists, creating new", "chatid", message.Chat.ID)
		stateMachine = initStateMachine()

		b.chatStateLock.Lock()
		b.chatStateMachines[message.Chat.ID] = stateMachine
		b.chatStateLock.Unlock()

		stateMachine.Event(greetCommand, message.Chat, b)
		return
	}

	b.Logger.Log("msg", "state machine exists, attempting transition", "chatid", message.Chat.ID,
		"current_state", stateMachine.Current(), "command", message.Text)
	err := stateMachine.Event(message.Text, message.Chat, b)
	if err != nil && err.Error() != "no transition" {
		b.Logger.Log("msg", "something is wrong with the transition, will return to start",
			"current_state", stateMachine.Current(), "chatid", message.Chat.ID, "error", err)
		stateMachine.Event(unknownCommand, message.Chat, b)
	}
}
