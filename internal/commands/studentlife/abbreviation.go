package studentlife

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jackmerrill/hampbot/internal/utils/config"
	"github.com/jackmerrill/hampbot/internal/utils/embed"
	"github.com/zekroTJA/shireikan"
)

type Abbreviation struct {
}

func (c *Abbreviation) GetInvokes() []string {
	return []string{"abbreviation", "abbrev", "abbr"}
}

func (c *Abbreviation) GetDescription() string {
	return "WHAT DOES IT MEAN???"
}

func (c *Abbreviation) GetHelp() string {
	return "`abbreviation [acronym]` - Returns the meaning of a Hampshire College acronym"
}

func (c *Abbreviation) GetGroup() string {
	return config.GroupStudentLife
}

func (c *Abbreviation) GetDomainName() string {
	return "hamp.studentlife.abbreviation"
}

func (c *Abbreviation) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}
func (c *Abbreviation) IsExecutableInDMChannels() bool {
	return true
}

var Abbreviations = map[string]string{
	"APL":    "Airport Lounge",
	"ASH":    "Adele Simmons Hall",
	"BOT":    "Board of Trustees",
	"CAPES":  "Community Advocacy, Prevention and Education, Safety",
	"CASA":   "Center for Academic Support and Advising",
	"CC":     "Cultural Center (short for Lebrón-Wiggins-Pran Cultural Center)",
	"CFF":    "Center for Feminisms",
	"CMC":    "Collaborative Modeling Center",
	"CoCo":   "Community Commons (located in the library—off of the Airport Lounge, and formerly called the Knowledge Commons)",
	"CRB":    "Community Review Board",
	"CTL":    "Center for Teaching and Learning in Cole Science Center on the 1st Floor",
	"DLR":    "Dakin Living Room (inside Dakin Student Life Center)",
	"DSLC":   "Dakin Student Life Center",
	"EDH":    "Emily Dickinson Hall",
	"ELC":    "Early Learning Center",
	"ELH":    "East Lecture Hall (room in Franklin Patterson Hall)",
	"CEPC":   "Curriculum and Educational Policy Committee",
	"ESSP":   "Environmental Studies and Sustainability Program",
	"FPH":    "Franklin Patterson Hall",
	"FPV":    "Film, Photography, and Video Program",
	"G":      "Greenwich and Enfield Houses (one shared office)",
	"E":      "Greenwich and Enfield Houses (one shared office)",
	"GEO":    "Global Education Office",
	"HCFC":   "Hampshire College Farm Center",
	"HCSU":   "Hampshire College Student Union",
	"HLP":    "Holistic Learning Program",
	"IDBM":   "Identity-Based Mod/Hall",
	"IDBH":   "Identity-Based Hall",
	"RLSE":   "Residence Life and Student Engagement",
	"HR":     "Human Resources",
	"HRP":    "Hampshire Research Project",
	"ISS":    "International Student Services",
	"IT":     "Information Technology",
	"JEA":    "Justice, Equity, and Antiracism",
	"JLC":    "Jerome Liebling Center for Film, Photography, and Video",
	"KERN":   "The R.W. Kern Center",
	"LC":     "Learning Collaborative",
	"MDB":    "Music and Dance Building",
	"MLR":    "Merrill Living Room (inside Merrill Student Life Center)",
	"NSE":    "New Student Experience",
	"OARS":   "Office of Accessibility Resources and Services",
	"OPRA":   "Outdoor Programs, Recreation, and Athletics",
	"PAWSS":  "Peace and World Security Studies",
	"QCAC":   "Queer Community Alliance Center",
	"RCC":    "Robert Crown Center",
	"SAC":    "Staff Advocacy Committee",
	"SOURCE": "Students of Under-Represented Cultures and Ethnicities",
	"SPARC":  "Supporting Your Purpose through Action, Resources, and Connections (formerly CORC = Career Options Resource Center)",
	"WLH":    "West Lecture Hall (room in Franklin Patterson Hall)",
}

func (c *Abbreviation) Exec(ctx shireikan.Context) error {
	acronym := ctx.GetArgs().Get(0).AsString()

	if acronym == "" {
		ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
			Embed:     embed.NewErrorEmbed(ctx).SetDescription("Please specify an acronym to get the meaning of.\n\n**Usage:** `abbreviation [acronym]`").MessageEmbed,
			Reference: ctx.GetMessage().Reference(),
		})
		return nil
	}

	acronym = strings.ToUpper(acronym)

	if meaning, ok := Abbreviations[acronym]; ok {
		ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
			Embed: embed.NewEmbed().
				SetTitle("Hampshire College Acronym").
				SetDescription(fmt.Sprintf("**%s**\n\n%s", acronym, meaning)).
				SetColor(0x00ff00).
				MessageEmbed,
			Reference: ctx.GetMessage().Reference(),
		})
		return nil
	}

	ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
		Embed:     embed.NewErrorEmbed(ctx).SetDescription("Please specify a valid acronym to get the meaning of.\n\n**Usage:** `abbreviation [acronym]`").MessageEmbed,
		Reference: ctx.GetMessage().Reference(),
	})
	return nil
}
