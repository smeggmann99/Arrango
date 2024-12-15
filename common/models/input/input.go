// common/models/input/input.go
package input

/* Definitions
Division: A division is a group of students, with each division having a set
of subjects that need to be scheduled, each division has a weight that determines how important it is
to satisfy its constraints, the higher the weight, the more important it is to satisfy the constraints.
They can be split into groups for each subject, e.g. english could be split into two groups, one group could
be taught in the morning and the other in the afternoon etc.

Subject: A subject is a topic that is taught to a division, with each subject having a set of constraints
that need to be satisfied, such as the number of hours of that subject that should be placed in the timetable,
the placement of the subject in the timetable, the teacher that should teach the subject, the classroom that
the subject should be taught in, and the group that the division is split into for that subject.

Teacher: A teacher is a person that teaches a subject to a division, with each teacher having a set of
constraints that need to be satisfied, such as the number of hours that the teacher should be teaching,
the divisions that the teacher should be teaching, the classrooms that the teacher should be teaching in,
and the subjects that the teacher should be teaching.

Classroom: A classroom is a room that a subject is taught in.
*/ 

/* Scheduling Algorithm
The scheduling algorithm is a genetic algorithm that generates a timetable for each division,
with each division having a set of subjects that need to be scheduled. The algorithm generates
a population of timetables, each with a random assignment of subjects to timeslots, I thought
about using GPT-4o-mini to generate the initial population, but I decided to use a random assignment
as a starting point, so maybe in the future I will implement it as an option. The algorithm
evaluates the fitness of each timetable in the population and evolves the population over a number
of generations to find the best timetable that satisfies the constraints of the input data.
The algorithm uses a fitness function to determine how well a timetable satisfies the constraints,
with a lower fitness value indicating a better timetable. The algorithm uses selection, crossover,
and mutation to evolve the population, with the best timetables from each generation being selected
to form the next generation. The algorithm stops when a timetable with a fitness of 0 is found,
indicating that all constraints have been satisfied. The algorithm outputs the best timetables found.
*/

/* Scheduling Rules
Hard Constraints (must be satisfied, otherwise the timetable is invalid):
No gaps in timetables of divisions - no division ever has a gap during the day.
No teacher overlaps                - no teacher is ever in two different divisions or classrooms at the same time.
No classroom overlaps              - no classroom is ever assigned to two different divisions at the same time.
Block lessons                      - lessons are in the exact allocations per week that have been assigned. No more, no less.
Preferable classrooms              - each division has a list of classrooms that are preferable for each subject.
Teacher preferences                - each division has a list of teachers that are preferable for each subject.

Soft Constraints (should be minimized if possible):
No gaps in timetables of teachers  - no teacher ever has a gap during the day.
No unbalanced timetables           - no division has a significantly different number of hours per day compared to other days.
*/

// Determines where the subject should be placed in the timetable
type SubjectPlacementType string

const (
	SubjectPlacementAny    SubjectPlacementType = "any"    // Anywhere in the timetable
	SubjectPlacementEdges  SubjectPlacementType = "edges"  // At the beginning or end of the timetable
	SubjectPlacementCenter SubjectPlacementType = "middle" // In the middle of the timetable
)

type SubjectsGroupType string

const (
	SubjectsGroupNone    SubjectsGroupType = "none"
	SubjectsGroupOne     SubjectsGroupType = "one"
	SubjectsGroupTwo     SubjectsGroupType = "two"
	SubjectsGroupThree   SubjectsGroupType = "three"
)

type GlobalSubject string
type Classroom string
type Teacher string

type Subject struct {
	GlobalSubject *GlobalSubject       `json:"global_subject,omitempty"`
	// The number of consecutive hours that the subject should be placed in the timetable, indexed by the day of the week,
	// e.g. [2, 1, 2, 1, 2] means that the subject should be placed in two consecutive hours on any day of the week, one hour
	// on any other day of the week, two consecutive hours on any day of the week, one hour on any other day of the week,
	// and two consecutive hours on any day of the week, respectively, it can't be placed in the same day twice
	// e.g. [2, 1] means that the subject should be placed in two consecutive hours on any day of the week and one hour on any other day of the week
	Allocation    [5]uint              `json:"allocation,omitempty"`
	// Determines where the subject should be placed in the timetable
	Placement     SubjectPlacementType `json:"placement,omitempty"`
	// The teacher that should teach the subject in that division
	Teacher       *Teacher             `json:"teacher,omitempty"`
	// The classrooms that the subject can be taught in, if it's empty, then any available classroom can be used, otherwise, the subject should be taught in one of the classrooms
	Classrooms    []*Classroom         `json:"classrooms,omitempty"`
	// The group that the division is split into for that subject
	// e.g. english could be split into two groups, one group could be taught in the morning and the other in the afternoon
	// e.g. electronics could be split into three groups, one group could be taught on Monday, the second on Wednesday, and the third on Friday
	// e.g. polish is not split into groups, so the group is none, and the subject is taught to the whole division at the same time
	Group         SubjectsGroupType    `json:"group,omitempty"`
}

