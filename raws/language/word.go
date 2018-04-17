package language // import "github.com/BenLubar/dfide/raws/language"

type Word struct {
	ID        string     `raws:"1"`
	Noun      *Noun      `raws:"NOUN"`
	Adjective *Adjective `raws:"ADJ"`
	Prefix    *Prefix    `raws:"PREFIX"`
	Verb      *Verb      `raws:"VERB"`
}

type Noun struct {
	Singular string `raws:"1"`
	Plural   string `raws:"2"`

	FrontCompoundSingular bool `raws:"FRONT_COMPOUND_NOUN_SING"`
	FrontCompoundPlural   bool `raws:"FRONT_COMPOUND_NOUN_PLUR"`
	RearCompoundSingular  bool `raws:"REAR_COMPOUND_NOUN_SING"`
	RearCompoundPlural    bool `raws:"REAR_COMPOUND_NOUN_PLUR"`
	TheSingular           bool `raws:"THE_NOUN_SING"`
	ThePlural             bool `raws:"THE_NOUN_PLUR"`
	TheCompoundSingular   bool `raws:"THE_COMPOUND_NOUN_SING"`
	TheCompoundPlural     bool `raws:"THE_COMPOUND_NOUN_PLUR"`
	OfSingular            bool `raws:"OF_NOUN_SING"`
	OfPlural              bool `raws:"OF_NOUN_PLUR"`
}

type Adjective struct {
	Adjective string `raws:"1"`
	Distance  uint8  `raws:"ADJ_DIST.1" min:"1" max:"7" default:"1"`

	FrontCompound bool `raws:"FRONT_COMPOUND_ADJ"`
	RearCompound  bool `raws:"REAR_COMPOUND_ADJ"`
	TheCompound   bool `raws:"THE_COMPOUND_ADJ"`
}

type Prefix struct {
	Prefix string `raws:"1"`

	FrontCompound bool `raws:"FRONT_COMPOUND_PREFIX"`
	TheCompound   bool `raws:"THE_COMPOUND_PREFIX"`
}

type Verb struct {
	PresentFirstPerson string `raws:"1"`
	PresentThirdPerson string `raws:"2"`
	Preterite          string `raws:"3"`
	PastParticiple     string `raws:"4"`
	PresentParticiple  string `raws:"5"`

	StandardVerb bool `raws:"STANDARD_VERB"`
}
