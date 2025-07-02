package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math"
	"net"
	"net/http"
	"net/url"
	"os"
	"slices"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
)

// This function serves
func getPeerList(peerID, infohash string) ([]string, error) {
	if len(infohash) != 20 || len(peerID) != 20 {
		return []string{}, nil
	}

	newURL, _ := url.Parse("http://bittorrent-test-tracker.codecrafters.io/announce")
	values := newURL.Query()
	values.Add("info_hash", infohash)
	values.Add("peer_id", peerID)
	values.Add("port", "6881")
	values.Add("uploaded", "0")
	values.Add("downloaded", "0")
	values.Add("left", "10000")
	values.Add("compact", "1")
	newURL.RawQuery = values.Encode()

	yy := newURL.String()
	resp, _ := http.Get(yy)
	rawResponse, _ := io.ReadAll(resp.Body)

	type hoho struct {
		MinInterval int    `bencode:"min interval"`
		Peers       []byte `bencode:"peers"`
		Complete    int    `bencode:"complete"`
		Interval    int    `bencode:"interval"`
		Incomplete  int    `bencode:"incomplete"`
	}
	var ppp hoho
	bencode.Unmarshal(rawResponse, &ppp)

	peers := []string{}
	idx := 0
	for idx < len(ppp.Peers) {
		rawSingleIPAddr := ppp.Peers[idx : idx+6]
		singleIPAddr := fmt.Sprintf(
			"%d.%d.%d.%d:%d",
			rawSingleIPAddr[0],
			rawSingleIPAddr[1],
			rawSingleIPAddr[2],
			rawSingleIPAddr[3],
			binary.BigEndian.Uint16(rawSingleIPAddr[4:6]),
		)
		peers = append(peers, singleIPAddr)
		idx += 6
	}
	return peers, nil
}

type TorrentFile struct {
	Filename    string
	PieceCheck  []bool
	PieceLength int
	Pieces      [][]byte
	Length      int
}

func NewTorrentFile(filename string, length, pieceLength int, pieces []byte) TorrentFile {
	idx := 0
	splittedPieces := [][]byte{}
	for idx < len(pieces) {
		splittedPieces = append(splittedPieces, pieces[idx:(idx+20)])
		idx += 20
	}
	return TorrentFile{
		Filename:    filename,
		PieceLength: pieceLength,
		Length:      length,
		Pieces:      splittedPieces,
		PieceCheck:  slices.Repeat[[]bool, bool]([]bool{false}, len(splittedPieces)),
	}
}

func (tf *TorrentFile) ExecuteFileCheck() {
	for idx, _ := range tf.Pieces {
		tf.ExecutePieceCheck(idx)
	}
}

func (tf *TorrentFile) ExecutePieceCheck(idx int) bool {
	ff, _ := os.Open(tf.Filename)
	defer ff.Close()

	rawPieceContent := make([]byte, tf.PieceLength)
	ff.ReadAt(rawPieceContent, int64(idx*tf.PieceLength))
	ss2 := sha1.New()
	ss2.Write(rawPieceContent)
	oo2 := ss2.Sum(nil)
	contentCheck := bytes.Equal(oo2, tf.Pieces[idx])
	tf.PieceCheck[idx] = contentCheck
	return contentCheck
}

