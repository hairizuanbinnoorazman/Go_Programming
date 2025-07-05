package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	badger "github.com/dgraph-io/badger/v4"
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

type SingleTorrentState struct {
	Blocks           []bool `json:"blocks"`
	DownloadingBlock int    `json:"downloading_block" `
}

func (s *SingleTorrentState) fromBytes(val []byte) {
	var sts SingleTorrentState
	json.Unmarshal(val, &sts)
	s.Blocks = sts.Blocks
	s.DownloadingBlock = sts.DownloadingBlock
}

func (s *SingleTorrentState) toBytes() []byte {
	val, _ := json.Marshal(s)
	return val
}

func nextBlock(currentBlocks, otherPeerBlocks []bool) (nextBlock int, completed bool) {
	// Convert currentBlocks with false to int slice
	allFalseCurrentBlock := []int{}
	for k, v := range currentBlocks {
		if v == false {
			allFalseCurrentBlock = append(allFalseCurrentBlock, k)
		}
	}
	if len(allFalseCurrentBlock) == 0 {
		return -1, true
	}
	// Convert otherPeerBlocks with true to int slice
	allTrueOtherPeerBlock := []int{}
	for k, v := range otherPeerBlocks {
		if v {
			allTrueOtherPeerBlock = append(allTrueOtherPeerBlock, k)
		}
	}
	// See matching int slice values
	possibleBlocks := []int{}
	for _, v := range allFalseCurrentBlock {
		for _, m := range allTrueOtherPeerBlock {
			if v == m {
				possibleBlocks = append(possibleBlocks, v)
			}
		}
	}
	fmt.Println(possibleBlocks)
	if len(possibleBlocks) == 0 {
		return -1, false
	}
	return possibleBlocks[rand.Intn(len(possibleBlocks))], false
}

type TorrentProgresState struct {
	db *badger.DB
}

func (t *TorrentProgresState) ViewPeerState(infohash, peerID []byte) SingleTorrentState {
	key := append(infohash, peerID...)
	var sts SingleTorrentState
	t.db.View(
		func(txn *badger.Txn) error {
			val, err := txn.Get(key)
			rawVal := make([]byte, val.ValueSize())
			val.ValueCopy(rawVal)
			sts.fromBytes(rawVal)
			fmt.Printf("Infohash: %x :: PeerID: %x State: %+v\n", infohash, peerID, sts)
			return err
		})
	return sts
}

func (t *TorrentProgresState) AddPeer(infohash, peerID []byte, blocks []bool) {
	sts := SingleTorrentState{
		Blocks:           blocks,
		DownloadingBlock: -1,
	}
	raw, _ := json.Marshal(sts)
	key := append(infohash, peerID...)
	t.db.Update(
		func(txn *badger.Txn) error {
			err := txn.Set(key, raw)
			err = txn.Commit()
			return err
		})
}

func (t *TorrentProgresState) StartDownload(infohash, currentPeerID, peerID []byte) int {
	key := append(infohash, peerID...)
	currentPeerKey := append(infohash, currentPeerID...)
	var toDownload int
	t.db.Update(
		func(txn *badger.Txn) error {
			val, _ := txn.Get(key)
			rawVal := make([]byte, val.ValueSize())
			val.ValueCopy(rawVal)
			var sts SingleTorrentState
			sts.fromBytes(rawVal)

			val, _ = txn.Get(currentPeerKey)
			rawVal = make([]byte, val.ValueSize())
			val.ValueCopy(rawVal)
			var currentSts SingleTorrentState
			currentSts.fromBytes(rawVal)

			fmt.Println(sts)
			fmt.Println(currentSts)
			toDownload, _ = nextBlock(currentSts.Blocks, sts.Blocks)

			sts.DownloadingBlock = toDownload
			raw, _ := json.Marshal(sts)
			err := txn.Set(key, raw)
			txn.Commit()
			return err
		})
	return toDownload
}

