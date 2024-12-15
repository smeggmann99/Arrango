// common/models/output/output.go
package output

/* Definitions
Timetable: A timetable is a schedule of subjects for a division for a week, with each day having
a set of sets of subjects that need to be scheduled, sometimes a subject can be split into multiple
groups that can be taught at the same time, e.g. english could be split into two groups, one group could
be taught in the morning and the other in the afternoon, or electronics could be split into three groups,
it might be possible to schedule them at the same time, but in different classrooms and with different teachers.

Timetables: The timetables for each division, indexed by the division ID.
*/

import (
	"smuggr.xyz/arrango/common/models/input"
)

type Subject struct {
	GlobalSubject *input.GlobalSubject     `json:"global_subject,omitempty"`
	Teacher       *input.Teacher           `json:"teacher,omitempty"`
	Classroom     *input.Classroom         `json:"classroom,omitempty"`
	Group         *input.SubjectsGroupType `json:"group,omitempty"`
}

type SubjectsGroup [3]Subject       // A group of subjects, which are taught at the same time, maximum 3
type Day           []SubjectsGroup  // A day's timetable
type Days          [5]Day           // A week's timetable

type OutputData struct {
	// The timetables for each division, indexed by the division index
	DivisionsTimetables []Days `json:"timetables,omitempty"`
}