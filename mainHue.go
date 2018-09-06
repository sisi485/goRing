package main

import (
	"fmt"
	"os"
	"time"
	"ring/hue"
	"github.com/rhymen/go-whatsapp"
	"encoding/gob"
	"github.com/Baozisoftware/qrcode-terminal-go"
)


const (
	AWAITTIME      = 800
	TRANSITIONTIME = float64(AWAITTIME) / 1000 / 0.1
	//SENDMSGTO 	   = "491787807924"
	SENDMSGTO 	   = "4915154689867"
	MSGTEXT 	   = "mhhh et klingelt bruder"
)


func main() {

	err := sendMsg(SENDMSGTO)

	if err != nil {
		fmt.Fprintf(os.Stderr, "can not send msg: %v\n", err)
		return
	}

	return

	lights, err := hue.GetLights()
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not get lgihts: %v\n", err)
		return
	}

	keys := make([]string, 0, len(lights))
	for k := range lights {
		keys = append(keys, k)
	}

	scene, err := hue.CreateScene("test", keys)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not create scene: %v\n", err)
		return
	}

	lightsColor := map[string]interface{}{
		"on":             true,
		"sat":            255,
		"hue":            34000,
		"bri":            255,
		"transitiontime": TRANSITIONTIME,
	}

	lightsOn := map[string]interface{}{
		"on":             true,
		"sat":            255,
		"hue":            34000,
		"bri":            255,
		"transitiontime": TRANSITIONTIME,
	}

	lightsOff := map[string]interface{}{
		"bri":            0,
		"transitiontime": TRANSITIONTIME,
	}

	id := scene[0].(map[string]interface{})["success"].(map[string]interface{})["id"].(string)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not convert id: %v\n", err)
		return
	}

	lightsDefault := map[string]interface{}{
		"scene":          id,
		"transitiontime": TRANSITIONTIME,
	}

	hue.SetGroup(1, lightsOff)
	<-time.After(time.Millisecond * AWAITTIME)
	hue.SetGroup(1, lightsColor)
	<-time.After(time.Millisecond * AWAITTIME)
	hue.SetGroup(1, lightsOff)
	<-time.After(time.Millisecond * AWAITTIME)
	hue.SetGroup(1, lightsOn)
	<-time.After(time.Millisecond * AWAITTIME)
	hue.SetGroup(1, lightsOff)
	<-time.After(time.Millisecond * AWAITTIME)
	hue.SetGroup(1, lightsDefault)

	_, err = hue.DeleteScene(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can not delete scene: %v\n", err)
		return
	}

	fmt.Println("im done, closing..")
}

func sendMsg(sendMsgTo string) error {

	wac, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating connection: %v\n", err)
		return err
	}

	err = login(wac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error logging in: %v\n", err)
		return err
	}

	<-time.After(3 * time.Second)

	msg := whatsapp.TextMessage{
		Info: whatsapp.MessageInfo{
			RemoteJid: fmt.Sprintf("%s@s.whatsapp.net", sendMsgTo),
		},
		Text: MSGTEXT,
	}

	err = wac.Send(msg)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error sending message: %v", err)
	}

	return err

}

func login(wac *whatsapp.Conn) error {
	//load saved session
	session, err := readSession()
	if err == nil {
		//restore session
		session, err = wac.RestoreSession(session)
		if err != nil {
			return fmt.Errorf("restoring failed: %v\n", err)
		}
	} else {
		//no saved session -> regular login
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()
		session, err = wac.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %v\n", err)
		}
	}

	//save session
	err = writeSession(session)
	if err != nil {
		return fmt.Errorf("error saving session: %v\n", err)
	}
	return nil
}

func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}
	return session, nil
}

func writeSession(session whatsapp.Session) error {
	file, err := os.Create(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(session)
	if err != nil {
		return err
	}
	return nil
}