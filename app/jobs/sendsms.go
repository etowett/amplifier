package jobs

import (
	"amplifier/app/db"
	"amplifier/app/entities"
	"amplifier/app/models"
	"amplifier/app/providers"
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/revel/modules/jobs/app/jobs"
	"github.com/revel/revel"
)

// Periodically count the users in the database.
type SMSSender struct{}

type Recipient struct {
	Number  string `json:"number"`
	Message string `json:"message"`
}

func (j SMSSender) Run() {
	if revel.Config.StringDefault("app.env", "local") == "local" {
		revel.AppLog.Infof("Can not send cron sms in local")
		return
	}

	rand.Seed(time.Now().UnixNano())
	reqNum := rand.Intn(100-4) + 2

	recipients := make([]*entities.SMSRecipient, 0)
	for _, rec := range j.getRecipients(reqNum) {
		recipients = append(recipients, &entities.SMSRecipient{
			Phone:   rec.Number,
			Message: rec.Message,
		})
	}

	// err := j.sendAt(recipients, reqNum)
	// if err != nil {
	// 	revel.AppLog.Infof("could not send to at: %v", err)
	// }
	// revel.AppLog.Infof("done sending at message")

	err := j.sendKomsner(recipients, reqNum)
	if err != nil {
		revel.AppLog.Infof("could not send to komsner: %v", err)
	}
	revel.AppLog.Infof("done sending komsner message")
}

func init() {
	revel.OnAppStart(func() {
		jobs.Schedule("@every 15m", SMSSender{})
	})
}

// func (j SMSSender) sendAt(
// 	recipients []*entities.SMSRecipient,
// 	reqNum int,
// ) error {
// 	ctx := context.Background()
// 	newCredential := &models.Credential{}
// 	credential, err := newCredential.ByApp(ctx, db.DB(), "apisim")
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			revel.AppLog.Infof("could not find any at creds for at send")
// 			return nil
// 		}
// 		return fmt.Errorf("could not get creds by app %v: %v", "apisim", err)
// 	}
// 	africasTalkingSender := providers.NewAfricasTalkingSenderWithParameters(
// 		credential.Url,
// 		credential.Username,
// 		credential.Password,
// 	)
// 	_, err = africasTalkingSender.Send(&entities.SendRequest{
// 		SenderID: "AmplifierJob",
// 		Message:  fmt.Sprintf("Hello to %v from amplifier job at %v", reqNum, time.Now()),
// 		To:       recipients,
// 	})
// 	if err != nil {
// 		return fmt.Errorf("failed to send to at: %v", err)
// 	}
// 	return nil
// }

func (j SMSSender) sendKomsner(
	recipients []*entities.SMSRecipient,
	reqNum int,
) error {
	ctx := context.Background()
	newCredential := &models.Credential{}
	credential, err := newCredential.ByApp(ctx, db.DB(), "komsner")
	if err != nil {
		if err == sql.ErrNoRows {
			revel.AppLog.Infof("could not find any komsner creds for komsner send")
			return nil
		}
		return fmt.Errorf("could not get creds by app %v: %v", "komsner", err)
	}

	values := []bool{true, false}
	rand.Seed(time.Now().UnixNano())
	isMulti := values[rand.Intn(len(values))]
	smsBody := fmt.Sprintf(
		"Send multi: %v to %d recs at %s.",
		isMulti,
		reqNum,
		time.Now().String()[0:19],
	)

	komsnerSender := providers.NewKomsnerSenderWithParameters(
		credential.Url,
		fmt.Sprintf("%v:%v", credential.Username, credential.Password),
	)
	_, err = komsnerSender.Send(&entities.SendRequest{
		SenderID: "AmplifierJob",
		Message:  smsBody,
		To:       recipients,
		Multi:    isMulti,
	})
	if err != nil {
		return fmt.Errorf("failed to send to komsner: %v", err)
	}
	return nil
}

func (j SMSSender) getRecipients(num int) (recipients []*Recipient) {
	for i := 0; i < num; i++ {
		recipients = append(recipients, &Recipient{
			Number:  j.getPhone(),
			Message: fmt.Sprintf("Hello to you %v", i),
		})
	}
	return recipients
}

func (j SMSSender) getPhone() string {
	prefix := []string{"+2547", "+2541"}
	// prefix := []string{"+2557", "+2536", "+2547", "+2119", "+2568", "+2541"}
	net := []string{
		"16", "17", "18", "20", "21", "22", "23", "25", "96", "27",
	}

	rand.Seed(time.Now().UnixNano())
	destCtry := prefix[rand.Intn(len(prefix))]
	destNet := net[rand.Intn(len(net))]
	randNum := strconv.Itoa(111111 + rand.Intn(999999-111111))

	return destCtry + destNet + randNum
}
