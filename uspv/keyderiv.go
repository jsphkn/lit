package uspv

import (
	"log"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
	"github.com/mit-dci/lit/portxo"
)

/*
Key derivation for a TxStore has 3 levels: use case, peer index, and keyindex.
Regular wallet addresses are use 0, peer 0, and then a linear index.
The identity key is use 11, peer 0, index 0.
Channel multisig keys are use 2, peer and index per peer and channel.
Channel refund keys are use 3, peer and index per peer / channel.
*/

// =====================================================================
// OK only use these now

// PathPrivkey returns a private key by descending the given path
// Returns nil if there's an error.
func (t *TxStore) PathPrivkey(kg portxo.KeyGen) *btcec.PrivateKey {
	// in uspv, we require path depth of 5
	if kg.Depth != 5 {
		return nil
	}
	priv, err := kg.DerivePrivateKey(t.rootPrivKey)
	if err != nil {
		log.Printf("PathPrivkey err %s", err.Error())
		return nil
	}
	return priv
}

// PathPubkey returns a public key by descending the given path.
// Returns nil if there's an error.
func (t *TxStore) PathPubkey(kg portxo.KeyGen) *btcec.PublicKey {
	priv := t.PathPrivkey(kg)
	if priv == nil {
		return nil
	}
	return t.PathPrivkey(kg).PubKey()
}

// PathPubHash160 returns a 20 byte pubkey hash for the given path
// It'll always return 20 bytes, or a nil if there's an error.
func (t *TxStore) PathPubHash160(kg portxo.KeyGen) []byte {
	pub := t.PathPubkey(kg)
	if pub == nil {
		return nil
	}
	return btcutil.Hash160(pub.SerializeCompressed())
}

// ------------- end of 2 main key deriv functions

// get a private key from the regular wallet
func (t *TxStore) GetWalletPrivkey(idx uint32) *btcec.PrivateKey {
	var kg portxo.KeyGen
	kg.Depth = 5
	kg.Step[0] = 44 | 1<<31
	kg.Step[1] = 0 | 1<<31
	kg.Step[2] = 0 | 1<<31
	kg.Step[3] = 0 | 1<<31
	kg.Step[4] = idx | 1<<31
	return t.PathPrivkey(kg)
}

// GetWalletKeygen returns the keygen for a standard wallet address
func GetWalletKeygen(idx uint32) portxo.KeyGen {
	var kg portxo.KeyGen
	kg.Depth = 5
	kg.Step[0] = 44 | 1<<31
	kg.Step[1] = 0 | 1<<31
	kg.Step[2] = 0 | 1<<31
	kg.Step[3] = 0 | 1<<31
	kg.Step[4] = idx | 1<<31
	return kg
}

// get a public key from the regular wallet
func (t *TxStore) GetWalletAddress(idx uint32) *btcutil.AddressWitnessPubKeyHash {
	if t == nil {
		log.Printf("GetAddress %d nil txstore\n", idx)
		return nil
	}
	priv := t.GetWalletPrivkey(idx)
	if priv == nil {
		log.Printf("GetAddress %d made nil pub\n", idx)
		return nil
	}
	adr, err := btcutil.NewAddressWitnessPubKeyHash(
		btcutil.Hash160(priv.PubKey().SerializeCompressed()), t.Param)
	if err != nil {
		log.Printf("GetAddress %d made nil pub\n", idx)
		return nil
	}
	return adr
}

// GetUsePrive generates a private key for the given use case & keypath
func (t *TxStore) GetUsePriv(kg portxo.KeyGen, use uint32) *btcec.PrivateKey {
	kg.Step[2] = use
	return t.PathPrivkey(kg)
}

// GetUsePub generates a pubkey for the given use case & keypath
func (t *TxStore) GetUsePub(kg portxo.KeyGen, use uint32) [33]byte {
	var b [33]byte
	pub := t.GetUsePriv(kg, use).PubKey()
	if pub != nil {
		copy(b[:], pub.SerializeCompressed())
	}
	return b
}

// IdKey returns the identity private key
func (t *TxStore) IdKeyx() *btcec.PrivateKey {
	var kg portxo.KeyGen
	kg.Depth = 5
	kg.Step[0] = 44 | 1<<31
	kg.Step[1] = 0 | 1<<31
	kg.Step[2] = 9 | 1<<31
	kg.Step[3] = 0 | 1<<31
	kg.Step[4] = 0 | 1<<31
	return t.PathPrivkey(kg)
}