func (t *TorrentProgresState) StopDownload(infohash, peerID []byte) {
	key := append(infohash, peerID...)
	t.db.Update(
		func(txn *badger.Txn) error {
			val, _ := txn.Get(key)
			rawVal := make([]byte, val.ValueSize())
			val.ValueCopy(rawVal)
			var sts SingleTorrentState
			sts.fromBytes(rawVal)

			sts.DownloadingBlock = -1
			raw, _ := json.Marshal(sts)
			err := txn.Set(key, raw)
			txn.Commit()
			return err
		})
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

func (tf *TorrentFile) ExecuteFileCheck() []bool {
	aa := []bool{}
	for idx, _ := range tf.Pieces {
		a := tf.ExecutePieceCheck(idx)
		aa = append(aa, a)
	}
	return aa
}

func (tf *TorrentFile) ExecutePieceCheck(idx int) bool {
	ff, err := os.Open(tf.Filename)
	if err != nil {
		fmt.Println("Error - unable to read file")
		return false
	}
	defer ff.Close()

	rawPieceContent := make([]byte, tf.PieceLength)
	if (idx+1)*tf.PieceLength > tf.Length {
		extraBytes := (idx+1)*tf.PieceLength - tf.Length
		fmt.Printf("Extra bytes: %v\n", extraBytes)
		rawPieceContent = make([]byte, tf.PieceLength-(extraBytes))
	}

	ff.ReadAt(rawPieceContent, int64(idx*tf.PieceLength))
	ss2 := sha1.New()
	ss2.Write(rawPieceContent)
	oo2 := ss2.Sum(nil)
	contentCheck := bytes.Equal(oo2, tf.Pieces[idx])
	tf.PieceCheck[idx] = contentCheck
	return contentCheck
}

func (tf *TorrentFile) GetContent(idx, offset, blockLength int) []byte {
	ff, _ := os.Open(tf.Filename)
	defer ff.Close()

	rawBlockContent := make([]byte, blockLength)
	ff.ReadAt(rawBlockContent, int64(idx*tf.PieceLength+offset))
	return rawBlockContent
}

func generatePeerID() []byte {
	initialVal := rand.Int31()
	ss := sha1.New()
	rawInitialVal := make([]byte, 4)
	binary.BigEndian.PutUint32(rawInitialVal, uint32(initialVal))
	ss.Write(rawInitialVal)
	peerID := ss.Sum(nil)
	return peerID
}

func main() {
	tcpListener, _ := os.LookupEnv("BT_PORT")
	httpListener, _ := os.LookupEnv("HTTP_PORT")
	peerURL, _ := os.LookupEnv("PEER_URL")
	folderPath, _ := os.LookupEnv("FOLDER_PATH")

	if !strings.HasSuffix(folderPath, "/") && folderPath != "" {
		folderPath += "/"
	}

	opt := badger.DefaultOptions("").WithInMemory(true)
	db, err := badger.Open(opt)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	aa := TorrentProgresState{
		db: db,
	}
	// infohash := []byte("zzz")
	// currentPeerID := []byte("p1")
	// otherPeerID := []byte("p2")
	// aa.AddPeer(infohash, currentPeerID, []bool{false, false, false})
	// aa.AddPeer(infohash, otherPeerID, []bool{true, true, true})
	// hoho := aa.StartDownload(infohash, currentPeerID, otherPeerID)
	// fmt.Println(hoho)
	// aa.ViewPeerState(infohash, currentPeerID)
	// aa.ViewPeerState(infohash, otherPeerID)
	// aa.StopDownload(infohash, otherPeerID)
	// aa.ViewPeerState(infohash, currentPeerID)
	// aa.ViewPeerState(infohash, otherPeerID)

	peerID := generatePeerID()
	slog.Info(fmt.Sprintf("PeerID: %x", peerID))

	pp, _ := hex.DecodeString("d69f91e6b2ae4c542468d1073a71d4ea13879a7f")
	tm := TorrentMessage{
		peerID:   peerID,
		infohash: pp,
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", tcpListener))
	if err != nil {
		panic("Unable to listen to port")
	}
	defer listener.Close()

	go handler(listener, tm, aa, peerID)

	http.Handle("/add-torrent", &addTorrentHandler{currentPeerID: peerID, tps: aa, folderPath: folderPath})
	http.Handle("/start", startHandler{PeerURL: peerURL, TM: tm, TPS: aa, FolderPath: folderPath, CurrentPeerID: peerID})

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

type startHandler struct {
	PeerURL       string
	CurrentPeerID []byte
	TM            TorrentMessage
	TPS           TorrentProgresState
	FolderPath    string
}

func (f startHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("start foo handler")
	defer slog.Info("end foo handler")

	go f.BeginTorrent()
}

func (f startHandler) BeginTorrent() {
	conn, err := net.Dial("tcp", f.PeerURL)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	yahoo, _ := os.ReadFile("zzz3.torrent")
	var zz zzz
	bencode.Unmarshal(yahoo, &zz)
	// TODO: Allow modify with Folder path
	tf := NewTorrentFile("hoho.mkv", zz.Info.Length, zz.Info.PieceLength, zz.Info.Pieces)

	val, _ := bencode.Marshal(zz.Info)
	ss := sha1.New()
	ss.Write(val)
	infohash := ss.Sum(nil)

	f.TM.infohash = infohash

	conn.Write(f.TM.Handshake())
	// Skip the "handshake section"
	// First 20 bytes is protocol info
	// Next 8 bytes is extension
	// Next 20 bytes is infohash
	// Next 20 bytes is peerid
	// Next 4 bytes Length of message
	// Next 1 byte is message id
	tmp := make([]byte, 28)
	conn.Read(tmp)

	rawInfohash := make([]byte, 20)
	conn.Read(rawInfohash)
	fmt.Printf("Infohash: %x\n", rawInfohash)

	rawPeerID := make([]byte, 20)
	conn.Read(rawPeerID)
	fmt.Printf("Peer ID: %x\n", rawPeerID)

	pieceID := -1
	// rawContent := []byte{}
	blockOffset := 0
	// initialPieceState := slices.Repeat([]bool{false}, 377)

	for {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))

		rawMessageLength := make([]byte, 4)
		_, err := conn.Read(rawMessageLength)
		if err != nil {
			fmt.Println("connection probably failed. let's retry again next time")
			return
		}
		messageLength := int(binary.BigEndian.Uint32(rawMessageLength))
		fmt.Printf("%d - %x\n", messageLength, rawMessageLength)

		rawMessageType := make([]byte, 1)
		_, err = conn.Read(rawMessageType)
		if err != nil {
			fmt.Println("connection probably failed. let's retry again next time")
			return
		}
		messageType := int(rawMessageType[0])

		switch messageType {
		case 5:
			fmt.Println("Bitfield type")
			bitfieldMessage := make([]byte, messageLength-1)
			conn.Read(bitfieldMessage)
			f.TPS.AddPeer(rawInfohash, rawPeerID, f.TM.allReverseBitField(bitfieldMessage))
		case 7:
			fmt.Println("Piece type")

			rawPieceIndex := make([]byte, 4)
			conn.Read(rawPieceIndex)
			pieceIndex := int(binary.BigEndian.Uint32(rawPieceIndex))

			rawOffset := make([]byte, 4)
			conn.Read(rawOffset)
			offset := int(binary.BigEndian.Uint32(rawOffset))

			fmt.Printf("PieceIndex: %v :: Offset: %v\n", pieceIndex, offset)

			rawContent := make([]byte, messageLength-1-4-4)
			aa, _ := conn.Read(rawContent)
			fmt.Printf("Obtained %v bytes from other peer", aa)
		default:
			fmt.Println("message type undertermined")
			time.Sleep(10 * time.Second)
			continue
		}

		if pieceID == -1 {
			pieceID = f.TPS.StartDownload(infohash, f.CurrentPeerID, rawPeerID)
		}

		if blockOffset == tf.PieceLength {
			pieceID += f.TPS.StartDownload(infohash, f.CurrentPeerID, rawPeerID)
			blockOffset = 0
		}

		blockLength := 16384
		if pieceID*tf.PieceLength+blockOffset+blockLength > tf.Length {
			blockLength = blockLength - (tf.Length - pieceID*tf.PieceLength + blockOffset + blockLength)
		}

		conn.Write(f.TM.Request(pieceID, blockOffset, blockLength))
		blockOffset += 16384
	}
}

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

type addTorrentHandler struct {
	currentPeerID []byte
	tps           TorrentProgresState
	folderPath    string
}

func (f *addTorrentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("add torrent handler")
	defer slog.Info("end torrent handler")

	yahoo, _ := os.ReadFile("zzz3.torrent")

	var zz zzz
	oo := bencode.Unmarshal(yahoo, &zz)
	fmt.Println(oo)

	val, _ := bencode.Marshal(zz.Info)
	ss := sha1.New()
	ss.Write(val)
	infohash := ss.Sum(nil)

	filePath := f.folderPath + zz.Info.Name

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		aa, err := os.Create(filePath)
		if err != nil {
			fmt.Println(err)
		}
		emptyContent := make([]byte, zz.Info.Length)
		aa.Write(emptyContent)
		aa.Sync()
	}

	tf := NewTorrentFile(filePath, zz.Info.Length, zz.Info.PieceLength, zz.Info.Pieces)
	vals := tf.ExecuteFileCheck()
	f.tps.AddPeer(infohash, f.currentPeerID, vals)
	fmt.Printf("Torrent Added. Infohash: %x. Pieces detected: %v\n", infohash, len(vals))
}

