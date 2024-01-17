package commands

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/jackmerrill/hampbot/internal/utils/config"
	"github.com/jackmerrill/hampbot/internal/utils/embed"
	"github.com/zekroTJA/shireikan"
)

type VerifyCommand struct {
}

func (c *VerifyCommand) GetInvokes() []string {
	return []string{"verify"}
}

func (c *VerifyCommand) GetDescription() string {
	return "Verify your student/staff/faculty status with your school email."
}

func (c *VerifyCommand) GetHelp() string {
	return "`verify` - `verify [email]`"
}

func (c *VerifyCommand) GetGroup() string {
	return config.GroupUtil
}

func (c *VerifyCommand) GetDomainName() string {
	return "hamp.util.verify"
}

func (c *VerifyCommand) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}
func (c *VerifyCommand) IsExecutableInDMChannels() bool {
	return true
}

type Verification struct {
	User           *discordgo.User
	Email          string
	Expires        time.Time
	IsStaffFaculty bool
}

var VerificationCodes map[string]Verification = make(map[string]Verification)

func (c *VerifyCommand) Exec(ctx shireikan.Context) error {
	// first check if the user is already verified (has the role)
	guildMember, err := ctx.GetSession().GuildMember(config.BotGuild, ctx.GetUser().ID)

	if err != nil {
		return err
	}

	for _, role := range guildMember.Roles {
		if role == config.VerifiedRoleId {
			ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
				Reference: ctx.GetMessage().Reference(),
				Embed:     embed.NewErrorEmbed(ctx).SetTitle("Already Verified").SetDescription("You are already verified.").MessageEmbed,
			})
			return nil
		}
	}

	// first arg is the email
	email := ctx.GetArgs().Get(0).AsString()

	// check if the email is valid, and if it is a hampshire.edu email
	if !strings.Contains(email, "@hampshire.edu") {
		ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
			Reference: ctx.GetMessage().Reference(),
			Embed:     embed.NewErrorEmbed(ctx).SetTitle("Invalid Email").SetDescription("Please use a valid @hampshire.edu email.").MessageEmbed,
		})
		return nil
	}

	code := uuid.New().String()

	err = SendEmail([]string{email}, code, ctx.GetUser())

	if err != nil {
		return err
	}

	// expires in 5 minutes
	expires := time.Now().Add(time.Minute * 5)

	m, err := ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
		Reference: ctx.GetMessage().Reference(),
		Embed:     embed.NewSuccessEmbed(ctx).SetTitle("Sent Verification Email").SetDescription("Waiting for you to verify...").AddField("Expires", fmt.Sprintf("<t:%d:R>", expires.Unix()), false).MessageEmbed,
	})

	if err != nil {
		return err
	}

	// check if the email is student or staff/faculty
	// student emails are in the format of their initials plus their first year (e.g. jmm18). To be safe, lets allow up to six initials plus the year (e.g. jmmmmmm18)
	// staff/faculty emails are in the format of their first initial plus their last name (e.g. jsmith), or their first and last initial plus their department (e.g. jsIA)

	hampnetUser := strings.Split(email, "@")[0]

	// use regex to check
	// student
	if match, _ := regexp.MatchString(`^[a-z]{1,6}[0-9]{2}$`, hampnetUser); match {
		VerificationCodes[code] = Verification{
			User:           ctx.GetUser(),
			Email:          email,
			Expires:        expires,
			IsStaffFaculty: false,
		}
	} else if match, _ := regexp.MatchString(`^[a-z]{1,2}[a-z]{1,2}[a-z]{1,2}[a-z]{1,2}$`, hampnetUser); match {
		VerificationCodes[code] = Verification{
			User:           ctx.GetUser(),
			Email:          email,
			Expires:        expires,
			IsStaffFaculty: true,
		}
	} else {
		// potentially student, but not in the format of a student email
		VerificationCodes[code] = Verification{
			User:           ctx.GetUser(),
			Email:          email,
			Expires:        expires,
			IsStaffFaculty: false,
		}
	}

	// wait for the user to verify
	for {
		if time.Now().After(expires) {
			ctx.GetSession().ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:      m.ID,
				Channel: ctx.GetChannel().ID,
				Embed:   embed.NewErrorEmbed(ctx).SetTitle("Verification Expired").SetDescription("Please try again.").MessageEmbed,
			})
		}

		_, ok := VerificationCodes[code]

		if !ok {
			ctx.GetSession().ChannelMessageEditComplex(&discordgo.MessageEdit{
				ID:      m.ID,
				Channel: ctx.GetChannel().ID,
				Embed:   embed.NewSuccessEmbed(ctx).SetTitle("Verified!").SetDescription("You can access the server now.").MessageEmbed,
			})
			break
		}

		time.Sleep(time.Second * 5)
	}

	return nil
}

func StartWebserver(session *discordgo.Session) {
	// listen to /verify

	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		if code == "" {
			http.Error(w, "Code is missing", http.StatusBadRequest)
			return
		}

		v, ok := VerificationCodes[code]

		if !ok {
			http.Error(w, "Code is invalid", http.StatusBadRequest)
			return
		}

		err := session.GuildMemberRoleAdd(config.BotGuild, v.User.ID, config.VerifiedRoleId)

		if err != nil {
			http.Error(w, "Failed to add role", http.StatusInternalServerError)
			return
		}

		// remove the code from the map
		delete(VerificationCodes, code)

		fmt.Fprintf(w, "Verified! You can close this tab now.")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK!")
	})

	http.ListenAndServe(":8080", nil)
}

func SendEmail(to []string, code string, discordUser *discordgo.User) error {
	sender := "automated@jackmerrill.com"

	user := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASSWORD")

	subject := "Hampshire Hangout Verification"
	tmpl, err := template.ParseFiles("verify.html")

	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)

	err = tmpl.Execute(buf, struct {
		Code            string
		DiscordUserName string
	}{
		Code:            code,
		DiscordUserName: discordUser.Username,
	})

	if err != nil {
		return err
	}

	request := Mail{
		Sender:  sender,
		To:      to,
		Subject: subject,
		Body:    buf.String(),
	}

	addr := "smtp.gmail.com:587"
	host := "smtp.gmail.com"

	msg := BuildMessage(request)
	auth := smtp.PlainAuth("", user, password, host)
	err = smtp.SendMail(addr, auth, sender, to, []byte(msg))

	if err != nil {
		return err
	}

	return nil
}

func BuildMessage(mail Mail) string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)

	return msg
}

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}
