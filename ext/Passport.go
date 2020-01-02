package ext

import (
	"encoding/json"
	"net/url"
)

type PassportData struct {
	Data        []EncryptedPassportElement `json:"data"`
	Credentials EncryptedCredentials       `json:"credentials"`
}

type PassportFile struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	FileDate     int    `json:"file_date"`
}

type EncryptedPassportElement struct {
	Type        string         `json:"type"`
	Data        string         `json:"data"`
	PhoneNumber string         `json:"phone_number"`
	Email       string         `json:"email"`
	Files       []PassportFile `json:"files"`
	FrontSide   PassportFile   `json:"front_side"`
	ReverseSide PassportFile   `json:"reverse_side"`
	Selfie      PassportFile   `json:"selfie"`
	Translation []PassportFile `json:"translation"`
	Hash        string         `json:"hash"`
}

type EncryptedCredentials struct {
	Data   string `json:"data"`
	Hash   string `json:"hash"`
	Secret string `json:"secret"`
}

type PassportElementError interface {
	Data() (url.Values, error) // return interface as url values
}

func (b Bot) SetPassportDataErrors(userId int, errors ...PassportElementError) {

}

type BasePassportError struct {
	Source  string
	Type    string
	Message string
}

type PassportElementErrorDataField struct {
	BasePassportError
	FieldName string
	DataHash  string
}

func (pe PassportElementErrorDataField) Data() (url.Values, error) {
	v := pe.BasePassportError.values()
	v.Add("field_name", pe.FieldName)
	v.Add("data_hash", pe.DataHash)
	return v, nil
}

type PassportElementErrorFrontSide struct {
	BasePassportError
	FileHash string
}

func (pe PassportElementErrorFrontSide) Data() (url.Values, error) {
	v := pe.BasePassportError.values()
	v.Add("file_hash", pe.FileHash)
	return v, nil
}

type PassportElementErrorReverseSide struct {
	BasePassportError
	FileHash string
}

func (pe PassportElementErrorReverseSide) Data() (url.Values, error) {
	v := pe.BasePassportError.values()
	v.Add("file_hash", pe.FileHash)
	return v, nil
}

type PassportElementErrorSelfie struct {
	BasePassportError
	FileHash string
}

func (pe PassportElementErrorSelfie) Data() (url.Values, error) {
	v := pe.BasePassportError.values()
	v.Add("file_hash", pe.FileHash)
	return v, nil
}

type PassportElementErrorFile struct {
	BasePassportError
	FileHash string
}

func (pe PassportElementErrorFile) Data() (url.Values, error) {
	v := pe.BasePassportError.values()
	v.Add("file_hash", pe.FileHash)
	return v, nil
}

type PassportElementErrorFiles struct {
	BasePassportError
	FileHashes []string
}

func (pe PassportElementErrorFiles) Data() (url.Values, error) {
	v := pe.BasePassportError.values()
	bs, err := json.Marshal(pe.FileHashes)
	if err != nil {
		return nil, err
	}
	v.Add("file_hashes", string(bs))
	return v, nil
}

type PassportElementErrorTranslationFile struct {
	BasePassportError
	FileHash string
}

func (pe PassportElementErrorTranslationFile) Data() (url.Values, error) {
	v := pe.BasePassportError.values()
	v.Add("file_hash", pe.FileHash)
	return v, nil
}

type PassportElementErrorTranslationFiles struct {
	BasePassportError
	FileHashes []string
}

func (pe PassportElementErrorTranslationFiles) Data() (url.Values, error) {
	v := pe.BasePassportError.values()
	bs, err := json.Marshal(pe.FileHashes)
	if err != nil {
		return nil, err
	}
	v.Add("file_hashes", string(bs))
	return v, nil
}

type PassportElementErrorUnspecified struct {
	BasePassportError
	ElementHash string
}

func (pe PassportElementErrorUnspecified) Data() (url.Values, error) {
	v := pe.BasePassportError.values()
	v.Add("element_hash", pe.ElementHash)
	return v, nil
}

func (b BasePassportError) values() (v url.Values) {
	v.Add("source", b.Source)
	v.Add("type", b.Type)
	v.Add("message", b.Message)
	return v
}