func main() {
	yahoo, _ := os.ReadFile("zzz.torrent")
	type aaa struct {
		Length      int    `bencode:"length"`
		Name        string `bencode:"name"`
		PieceLength int    `bencode:"piece length"`
		Pieces      []byte `bencode:"pieces"`
	}
	type zzz struct {
		Announce     string     `bencode:"announce,omitempty"`
		AnnounceList [][]string `bencode:"announce-list"`
		Comment      string     `bencode:"comment"`
		CreatedBy    string     `bencode:"created by"`
		CreationDate int        `bencode:"creation date"`
		Encoding     string     `bencode:"encoding"`
		Info         aaa        `bencode:"info"`
	}
	var zz zzz
	bencode.Unmarshal(yahoo, &zz)

	// fmt.Printf("%+v\n", zz)

	val, _ := bencode.Marshal(zz.Info)
	ss := sha1.New()
	ss.Write(val)
	oo := ss.Sum(nil)
	fmt.Printf("%x\n", oo)

	tf := NewTorrentFile("hoho.mkv", zz.Info.Length, zz.Info.PieceLength, zz.Info.Pieces)
	fmt.Println(tf.ExecutePieceCheck(0))

	tcpListener, _ := os.LookupEnv("BT_PORT")
	httpListener, _ := os.LookupEnv("HTTP_PORT")
	peerURL, _ := os.LookupEnv("PEER_URL")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", tcpListener))
	if err != nil {
		panic("Unable to listen to port")
	}
	defer listener.Close()

	go handler(listener)

	hoho := TorrentMessage{}
	zpa := hoho.allConvertBitField(tf.PieceCheck)
	fmt.Printf("%v\n", zpa)

	http.Handle("/foo", fooHandler{PeerURL: peerURL})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", httpListener), nil))

	// pp, _ := hex.DecodeString("d69f91e6b2ae4c542468d1073a71d4ea13879a7f")
	// peerID := "11111111111111111111"

	// peers, _ := getPeerList(peerID, string(pp))
	// fmt.Printf("Peers: %v\n", peers)

	// conn, err := net.Dial("tcp", peers[0])
	// if err != nil {
	// 	fmt.Println("Error connecting:", err)
	// 	return
	// }
	// defer conn.Close()

	// pp2 := TorrentMessage{peerID: []byte(peerID), infohash: pp}
	// yy := TorrentHandler{conn: conn}
	// yy.handshake(pp2.handshake())
	// yy.handshake(pp2.Interested())
	// yy.handshake(pp2.Request())
	// fmt.Printf("%x", zz.Info.Pieces[0:20])
}

type fooHandler struct {
	PeerURL string
}

func (f fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("start foo handler")
	defer slog.Info("end foo handler")

	conn, err := net.Dial("tcp", f.PeerURL)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	conn.Write([]byte("Testing this crappy thing"))
	conn.Write([]byte(io.EOF.Error()))
	tmp := make([]byte, 1024)
	conn.Read(tmp)
	fmt.Println(string(tmp))
}

func handler(listener net.Listener) {
	slog.Info("Start listener function")
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Handle error (e.g., connection closed)
			continue
		}
		go handleConnection(conn) // Handle connection in a new goroutine
	}
}

func handleConnection(conn net.Conn) {
	slog.Info("Handling connection")
	// make a temporary bytes var to read from the connection
	tmp := make([]byte, 1024)
	// make 0 length data bytes (since we'll be appending)
	conn.Read(tmp)
	slog.Info(string(tmp))

	// loop through the connection stream, appending tmp to data
	conn.Write([]byte("testing"))
}

func mainzz() {
	pp, _ := hex.DecodeString("13426974546f7272656e742070726f746f636f6c0000000000000004d69f91e6b2ae4c542468d1073a71d4ea13879a7f2d524e302e302e302d9d5e0226875575c6233a130000000205e0")
	tt := NewTorrentMessage(pp)
	fmt.Printf("%v", tt)

	zz := tt.Request(0, 0)
	fmt.Println(hex.EncodeToString(zz))

}

// https://wiki.theory.org/BitTorrentSpecification#have:_.3Clen.3D0005.3E.3Cid.3D4.3E.3Cpiece_index.3E
type TorrentMessage struct {
	peerID        []byte
	infohash      []byte
	messageLength int32
	messageType   int8
}

func NewTorrentMessage(raw []byte) TorrentMessage {
	// First 20 bytes is protocol info
	// Next 8 bytes is extension
	// Next 20 bytes is infohash
	// Next 20 bytes is peerid
	// Next 4 bytes Length of message
	// Next 1 byte is message id
	// The rest is contents
	rawMessageID := raw[68:72]
	rawMessageType := raw[72]
	messageLength := int32(binary.BigEndian.Uint32(rawMessageID))
	return TorrentMessage{
		infohash:      raw[28:48],
		peerID:        raw[48:68],
		messageLength: messageLength,
		messageType:   int8(rawMessageType),
	}
}

func (t TorrentMessage) String() string {
	output := ""
	output += fmt.Sprintf("Infohash: %+x\n", t.infohash)
	output += fmt.Sprintf("Peerid: %+x\n", t.peerID)
	output += fmt.Sprintf("MessageLength: %d\n", t.messageLength)
	output += fmt.Sprintf("MessageType: %d\n", t.messageType)
	return output
}

func (t TorrentMessage) convertToBytes(num int) []byte {
	// Int32 - refers to number represented by 32 bits
	// 1 byte represents 8 bits
	bBigEndian := make([]byte, 4)
	binary.BigEndian.PutUint32(bBigEndian, uint32(num))
	return bBigEndian
}