type Division struct {
	Name     string    `json:"name,omitempty"`
	// The weight of the division, used to determine how important it is to satisfy the constraints of the division
	// the higher the weight, the more important it is to satisfy the constraints of the division and the earlier
	// the division is scheduled in the timetable (that division should be scheduled first, so they start their day early)
	Weight   uint      `json:"weight,omitempty"`
	// The grouping of the division for each subject, indexed by the subject ID
	Subjects []Subject `json:"subjects,omitempty"` // The subjects that the division has
}

type InputData struct {
	// The global subjects that are available, each division has a subset of these subjects with different allocations etc.
	GlobalSubjects         []GlobalSubject `json:"global_subjects,omitempty"`
	Classrooms             []Classroom     `json:"classrooms,omitempty"`
	Teachers               []Teacher       `json:"teachers,omitempty"`
	Divisions              []Division      `json:"divisions,omitempty"`
}

var GlobalSubjects = []GlobalSubject{
	"Zajęcia w ZPKZ",
	"matematyka",
	"urz.i.syst.m",
	"j.niemiecki",
	"j.polski",
	"historia",
	"godz.wych",
	"religia",
	"wf",
	"fizyka",
	"WOS",
	"j.ang",
	"prac.apk.mob",
	"pr.te.do.apk",
	"prog.apk.web",
	"prog.apk.mob",
	"prog.str.obi",
}

var Classrooms = []Classroom{
	"sg4", "sg3", "sj1", "sj7", "14", "12", "47", "44", "4", "SKat", "7", "sj2", 
	"sj6", "ckz", "39", "107", "108", "42", "45", "38", "52", "40", "46",
}

var Teachers = []Teacher{
	"Be", "gr", "Sw", "kl", "LJ", "PO", "Su", "Kc", "LW", "Na", "Ba", "Bm", 
	"Ckz", "WG", "Kv", "Mw", "LI", "Sr", "GÓ", "Mt", "Aw", "Kł", "Wo", "tl",
}

