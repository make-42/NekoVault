package main

import (
  "log"
  "github.com/gotk3/gotk3/gtk"
  "crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
  "crypto/sha256"
	"fmt"
  "io/ioutil"
	"io"
  "os/exec"
  "os"
  "path/filepath"
)

func check(e error) {
    if e != nil {
        log.Fatal(e);
        panic(e)
    }
}

func encryptfile(filenamepointer *string, s string, progressbar *gtk.ProgressBar){
  if (filenamepointer == nil){
    return
  }
  filename := *filenamepointer;
  fi, err := os.Stat(filename);
  var cmd *exec.Cmd;
  check(err);
  switch mode := fi.Mode(); {
  case mode.IsDir():
      cmd = exec.Command("zip","-j","-r","nekotemp.zip",filename);
  case mode.IsRegular():
      cmd = exec.Command("zip","-j","nekotemp.zip",filename);
  }
  cmd.Run();
  progressbar.SetFraction(0.5);
  dat, _:= ioutil.ReadFile("nekotemp.zip")

  h := sha256.New()
  h.Write([]byte(s))
  key := hex.EncodeToString(h.Sum(nil))

  encrypted := encrypt(dat, key);

  f, err := os.OpenFile(filepath.Dir(filename)+"/output.neko", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600);
  check(err);
  defer f.Close()
  f.WriteString(encrypted);
  os.Remove("nekotemp.zip");
}

func decryptfile(filenamepointer *string,s string,progressbar *gtk.ProgressBar){
  if (filenamepointer == nil){
    return
  }
  filename := *filenamepointer;
  progressbar.SetFraction(0.5);
  dat, _:= ioutil.ReadFile(filename);

  h := sha256.New()
  h.Write([]byte(s))
  key := hex.EncodeToString(h.Sum(nil))

	decrypted := decrypt(string(dat), key);
  f, err := os.OpenFile("nekotemp.zip", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600);
  check(err);
  defer f.Close()
  f.WriteString(decrypted);
  cmd := exec.Command("unzip","nekotemp.zip","-d",filepath.Dir(filename));
  cmd.Run();
  os.Remove("nekotemp.zip");
}


func encrypt(plaintext []byte, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	check(err);

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	check(err);

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

func decrypt(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	check(err);

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	check(err);

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	check(err);

	return fmt.Sprintf("%s", plaintext)
}
