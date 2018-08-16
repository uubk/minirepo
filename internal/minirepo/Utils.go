/*
 * Copyright 2018 The minirepo authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package minirepo

import (
	"crypto"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"golang.org/x/crypto/openpgp/s2k"
	"os"
	"time"
)

// loadKeyFromFile loads an GPG ASCII-armored RSA key from 'file'. 'pubkey' decideds whether this is the public or the private part
func loadKeyFromFile(file string, pubkey bool) packet.Packet {
	pubkeyFD, err := os.Open(file)
	if err != nil {
		log.WithFields(log.Fields{
			"file": file,
		}).WithError(err).Fatal("Couldn't open key file")
	}
	res, err := armor.Decode(pubkeyFD)
	if err != nil {
		log.WithFields(log.Fields{
			"file": file,
		}).WithError(err).Fatal("Couldn't read key file")
	}

	if pubkey {
		if res.Type != openpgp.PublicKeyType {
			log.WithFields(log.Fields{
				"file": file,
				"type": res.Type,
			}).Fatal("Public key file didn't contain public key")
		}
		pkReader := packet.NewReader(res.Body)
		pkPacket, err := pkReader.Next()
		if res.Type != openpgp.PublicKeyType {
			log.WithFields(log.Fields{
				"file": file,
				"type": res.Type,
			}).WithError(err).Fatal("Couldn't read public key packet")
		}
		_, ok := pkPacket.(*packet.PublicKey)
		if !ok {
			log.WithFields(log.Fields{
				"file": file,
			}).Fatal("Public key is not a public key")
		}
		return pkPacket
	}
	if res.Type != openpgp.PrivateKeyType {
		log.WithFields(log.Fields{
			"file": file,
			"type": res.Type,
		}).Fatal("Private key file didn't contain private key")
	}
	pkReader := packet.NewReader(res.Body)
	pkPacket, err := pkReader.Next()
	if res.Type != openpgp.PrivateKeyType {
		log.WithFields(log.Fields{
			"file": file,
			"type": res.Type,
		}).WithError(err).Fatal("Couldn't read private key packet")
	}
	_, ok := pkPacket.(*packet.PrivateKey)
	if !ok {
		log.WithFields(log.Fields{
			"file": file,
		}).Fatal("Private key is not a private key")
	}
	return pkPacket
}

// fakeEntity creates a fake opengpg.Entity from a public and private key that can be used to sign files
func fakeEntity(pubkey *packet.PublicKey, privkey *packet.PrivateKey) *openpgp.Entity {
	cfg := packet.Config{
		DefaultCipher:          packet.CipherAES128,
		DefaultHash:            crypto.SHA384,
		DefaultCompressionAlgo: packet.CompressionNone,
		RSABits:                2048,
	}

	entity := openpgp.Entity{
		PrimaryKey: pubkey,
		PrivateKey: privkey,
		Identities: make(map[string]*openpgp.Identity),
	}

	// Create fake identity
	// We need an identity, not because it is used or important in any way but because Sign(...) will use the
	// same hash function (hardcoded) or fall back to RIPEMD...
	trueVal := true
	hid, _ := s2k.HashToHashId(crypto.SHA256)
	entity.Identities[""] = &openpgp.Identity{
		Name:   "",
		UserId: packet.NewUserId("", "", ""),
		SelfSignature: &packet.Signature{
			CreationTime:  time.Now(),
			SigType:       packet.SigTypePositiveCert,
			PubKeyAlgo:    packet.PubKeyAlgoRSA,
			Hash:          cfg.Hash(),
			IsPrimaryId:   &trueVal,
			FlagSign:      true,
			FlagsValid:    true,
			FlagCertify:   true,
			IssuerKeyId:   &pubkey.KeyId,
			PreferredHash: []uint8{hid},
		},
	}

	return &entity
}
