// client.go
package gapBotApi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/amirimatin/gapBotApi/v2/params"
	"io"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"strconv"

	"github.com/amirimatin/gapBotApi/v2/models"
)

func (bot *BotAPI) buildParams(in Params) url.Values {
	out := url.Values{}
	for key, value := range in {
		out.Set(key, value)
	}
	return out
}

func (bot *BotAPI) MakeRequest(endpoint string, params Params) (*models.APIResponse, error) {
	if bot.Debug {
		fmt.Printf("Endpoint: %s, Params: %+v\n", endpoint, params)
	}

	method := fmt.Sprintf(bot.apiEndpoint, endpoint)
	values := bot.buildParams(params)

	var apiResp models.APIResponse
	_, err := bot.Client.R().
		SetBody(values.Encode()).
		SetResult(&apiResp).
		Post(method)

	if err != nil {
		return nil, err
	}
	if apiResp.Error != "" {
		return &apiResp, &models.Error{Message: apiResp.Error}
	}

	return &apiResp, nil
}

func (bot *BotAPI) Request(c Chattable) (*models.APIResponse, error) {
	params, err := c.params()
	if err != nil {
		return nil, err
	}

	if t, ok := c.(Fileable); ok {
		file := t.file()
		if file.Data.NeedsUpload() {
			uFile, err := bot.UploadFile(params, file)
			if err != nil {
				return nil, err
			}
			mFile := models.FileDta{File: *uFile, Description: params["desc"]}
			data, err := json.Marshal(mFile)
			if err != nil {
				return nil, err
			}
			params["data"] = string(data)
		} else {
			params["data"] = file.Data.SendData()
		}
	}

	var apiResp models.APIResponse
	resp, err := bot.Client.R().
		SetFormData(params).
		SetHeader("token", bot.Token).
		Post(fmt.Sprintf(bot.apiEndpoint, c.method()))

	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &apiResp)
	if err != nil {
		return nil, err
	}
	if apiResp.Error != "" {
		return &apiResp, &models.Error{Message: apiResp.Error}
	}

	return &apiResp, nil
}

func (bot *BotAPI) UploadFile(params Params, file RequestFile) (*models.File, error) {
	w := &bytes.Buffer{}
	m := multipart.NewWriter(w)

	for field, value := range params {
		if err := m.WriteField(field, value); err != nil {
			return nil, err
		}
	}

	name, reader, err := file.Data.UploadData()
	if err != nil {
		return nil, err
	}

	part, err := m.CreateFormFile(file.Name, filepath.Base(name))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, reader); err != nil {
		return nil, err
	}
	m.Close()

	var uploadedFile models.File
	resp, err := bot.Client.R().
		SetHeader("Content-Type", m.FormDataContentType()).
		SetBody(w).
		SetResult(&uploadedFile).
		Post(fmt.Sprintf(bot.apiEndpoint, "upload"))
	if err != nil {
		return nil, err
	}

	if uploadedFile.SID == "" {
		var apiResp models.APIResponse
		if err := json.Unmarshal(resp.Body(), &apiResp); err != nil {
			return nil, err
		}
		if apiResp.Error != "" {
			return nil, &models.Error{Message: apiResp.Error}
		}
	}

	return &uploadedFile, nil
}

func (bot *BotAPI) Send(c Chattable) (models.Message, error) {
	resp, err := bot.Request(c)
	if err != nil {
		return models.Message{}, err
	}
	if resp.Error != "" {
		return models.Message{}, errors.New(resp.Error)
	}

	msg := models.Message{MessageID: resp.MessageId}
	params, err := c.params()
	if err == nil {
		if chatIDStr := params.GetParam("chat_id"); chatIDStr != "" {
			if chatID, err := strconv.ParseInt(chatIDStr, 10, 64); err == nil {
				msg.ChatID = chatID
			}
		}
	}
	return msg, nil
}

func (bot *BotAPI) MultiSend(chattables ...Chattable) ([]models.Message, []error) {
	var messages []models.Message
	var errs []error
	for _, c := range chattables {
		msg, err := bot.Send(c)
		if err != nil {
			errs = append(errs, err)
		}
		messages = append(messages, msg)
	}
	return messages, errs
}
