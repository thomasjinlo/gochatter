package client

//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"io"
//	"log"
//	"net/http"
//
//	"github.com/gorilla/websocket"
//)

//type Channel struct {
//    id string
//    name string
//    users []string
//
//    conn *websocket.Conn
//}

//func NewChannel(displayName string, channelBody ChannelResponse) *Channel {
//    dialer := websocket.DefaultDialer
//    header := http.Header{}
//    header.Add("DisplayName", displayName)
//    header.Add("ChannelId", channelBody.Id)
//
//    log.Println("Joining channel", fmt.Sprintf("ws://localhost:8443", channelBody.Id))
//    conn, _, err := dialer.Dial("ws://localhost:8443/", header)
//    if err != nil {
//        log.Println("CONNECT error", err)
//    }
//
//    channel := &Channel{
//        id: channelBody.Id,
//        name: channelBody.Name,
//        users: channelBody.Users,
//
//        conn: conn,
//    }
//
//    return channel
//}
//
//
//func (channel *Channel) ChannelName() string {
//    return channel.name
//}
//
//func (channel *Channel) GetUsers() []string {
//    return channel.users
//}
//
//func (channel *Channel) GetChatHistory() string {
//    return "line 1\nline2\nline3\nline4"
//}
//
//func (channel *Channel) Broadcast(message string) {
//    err := channel.conn.WriteMessage(websocket.BinaryMessage, []byte(message))
//    if err != nil {
//        log.Println(err)
//        return
//    }
//}
//
//func (channel *Channel) ReceiveFromServer(handleMessages func(message []byte)) {
//    for {
//        _, messageBytes, err := channel.conn.ReadMessage()
//        if err != nil {
//            log.Println(err)
//            return
//        }
//
//        handleMessages(messageBytes)
//    }
//}
//
//type ChannelResponse struct {
//    Id string
//    Name string
//    Users []string
//}
//
//type GetChannelsResponse struct {
//    Channels []ChannelResponse
//}
//
//func GetChannels(displayName string) []*Channel {
//    var channels []*Channel
//
//    //response, err := http.Get("https://gochatter.app/channels")
//    response, err := http.Get("http://localhost:8080/channels")
//    if err != nil {
//        log.Println("GET REQUEST ERR", err)
//        return channels
//    }
//
//    body, err := io.ReadAll(response.Body)
//    if err != nil {
//        log.Println(err)
//        return channels
//    }
//
//    var channelResponse GetChannelsResponse
//
//    err = json.Unmarshal(body, &channelResponse)
//    if err != nil {
//        log.Println(err)
//        return channels
//    }
//    log.Println(channelResponse)
//
//    for _, channelBody := range channelResponse.Channels {
//        channels = append(channels, NewChannel(displayName, channelBody))
//    }
//
//    return channels
//}
//
//func CreateChannel(name string) {
//    client := &http.Client{}
//
//    type RequestBody struct {
//        name string
//    }
//
//    requestBody := &RequestBody{name: name}
//    jsonData, err := json.Marshal(requestBody)
//    if err != nil {
//        log.Println(err)
//        return
//    }
//    body := bytes.NewReader(jsonData)
//    req, err := http.NewRequest("POST", "https://gochatter.app/channels", body)
//    if err != nil {
//        log.Println(err)
//        return
//    }
//
//    res, err := client.Do(req)
//    if err != nil {
//        log.Println(err)
//        return
//    }
//
//    log.Println(res)
//}
//
//type 
