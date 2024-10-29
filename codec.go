package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Decode func(src []byte) string
type Encode func(src string) []byte

func ParseCodePage(cp string) encoding.Encoding {
	switch cp {
	case "GB2312":
		return simplifiedchinese.HZGB2312
	case "GBK":
		return simplifiedchinese.GBK
	case "C936":
		return simplifiedchinese.GBK
	case "SHJIS":
		return japanese.ShiftJIS
	case "C932":
		return japanese.ShiftJIS
	case "UTF8":
		return unicode.UTF8
	default:
		return unicode.UTF8
	}
}

// create decoder
func Decoder(enc encoding.Encoding) Decode {
	codec := enc.NewDecoder()
	return func(src []byte) string {
		var reader io.Reader = bytes.NewReader(src)
		reader = codec.Reader(reader)
		out, err := io.ReadAll(reader)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		return string(out)
	}
}

// create encoder
func Encoder(enc encoding.Encoding) Encode {
	codec := enc.NewEncoder()
	return func(src string) []byte {
		var reader io.Reader = strings.NewReader(src)
		reader = transform.NewReader(reader, codec)
		out, err := io.ReadAll(reader)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		return out
	}
}

// C932 & C936 => unicode
func DecoderMix(enc encoding.Encoding) Decode {
	codec_jp := enc.NewDecoder()
	return func(src []byte) string {
		reader := bytes.NewReader(src)
		output := strings.Builder{}
		for {
			head, err := reader.ReadByte()
			if err == io.EOF {
				break
			}
			tail, err := reader.ReadByte()
			if err == io.EOF {
				output.WriteByte(head)
				// fmt.Printf("%X ", head)
				break
			}
			code := []byte{head, tail}
			if isShiftJIS(head, tail) {
				rst, err := codec_jp.Bytes(code)
				if err == nil {
					output.Write(rst)
					continue
				}
			}
			output.Write(code)
			// fmt.Printf("%X %X ", head, tail)
		}
		return output.String()
	}
}

// unicode => C932 & C936
func EncoderMix(enc encoding.Encoding) Encode {
	codec_jp := enc.NewEncoder()
	return func(src string) []byte {
		buffer := []byte(src)
		output := bytes.Buffer{}
		for {
			rst, num, err := transform.Bytes(codec_jp, buffer)
			if err == io.EOF {
				break
			} else if err == nil {
				output.Write(rst)
				break
			}
			if num > 0 {
				output.Write(rst)
				buffer = buffer[num:]
				// continue
			}
			// _, size := utf8.DecodeRune(buffer)
			size := 1
			output.Write(buffer[0:size])
			// fmt.Printf("%X ", buffer[0])
			buffer = buffer[size:]
		}
		return output.Bytes()
	}
}

func isShiftJIS(head uint8, tail uint8) bool {
	if head < 0x81 || 0xEE < head {
		return false
	}
	if tail < 0x40 || 0xFC < tail {
		return false
	}
	if 0x9F <= head && head < 0xE0 {
		return false
	}
	if 0x85 <= head && head < 0x87 {
		return false
	}
	if 0xEB <= head && head < 0xED {
		return false
	}
	return tail != 0x7F
}