func handler(listener net.Listener, tm TorrentMessage, tps TorrentProgresState, currentPeerID []byte) {
	slog.Info("Start listener function")
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Handle error (e.g., connection closed)
			continue
		}
		go handleConnection(conn, tm, tps, currentPeerID) // Handle connection in a new goroutine
	}
}

func handleConnection(conn net.Conn, tm TorrentMessage, tps TorrentProgresState, currentPeerID []byte) {
	slog.Info("Wait for the handshake first. We are receiver section of the codebase")

	// Skip the "handshake section"
	// First 20 bytes is protocol info
	// Next 8 bytes is extension
	// Next 20 bytes is infohash
	// Next 20 bytes is peerid
	tmp := make([]byte, 28)
	conn.Read(tmp)
	fmt.Println(string(tmp))

	rawInfohash := make([]byte, 20)
	conn.Read(rawInfohash)
	fmt.Printf("Infohash: %x\n", rawInfohash)

	rawPeerID := make([]byte, 20)
	conn.Read(rawPeerID)
	fmt.Printf("PeerID: %x\n", rawPeerID)

	yahoo, _ := os.ReadFile("zzz3.torrent")
	var zz zzz
	bencode.Unmarshal(yahoo, &zz)
	// TODO: Allow modify with Folder path
	tf := NewTorrentFile("hoho.mkv", zz.Info.Length, zz.Info.PieceLength, zz.Info.Pieces)

	val, _ := bencode.Marshal(zz.Info)
	ss := sha1.New()
	ss.Write(val)
	infohash := ss.Sum(nil)

	tm.infohash = infohash

	// Immediately return handshake
	prefix := tm.Handshake()
	state := tps.ViewPeerState(rawInfohash, currentPeerID)
	bitfield := tm.Bitfield(tm.allConvertBitField(state.Blocks))
	fullMessage := append(prefix, bitfield...)
	conn.Write(fullMessage)

	for {
		conn.SetReadDeadline(time.Now().Add(120 * time.Second))
		time.Sleep(5 * time.Second)

		rawMessageLength := make([]byte, 4)
		_, err := conn.Read(rawMessageLength)
		if err != nil {
			fmt.Println("connection probably failed. let's retry again next time")
			return
		}
		messageLength := int(binary.BigEndian.Uint32(rawMessageLength))
		fmt.Printf("Message Length: %d\n", messageLength)

		rawMessageType := make([]byte, 1)
		_, err = conn.Read(rawMessageType)
		if err != nil {
			fmt.Println("connection probably failed. let's retry again next time")
			return
		}
		messageType := int(rawMessageType[0])

		switch messageType {
		case 6:
			fmt.Println("Request type")

			rawPieceIndex := make([]byte, 4)
			conn.Read(rawPieceIndex)
			pieceIndex := int(binary.BigEndian.Uint32(rawPieceIndex))

			rawOffset := make([]byte, 4)
			conn.Read(rawOffset)
			offset := int(binary.BigEndian.Uint32(rawOffset))

			rawLength := make([]byte, 4)
			conn.Read(rawLength)
			contentLength := int(binary.BigEndian.Uint32(rawLength))

			fmt.Printf("PieceIndex: %v :: Offset: %v :: ContentLength: %v\n", pieceIndex, offset, contentLength)
			output := tf.GetContent(pieceIndex, offset, contentLength)

			response := tm.Piece(pieceIndex, offset, output)
			conn.Write(response)
		}
	}
}