func (t TorrentMessage) allReverseBitField(bb []byte) []bool {
	rawBools := []bool{}
	for _, v := range bb {
		rawBools = t.reverseBitfield(v)
	}
	return rawBools
}

func (t TorrentMessage) reverseBitfield(b byte) []bool {
	aa := []bool{}
	rawVal := int(b)
	idx := 7
	for idx >= 0 {
		tempVal := rawVal - int(math.Pow(2, float64(idx)))
		if tempVal >= 0 {
			aa = append(aa, true)
			rawVal = tempVal
		} else {
			aa = append(aa, false)
		}
		idx -= 1
	}
	return aa
}

func (t TorrentMessage) allConvertBitField(items []bool) []byte {
	idx := 0
	raw := []byte{}
	for idx < len(items) {
		var singleBitfield byte
		if idx+8 > len(items) {
			singleBitfield = t.convertBitfield(items[idx:])
		} else {
			singleBitfield = t.convertBitfield(items[idx:(idx + 8)])
		}
		idx += 8
		raw = append(raw, singleBitfield)
	}
	return raw
}

func (t TorrentMessage) convertBitfield(items []bool) byte {
	if len(items) > 8 {
		items = items[0:8]
	}

	// Assume number of items is 8
	bitValue := 0.0
	for idx, val := range items {
		if val {
			bitValue += math.Pow(2, float64(7-idx))
		}
	}
	return byte(bitValue)
}

func (t TorrentMessage) Handshake() []byte {
	aa := []byte{}
	// First part is defining the protocol
	aa = append(aa, byte(19))
	aa = append(aa, []byte("BitTorrent protocol")...)
	// Second part is extension section
	aa = append(aa, []byte{0, 0, 0, 0, 0, 0, 0, 0}...)
	// Third part is infohash
	aa = append(aa, t.infohash...)
	// Fourth part is peerID to be passed to the other peer
	aa = append(aa, t.peerID...)
	return aa
}

// emptyMessage - return message via bittorent protocol that doesn't have message
func (t TorrentMessage) emptyMessage(messageType int) []byte {
	message := []byte{byte(messageType)}
	messageLength := t.convertToBytes(len(message))
	output := append(messageLength, message...)
	return output
}

func (t TorrentMessage) Choke() []byte {
	return t.emptyMessage(0)
}

func (t TorrentMessage) Unchoke() []byte {
	return t.emptyMessage(1)
}

func (t TorrentMessage) Interested() []byte {
	return t.emptyMessage(2)
}

func (t TorrentMessage) NotInterested() []byte {
	return t.emptyMessage(3)
}

func (t TorrentMessage) Have(index int) []byte {
	messageType := []byte{4}
	message := t.convertToBytes(index)
	initialOutput := append(messageType, message...)
	messageLength := t.convertToBytes(len(message))
	output := append(messageLength, initialOutput...)
	return output
}

func (t TorrentMessage) Request(indexPiece, blockOffset int) []byte {
	message := []byte{6}
	rawIndexPiece := t.convertToBytes(0)
	rawBlockOffset := t.convertToBytes(0)
	// TODO: For last piece, last block, it will be different
	blockLength := t.convertToBytes(16384)
	message = append(message, rawIndexPiece...)
	message = append(message, rawBlockOffset...)
	message = append(message, blockLength...)

	messageLength := t.convertToBytes(len(message))
	output := append(messageLength, message...)
	return output
}

type TorrentHandler struct {
	conn net.Conn
}

func (t *TorrentHandler) handshake(input []byte) []byte {
	_, err := t.conn.Write(input)
	if err != nil {
		fmt.Println("Error writing:", err)
		return []byte{}
	}

	buffer := make([]byte, 1024)
	n, err := t.conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading:", err)
		return []byte{}
	}
	apeape := buffer[:n]
	// fmt.Printf("Received from server: %s\n", string(apeape))
	fmt.Printf("Received from server: %x\n", apeape)
	return apeape
}

func mainpp() {
	log.Println("Starting")
	c, _ := torrent.NewClient(nil)
	defer c.Close()
	log.Println("reading torrent file")
	tt, _ := c.AddTorrentFromFile("zzz.torrent")
	// <-tt.GotInfo()
	// if err != nil {
	// 	log.Println(err)
	// 	panic("unable to read torrent from file")
	// }
	log.Printf("Infohash: %+v", tt.InfoHash())
	tt.DownloadAll()
	c.WaitAll()
}
