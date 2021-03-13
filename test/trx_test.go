package tron

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/fbsobreira/gotron-sdk/pkg/common"
	"github.com/fbsobreira/gotron-sdk/pkg/proto/core"
	"github.com/gogo/protobuf/proto"
	"io/ioutil"
	"strings"
	"testing"
	"time"
)

func TestTrx(t *testing.T) {
	g := client.NewGrpcClient("13.124.62.58:50051")
	err := g.Start()
	if err != nil {
		panic(err)
	}
	filepath := "xiao"
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	account := ""
	amount := 0

	errlist := ""

	ary := strings.Split(string(content), "\n")
	for _, t := range ary {
		if len(t) == 0 {
			continue
		}
		tary := strings.Split(t, ",")
		err = sendtx(g, tary[0], tary[1], account, int64(amount))
		if err != nil {
			fmt.Println(tary[0] + " 错误: " + err.Error())
			errlist += tary[0] + "," + tary[1] + "\n"
			continue
		}
	}

	err = ioutil.WriteFile("errlist", []byte(errlist), 0644)
	if err != nil {
		panic(err)
	}
}

func sendtx(g *client.GrpcClient, from string, pwd string, to string, amount int64) (err error) {
	transferContract := new(core.TransferContract)
	transferContract.OwnerAddress, err = common.DecodeCheck(from)
	if err != nil {
		return err
	}
	transferContract.ToAddress, err = common.DecodeCheck(to)
	if err != nil {
		return err
	}
	transferContract.Amount = amount

	GrpcTimeout := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), GrpcTimeout)
	defer cancel()

	coreTx, err := g.Client.CreateTransaction(ctx, transferContract)
	if err != nil {
		return err
	}

	privateKeyBytes, err := hex.DecodeString(pwd)
	sk, _ := btcec.PrivKeyFromBytes(btcec.S256(), privateKeyBytes)
	key := sk.ToECDSA()

	err = SignTransaction(coreTx, key)
	if err != nil {
		return err
	}

	broadRest, err := g.Client.BroadcastTransaction(ctx, coreTx)
	if err != nil {
		return err
	}
	if !broadRest.Result {
		return errors.New("失败")
	}
	return nil
}

func SignTransaction(transaction *core.Transaction, key *ecdsa.PrivateKey) error {
	if transaction.GetRawData() == nil {
		return errors.New("签名错误，raw为空")
	}
	transaction.GetRawData().Timestamp = time.Now().UnixNano() / 1000000

	rawData, err := proto.Marshal(transaction.GetRawData())

	if err != nil {
		return fmt.Errorf("sign transaction error: %v", err)
	}

	h256h := sha256.New()
	h256h.Write(rawData)
	hash := h256h.Sum(nil)

	//fmt.Printf("%x\n", hash)
	//fmt.Println("hash", hex.EncodeToString(hash))

	contractList := transaction.GetRawData().GetContract()

	for range contractList {
		signature, err := crypto.Sign(hash, key)

		if err != nil {
			return fmt.Errorf("sign transaction error: %v", err)
		}

		transaction.Signature = append(transaction.Signature, signature)
	}

	return nil
}
