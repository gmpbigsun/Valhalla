package handlers

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"

	"github.com/Hucaru/Valhalla/common/connection"
	"github.com/Hucaru/Valhalla/common/constants"
	"github.com/Hucaru/Valhalla/common/crypt"
	"github.com/Hucaru/Valhalla/common/packet"
)

// HandlePacket -
func HandlePacket(conn connection.Connection, buffer packet.Packet, isHeader bool) int {
	size := constants.CLIENT_HEADER_SIZE

	if isHeader {
		// Reading encrypted header
		size = crypt.GetPacketLength(buffer)
	} else {
		// Handle data packet
		pos := 0

		opcode := buffer.ReadByte(&pos)

		switch opcode {
		case constants.LOGIN_OP:
			handleLoginRequest(buffer, &pos, conn)
		}
	}

	return size
}

func handleLoginRequest(p packet.Packet, pos *int, conn connection.Connection) {
	fmt.Println("Login packet received")
	usernameLength := p.ReadShort(pos)
	username := p.ReadString(pos, usernameLength)

	passwordLength := p.ReadShort(pos)
	password := p.ReadString(pos, passwordLength)

	// hash and salt the password#
	hasher := sha512.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	var userID uint32
	var user string
	var databasePassword string
	var isLogedIn bool
	var isBanned int
	var isAdmin bool

	err := connection.Db.QueryRow("SELECT userID, username, password, isLogedIn, isBanned, isAdmin FROM users WHERE username=?", username).
		Scan(&userID, &user, &databasePassword, &isLogedIn, &isBanned, &isAdmin)

	result := byte(0x00)

	if err != nil {
		result = 0x05
	} else if hashedPassword != databasePassword {
		result = 0x04
	} else if isLogedIn {
		result = 0x07
	} else if isBanned > 0 {
		result = 0x02
	}

	// -Banned- = 2
	// Deleted or Blocked = 3
	// Invalid Password = 4
	// Not Registered = 5
	// Sys Error = 6
	// Already online = 7
	// System error = 9
	// Too many requests = 10
	// Older than 20 = 11
	// Master cannot login on this IP = 13

	pac := packet.NewPacket()
	pac.WriteByte(0x01)
	pac.WriteByte(result)
	pac.WriteByte(0x00)
	pac.WriteInt(0)

	if result <= 0x01 {
		pac.WriteInt(userID)
		pac.WriteByte(0x00)
		pac.WriteByte(0x01)
		pac.WriteString(username)
	} else if result == 0x02 {
		// Being banned is not yet implemented
	}

	pac.WriteLong(0)
	pac.WriteLong(0)
	pac.WriteLong(0)
	conn.Write(pac)

	if result > 0x01 {
		return
	}

	pac = packet.NewPacket()
	pac.WriteByte(0x03)
	pac.WriteByte(0x04)
	pac.WriteByte(0x00)
	//conn.Write(pac)

	pac = packet.NewPacket()
	pac.WriteByte(constants.PLAYER_REQUEST_WORLD_LIST)
	pac.WriteString("hash")
	//conn.Write(pac)
}