var Divisions = []Division{
	{
		Name:   "Division 0",
		Weight: 1,
		Subjects: []Subject{
			// Zajęcia w ZPKZ
			{
				GlobalSubject: &GlobalSubjects[0],
				Allocation:    [5]uint{4, 4},
				Placement:     SubjectPlacementEdges,
				Teacher:       &Teachers[12],
				Classrooms:    []*Classroom{&Classrooms[13]},
				Group:         SubjectsGroupNone,
			},
			// matematyka
			{
				GlobalSubject: &GlobalSubjects[1],
				Allocation:    [5]uint{1, 2, 2},
				Placement:     SubjectPlacementCenter,
				Teacher:       &Teachers[4], // LJ
				Classrooms:    []*Classroom{&Classrooms[4], &Classrooms[10]}, // 14, 7
				Group:         SubjectsGroupNone,
			},
			// urz.i.syst.m
			{
				GlobalSubject: &GlobalSubjects[2],
				Allocation:    [5]uint{2, 2, 1},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[5], // PO
				Classrooms:    []*Classroom{&Classrooms[5], &Classrooms[4]}, // 12
				Group:         SubjectsGroupNone,
			},
			// j.niemiecki group 1
			{
				GlobalSubject: &GlobalSubjects[3],
				Allocation:    [5]uint{1},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[10], // Ba
				Classrooms:    []*Classroom{&Classrooms[11], &Classrooms[12]}, // sj2, sj6
				Group:         SubjectsGroupOne,
			},
			// j.niemiecki group 2
			{
				GlobalSubject: &GlobalSubjects[3],
				Allocation:    [5]uint{1},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[11], // Bm
				Classrooms:    []*Classroom{&Classrooms[11], &Classrooms[12]}, // sj2, sj6
				Group:         SubjectsGroupTwo,
			},
			// j.polski
			{
				GlobalSubject: &GlobalSubjects[4],
				Allocation:    [5]uint{2, 2},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[6], // Su
				Classrooms:    []*Classroom{&Classrooms[6]}, // 47
				Group:         SubjectsGroupNone,
			},
			// historia
			{
				GlobalSubject: &GlobalSubjects[6],
				Allocation:    [5]uint{1},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[7], // Kc
				Classrooms:    []*Classroom{&Classrooms[7]}, // 44
				Group:         SubjectsGroupNone,
			},
			// TODO: Implement placement constraints
			// godz.wych
			{
				GlobalSubject: &GlobalSubjects[7],
				Allocation:    [5]uint{1},
				Placement:     SubjectPlacementEdges,
				Teacher:       &Teachers[0], // Be
				Classrooms:    []*Classroom{&Classrooms[8]}, // 4
				Group:         SubjectsGroupNone,
			},
			// religia
			{
				GlobalSubject: &GlobalSubjects[8],
				Allocation:    [5]uint{2},
				Placement:     SubjectPlacementEdges,
				Teacher:       &Teachers[8], // LW
				Classrooms:    []*Classroom{&Classrooms[9]}, // SKat
				Group:         SubjectsGroupNone,
			},
			// TODO: Add classroom capacity constraints
			// wf group 1
			{
				GlobalSubject: &GlobalSubjects[9],
				Allocation:    [5]uint{2, 1},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[0], // Be
				Classrooms:    []*Classroom{&Classrooms[0], &Classrooms[1]}, // sg4, sg3
				Group:         SubjectsGroupOne,
			},
			// wf group 2
			{
				GlobalSubject: &GlobalSubjects[9],
				Allocation:    [5]uint{2, 1},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[1], // gr
				Classrooms:    []*Classroom{&Classrooms[0], &Classrooms[1]}, // sg4, sg3
				Group:         SubjectsGroupTwo,
			},
			// fizyka
			{
				GlobalSubject: &GlobalSubjects[10],
				Allocation:    [5]uint{2},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[9], // Na
				Classrooms:    []*Classroom{&Classrooms[10]}, // 7
				Group:         SubjectsGroupNone,
			},
			// WOS
			{
				GlobalSubject: &GlobalSubjects[12],
				Allocation:    [5]uint{1},
				Placement:     SubjectPlacementEdges,
				Teacher:       &Teachers[7], // Kc
				Classrooms:    []*Classroom{&Classrooms[7]}, // 44
				Group:         SubjectsGroupNone,
			},
			// j.ang group 1
			{
				GlobalSubject: &GlobalSubjects[13],
				Allocation:    [5]uint{1, 2},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[2], // Sw
				Classrooms:    []*Classroom{&Classrooms[2], &Classrooms[3]}, // sj1, sj7
				Group:         SubjectsGroupOne,
			},
			// j.ang group 2
			{
				GlobalSubject: &GlobalSubjects[13],
				Allocation:    [5]uint{1, 2},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[3], // kl
				Classrooms:    []*Classroom{&Classrooms[2], &Classrooms[3]}, // sj1, sj7
				Group:         SubjectsGroupTwo,
			},
		},
	},
	{
		Name:   "Division 1",
		Weight: 1,
		Subjects: []Subject{
			// r_matematyka
			{
				GlobalSubject: &GlobalSubjects[5], // r_matematyka
				Allocation:    [5]uint{1, 0, 0, 0, 0},
				Placement:     SubjectPlacementEdges,
				Teacher:       &Teachers[4], // Lj
				Classrooms:    []*Classroom{&Classrooms[4]}, // 14
				Group:         SubjectsGroupNone,
			},
			// matematyka
			{
				GlobalSubject: &GlobalSubjects[1], // matematyka
				Allocation:    [5]uint{0, 2, 1, 0, 0},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[4], // Lj
				Classrooms:    []*Classroom{&Classrooms[4]}, // 14
				Group:         SubjectsGroupNone,
			},
			// wf group 1
			{
				GlobalSubject: &GlobalSubjects[9], // wf
				Allocation:    [5]uint{1, 0, 0, 0, 1},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[22], // Kł
				Classrooms:    []*Classroom{&Classrooms[2]}, // sj1
				Group:         SubjectsGroupOne,
			},
			// wf group 2
			{
				GlobalSubject: &GlobalSubjects[9], // wf
				Allocation:    [5]uint{1, 0, 0, 0, 1},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[23], // Wo
				Classrooms:    []*Classroom{&Classrooms[3]}, // sj7
				Group:         SubjectsGroupTwo,
			},
			// j.polski
			{
				GlobalSubject: &GlobalSubjects[4], // j.polski
				Allocation:    [5]uint{2, 1, 0, 0, 0},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[6], // Su
				Classrooms:    []*Classroom{&Classrooms[6]}, // 47
				Group:         SubjectsGroupNone,
			},
			// historia
			{
				GlobalSubject: &GlobalSubjects[6], // historia
				Allocation:    [5]uint{0, 0, 1, 0, 0},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[7], // Kc
				Classrooms:    []*Classroom{&Classrooms[7]}, // 44
				Group:         SubjectsGroupNone,
			},
			// prog.str.obi
			{
				GlobalSubject: &GlobalSubjects[18], // prog.str.obi
				Allocation:    [5]uint{0, 0, 2, 0, 0},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[17], // Sr
				Classrooms:    []*Classroom{&Classrooms[5], &Classrooms[20]}, // Sr_12, 52
				Group:         SubjectsGroupNone,
			},
			// WOS
			{
				GlobalSubject: &GlobalSubjects[12], // WOS
				Allocation:    [5]uint{0, 1, 0, 0, 0},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[18], // GÓ
				Classrooms:    []*Classroom{&Classrooms[10]}, // 45
				Group:         SubjectsGroupNone,
			},
			// prog.apk.web
			{
				GlobalSubject: &GlobalSubjects[16], // prog.apk.web
				Allocation:    [5]uint{0, 0, 1, 1, 1},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[16], // LI
				Classrooms:    []*Classroom{&Classrooms[8], &Classrooms[21]}, // LI_7, 46
				Group:         SubjectsGroupNone,
			},
			// prog.apk.mob
			{
				GlobalSubject: &GlobalSubjects[17], // prog.apk.mob
				Allocation:    [5]uint{1, 0, 0, 0, 0},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[17], // Sr
				Classrooms:    []*Classroom{&Classrooms[19]}, // 38
				Group:         SubjectsGroupNone,
			},
			// pr.te.do.apk group 1
			{
				GlobalSubject: &GlobalSubjects[15], // pr.te.do.apk
				Allocation:    [5]uint{1, 0, 0, 0, 1},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[14], // WG
				Classrooms:    []*Classroom{&Classrooms[15]}, // 107
				Group:         SubjectsGroupOne,
			},
			// pr.te.do.apk group 2
			{
				GlobalSubject: &GlobalSubjects[15], // pr.te.do.apk
				Allocation:    [5]uint{1, 0, 0, 0, 1},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[15], // Kv
				Classrooms:    []*Classroom{&Classrooms[16]}, // 108
				Group:         SubjectsGroupTwo,
			},
			// religia
			{
				GlobalSubject: &GlobalSubjects[8], // religia
				Allocation:    [5]uint{1, 0, 0, 0, 0},
				Placement:     SubjectPlacementEdges,
				Teacher:       &Teachers[10], // LW
				Classrooms:    []*Classroom{&Classrooms[9]}, // SKat
				Group:         SubjectsGroupNone,
			},
			// godz.wych
			{
				GlobalSubject: &GlobalSubjects[7], // godz.wych
				Allocation:    [5]uint{0, 0, 0, 1, 0},
				Placement:     SubjectPlacementEdges,
				Teacher:       &Teachers[15], // Mw
				Classrooms:    []*Classroom{&Classrooms[17]}, // 42
				Group:         SubjectsGroupNone,
			},
			// j.ang group 1
			{
				GlobalSubject: &GlobalSubjects[13], // j.ang
				Allocation:    [5]uint{0, 2, 0, 0, 0},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[19], // Mt
				Classrooms:    []*Classroom{&Classrooms[2]}, // sj1
				Group:         SubjectsGroupOne,
			},
			// j.ang group 2
			{
				GlobalSubject: &GlobalSubjects[13], // j.ang
				Allocation:    [5]uint{0, 2, 0, 0, 0},
				Placement:     SubjectPlacementAny,
				Teacher:       &Teachers[20], // Aw
				Classrooms:    []*Classroom{&Classrooms[3]}, // sj7
				Group:         SubjectsGroupTwo,
			},
		},
	},
}

var ExampleInputData = InputData{
	GlobalSubjects: GlobalSubjects,
	Classrooms:     Classrooms,
	Teachers:       Teachers,
	Divisions:      Divisions,
}