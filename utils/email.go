package utils

import (
	"backup/conf"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type SendBody struct {
	Action string `json:"action,omitempty"`
	Attrs  Attrs  `json:"attrs"`
}

type Attrs struct {
	Account string   `json:"account,omitempty"`
	To      []string `json:"to,omitempty"`
	Subject string   `json:"subject,omitempty"`
	Content string   `json:"content,omitempty"`
}

type Response struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func SendEmail(subject, content string) error {
	var buf bytes.Buffer
	wr := multipart.NewWriter(&buf)
	fw, err := wr.CreateFormField("body")
	if err != nil {
		return err
	}
	var sendBody = SendBody{
		Action: "deliver",
		Attrs: Attrs{
			Account: "qhjfzhglpt@dc.icbc.com.cn",
			To:      conf.GlobalCfg.Email.Receivers,
			Subject: subject,
			Content: content,
		},
	}
	data, err := json.Marshal(sendBody)
	if err != nil {
		return err
	}
	i, err := fw.Write(data)
	if err != nil {
		return err
	}
	if i != len(data) {
		return errors.New("write  sendbody less")
	}
	if err = wr.Close(); err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, conf.GlobalCfg.Email.Endpoint, &buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", wr.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("response code not 200:%d", resp.StatusCode))
	}
	rb, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var re Response
	if err = json.Unmarshal(rb, &re); err != nil {
		return err
	}
	if re.Code == 200 {
		return nil
	}
	return fmt.Errorf("send mail error:%s", re.Message)
}
