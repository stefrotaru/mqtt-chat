package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type message struct {
	user string
	text string
	sentAt time.Time
}

type model struct {
	topic	  string
	name	  string
	client    mqtt.Client
	messages  []message
	textInput textinput.Model
	cursor    int
}

var received = make(chan mqtt.Message)

func initialModel(topic string, name string, client mqtt.Client) model {
	ti := textinput.New()
	ti.Placeholder = "Type your message..."
	ti.Focus()
	ti.Width = 80

	m := model{
		topic: topic,
		name: name,
		client: client,
		messages: []message{},
		textInput: ti,
		cursor:  0,
	}

	// message handler adds messages to a global channel
	var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		received <- msg
	}

	if token := client.Subscribe(topic + "#", 0, messagePubHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	return m
}

func (m model) Init() tea.Cmd {
    return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			token := m.client.Publish(m.topic + m.name, 0, false, m.textInput.Value())
			token.Wait()
			m.textInput.SetValue("")
			return m, nil
		}
    }

	select {
	case msg := <-received:
		m.messages = append(m.messages, message{
			user: strings.Split(msg.Topic(), "/")[1],
			text: string(msg.Payload()),
			sentAt: time.Now(),
		})
	default:
	}
	m.textInput, cmd = m.textInput.Update(msg)
    return m, cmd
}

func (m model) View() string {
    // The header
    s := fmt.Sprintf("--- %s ---\n\n", m.topic)

    // Iterate over messages
    for _, msg := range m.messages {
        // Render the row
        s += fmt.Sprintf("%s [%s]: %s\n", msg.user, msg.sentAt.Format("15:04:05"), msg.text)
    }

    // The footer
    s += "\n----" + strings.Repeat("-", len(m.topic)) + "----\n"
	s += m.textInput.View() + "\n"

    // Send the UI for rendering
    return s
}

func main() {
	var displayName string
	if len(os.Args) > 1 {
		displayName = os.Args[1]
	} else {
		displayName = "go-simple"
	}

	opts := mqtt.NewClientOptions().AddBroker("0.0.0.0:1883")
	opts.SetClientID(displayName)
	opts.SetCleanSession(true)
	
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	p := tea.NewProgram(initialModel("topicul_de_miercuri_seara/", displayName, c))
    if _, err := p.Run(); err != nil {
        fmt.Printf("Error: %v", err)
        os.Exit(1)
    }
}