func mainzz() {
	pp, _ := hex.DecodeString("13426974546f7272656e742070726f746f636f6c0000000000000004d69f91e6b2ae4c542468d1073a71d4ea13879a7f2d524e302e302e302d9d5e0226875575c6233a130000000205e0")
	tt := NewTorrentMessage(pp)
	fmt.Printf("%v", tt)

	// zz := tt.Request(0, 0)
	// fmt.Println(hex.EncodeToString(zz))

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

func (t TorrentMessage) Bitfield(bitfieldContent []byte) []byte {
	messageType := []byte{5}
	message := append(messageType, bitfieldContent...)
	messageLength := t.convertToBytes(len(message))
	output := append(messageLength, message...)
	return output
}

func (t TorrentMessage) Request(indexPiece, blockOffset, rawBlockLength int) []byte {
	message := []byte{6}
	rawIndexPiece := t.convertToBytes(indexPiece)
	rawBlockOffset := t.convertToBytes(blockOffset)
	// TODO: For last piece, last block, it will be different
	// blockLength := t.convertToBytes(16384)
	blockLength := t.convertToBytes(rawBlockLength)
	message = append(message, rawIndexPiece...)
	message = append(message, rawBlockOffset...)
	message = append(message, blockLength...)

	messageLength := t.convertToBytes(len(message))
	output := append(messageLength, message...)
	return output
}

func (t TorrentMessage) Piece(indexPiece, blockOffset int, content []byte) []byte {
	message := []byte{7}
	rawIndexPiece := t.convertToBytes(indexPiece)
	rawBlockOffset := t.convertToBytes(blockOffset)
	// TODO: For last piece, last block, it will be different
	message = append(message, rawIndexPiece...)
	message = append(message, rawBlockOffset...)
	message = append(message, content...)

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
	tt, _ := c.AddTorrentFromFile("zzz3.torrent")
	// <-tt.GotInfo()
	// if err != nil {
	// 	log.Println(err)
	// 	panic("unable to read torrent from file")
	// }
	log.Printf("Infohash: %+v", tt.InfoHash())
	tt.DownloadAll()
	c.WaitAll()
}
