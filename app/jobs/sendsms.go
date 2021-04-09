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

func (j SMSSender) Run() {
	err := j.sendAt()
	if err != nil {
		revel.AppLog.Infof("could not send to at: %v", err)
	}
	revel.AppLog.Infof("done sending at message")

	err = j.sendKomsner()
	if err != nil {
		revel.AppLog.Infof("could not send to komsner: %v", err)
	}
	revel.AppLog.Infof("done sending komsner message")
}

func init() {
	revel.OnAppStart(func() {
		jobs.Schedule("@every 1m", SMSSender{})
	})
}

func (j SMSSender) sendAt() error {
	ctx := context.Background()

	newCredential := &models.Credential{}
	credential, err := newCredential.ByApp(ctx, db.DB(), "apisim")
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("could not find any at creds for at send")
		}
		return fmt.Errorf("could not get creds by app %v: %v", "apisim", err)
	}

	africasTalkingSender := providers.NewAfricasTalkingSenderWithParameters(
		credential.Url,
		credential.Username,
		credential.Password,
	)

	rand.Seed(time.Now().UnixNano())
	recs := j.getRecs(rand.Intn(1000 - 10))

	recipients := make([]*entities.SMSRecipient, 0)
	for _, num := range recs {
		recipients = append(recipients, &entities.SMSRecipient{Phone: num})
	}

	_, err = africasTalkingSender.Send(&entities.SendRequest{
		SenderID: "AmplifierJob",
		Message:  fmt.Sprintf("Hello from amplifier job at %v", time.Now()),
		To:       recipients,
	})
	if err != nil {
		return fmt.Errorf("failed to send to at: %v", err)
	}
	return nil
}

func (j SMSSender) sendKomsner() error {
	return nil
}

func (j SMSSender) getRecs(num int) (numbers []string) {
	for i := 0; i < num; i++ {
		numbers = append(numbers, j.getPhone())
	}
	return numbers
}

func (j SMSSender) getPhone() string {
	prefix := []string{"+2557", "+2536", "+2547", "+2119", "+2568", "+2541"}
	net := []string{
		"16", "17", "18", "20", "21", "22", "23", "25", "96", "27",
	}

	rand.Seed(time.Now().UnixNano())
	destCtry := prefix[rand.Intn(len(prefix))]
	destNet := net[rand.Intn(len(net))]
	randNum := strconv.Itoa(111111 + rand.Intn(999999-111111))

	return destCtry + destNet + randNum
}
