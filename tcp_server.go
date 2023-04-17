package main

import (
	"io"
	"log"
	"net"
	"sync"
	"strings"
	"strconv"
)

var ALF_SEM03 []rune = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.,:; ")

func CelsiusToFahrenheit(value float64) float64 {
	var fahrenheit float64
	fahrenheit = (value * 1.8) + 32
	return fahrenheit
}

func checkKjevik(dekryptertMelding []rune) []rune {
    checkKjevik := string(dekryptertMelding)
    if strings.HasPrefix(checkKjevik, "Kjevik") {
        runeKjevik := []rune(checkKjevik)
        lastRune := runeKjevik[len(runeKjevik)-1]
        lastRuneValue, err := strconv.Atoi(string(lastRune))
        if err != nil {
            return dekryptertMelding
        }
        newRuneValue := CelsiusToFahrenheit(float64(lastRuneValue))
        newRuneString := strconv.FormatFloat(newRuneValue, 'f', -1, 64)
        newRunes := []rune(newRuneString)
        runeKjevik = append(runeKjevik[:len(runeKjevik)-1], newRunes...)
        kryptertMelding2 := Krypter(runeKjevik, ALF_SEM03, 4)
        return kryptertMelding2
    }
    return dekryptertMelding
}

func Krypter(melding []rune, alphabet []rune, chiffer int) []rune {
        kryptertMelding := make([]rune, len(melding))
        for i := 0; i < len(melding); i++ {
                indeks := sokIAlfabetet(melding[i], alphabet)
                newIndex := (indeks + chiffer) % len(alphabet)
                kryptertMelding[i] = alphabet[newIndex]
        }
        return kryptertMelding
}

func sokIAlfabetet(symbol rune, alfabet []rune) int {
	for i := 0; i < len(alfabet); i++ {
			if symbol == alfabet[i] {
					return i
					break
			}
	}
	return -1
}

func main() {

	var wg sync.WaitGroup

	server, err := net.Listen("tcp", "172.17.0.2:")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("bundet til %s", server.Addr().String())
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			log.Println("før server.Accept() kallet")
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				for {
					buf := make([]byte, 1024)
					n, err := c.Read(buf)
					if err != nil {
						if err != io.EOF {
							log.Println(err)
						}
						return // fra for løkke
					}
					dekryptertMelding := Krypter([]rune(string(buf[:n])), ALF_SEM03, len(ALF_SEM03)-4)
					log.Printf("dekrypt: %s", string(dekryptertMelding))
					kryptertMelding := Krypter(([]rune("pong")), ALF_SEM03, 4)
					log.Printf("mld1: %s", string(kryptertMelding))
					kryptertMelding2 := checkKjevik(dekryptertMelding)
					log.Printf("mld2: %s", string(kryptertMelding2))	
                                        switch msg := string(dekryptertMelding); msg {
                                        case "ping":
                                                _, err = c.Write([]byte(string(kryptertMelding)))
					case "Kjevik;SN39040;18.03.2022 01:50;6":
                                                _, err = c.Write([]byte(string(kryptertMelding2)))
                                        default:
                                                _, err = c.Write([]byte(string(dekryptertMelding)))
                                        }
                                        if err != nil {
                                                if err != io.EOF {
                                                        log.Println(err)
                                                }
                                                return // fra for l  kke
                                        }
                                }
                        }(conn)
                }
        }()
        wg.Wait()
}
