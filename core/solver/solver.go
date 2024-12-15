// core/solver/solver.go
package solver

import (
	"math/rand"
	"sort"
	"time"

	"smuggr.xyz/arrango/common/models/input"
	"smuggr.xyz/arrango/common/models/output"
)

type Solver struct {
	PopulationSize int
	Generations    int
	MutationRate   float64
}

type Individual struct {
	Timetables []output.Days // One timetable per division
}

func (s *Solver) Solve(in input.InputData) output.OutputData {
	rand.Seed(time.Now().UnixNano())

	pop := s.initializePopulation(in)

	bestIndividual := pop[0]
	bestFitness := s.fitness(bestIndividual, in)

	for g := 0; g < s.Generations; g++ {
		type fitInd struct {
			ind     Individual
			fitness int
		}
		fits := make([]fitInd, len(pop))
		for i, ind := range pop {
			f := s.fitness(ind, in)
			fits[i] = fitInd{ind, f}
			if f < bestFitness {
				bestFitness = f
				bestIndividual = ind
				if bestFitness == 0 {
					break
				}
			}
		}

		if bestFitness == 0 {
			break
		}

		sort.Slice(fits, func(i, j int) bool {
			return fits[i].fitness < fits[j].fitness
		})

		nextPop := make([]Individual, 0, s.PopulationSize)
		// selection: top half
		for i := 0; i < s.PopulationSize/2; i++ {
			nextPop = append(nextPop, fits[i].ind)
		}

		// Reproduction
		for len(nextPop) < s.PopulationSize {
			p1 := fits[rand.Intn(s.PopulationSize/2)].ind
			p2 := fits[rand.Intn(s.PopulationSize/2)].ind
			child := s.crossover(p1, p2)
			s.mutate(&child)
			nextPop = append(nextPop, child)
		}

		pop = nextPop
	}

	return output.OutputData{DivisionsTimetables: bestIndividual.Timetables}
}

// Extract chunks of subject allocations
type subjectChunk struct {
	subj input.Subject
	size uint
}

func (s *Solver) extractSubjectChunks(div input.Division) []subjectChunk {
	var chunks []subjectChunk
	for _, subj := range div.Subjects {
		for _, alloc := range subj.Allocation {
			if alloc > 0 {
				chunks = append(chunks, subjectChunk{
					subj: subj,
					size: alloc,
				})
			}
		}
	}
	return chunks
}

func (s *Solver) pickClassroom(subj input.Subject) *input.Classroom {
	if len(subj.Classrooms) > 0 {
		return subj.Classrooms[rand.Intn(len(subj.Classrooms))]
	}
	return nil
}

// Initialize a random individual with balanced day lengths for each division.
func (s *Solver) randomIndividual(in input.InputData) Individual {
	timetables := make([]output.Days, len(in.Divisions))

	for dIdx, div := range in.Divisions {
		// We start with empty days
		var divisionDays output.Days
		for i := 0; i < 5; i++ {
			divisionDays[i] = make([]output.SubjectsGroup, 0)
		}

		requiredChunks := s.extractSubjectChunks(div)

		// Place chunks in the day with the fewest groups so far, to keep balanced
		for _, chunk := range requiredChunks {
			// We need to place 'chunk.size' consecutive hours for the subject
			// Pick a day that currently has the least number of groups
			dayIdx := s.pickLeastLoadedDay(divisionDays)
			// Append chunk.size groups with this subject
			for i := uint(0); i < chunk.size; i++ {
				sg := output.SubjectsGroup{}
				sg[0] = output.Subject{
					GlobalSubject: chunk.subj.GlobalSubject,
					Teacher:       chunk.subj.Teacher,
					Classroom:     s.pickClassroom(chunk.subj),
					Group:         &chunk.subj.Group,
				}
				divisionDays[dayIdx] = append(divisionDays[dayIdx], sg)
			}
		}

		timetables[dIdx] = divisionDays
	}

	return Individual{Timetables: timetables}
}

// pickLeastLoadedDay returns the index of the day with the fewest subjects groups
func (s *Solver) pickLeastLoadedDay(days output.Days) int {
	minLoad := len(days[0])
	minDay := 0
	for i := 1; i < 5; i++ {
		if len(days[i]) < minLoad {
			minLoad = len(days[i])
			minDay = i
		}
	}
	return minDay
}

func (s *Solver) initializePopulation(in input.InputData) []Individual {
	pop := make([]Individual, s.PopulationSize)
	for i := 0; i < s.PopulationSize; i++ {
		pop[i] = s.randomIndividual(in)
	}
	return pop
}

func (s *Solver) fitness(ind Individual, in input.InputData) int {
	score := 0

	// Check teacher/classroom overlaps
	type slotKey struct {
		day  int
		slot int
	}
	teacherUsed := make(map[slotKey]map[input.Teacher]bool)
	classroomUsed := make(map[slotKey]map[input.Classroom]bool)

	for _, divTT := range ind.Timetables {
		for day := 0; day < 5; day++ {
			for slot, sg := range divTT[day] {
				tk := slotKey{day: day, slot: slot}
				for _, subj := range sg {
					if subj.GlobalSubject == nil {
						continue
					}
					if subj.Teacher != nil {
						if teacherUsed[tk] == nil {
							teacherUsed[tk] = make(map[input.Teacher]bool)
						}
						if teacherUsed[tk][*subj.Teacher] {
							score += 1000 // Teacher overlap
						} else {
							teacherUsed[tk][*subj.Teacher] = true
						}
					}
					if subj.Classroom != nil {
						if classroomUsed[tk] == nil {
							classroomUsed[tk] = make(map[input.Classroom]bool)
						}
						if classroomUsed[tk][*subj.Classroom] {
							score += 1000 // Classroom overlap
						} else {
							classroomUsed[tk][*subj.Classroom] = true
						}
					}
				}
			}
		}
	}

	// Check allocations are met
	for dIdx, div := range in.Divisions {
		requiredChunks := s.extractSubjectChunks(div)
		// Copy needed counts
		remaining := make([]subjectChunk, len(requiredChunks))
		copy(remaining, requiredChunks)

		for day := 0; day < 5; day++ {
			for _, sg := range ind.Timetables[dIdx][day] {
				for _, subj := range sg {
					if subj.GlobalSubject == nil {
						continue
					}
					for i := range remaining {
						if remaining[i].subj.GlobalSubject == subj.GlobalSubject &&
							remaining[i].subj.Teacher == subj.Teacher {
							// placed an hour
							if remaining[i].size > 0 {
								remaining[i].size--
							}
						}
					}
				}
			}
		}

		// penalty for not meeting required allocations
		for _, c := range remaining {
			if c.size > 0 {
				score += int(c.size) * 500
			}
		}
	}

	// No gaps in division timetables:
	// Since we directly appended chunks, no "empty slots" were created.
	// Each subjects group is consecutive. So no internal gaps by definition.
	// If we considered gaps as missing groups, we would have introduced them ourselves.
	// Hence no penalty needed here.

	// Soft constraints: Unbalanced day distribution within a division
	// Check difference in day loads (number of groups per day)
	for dIdx := range ind.Timetables {
		dayCounts := make([]int, 5)
		for day := 0; day < 5; day++ {
			dayCounts[day] = len(ind.Timetables[dIdx][day])
		}
		minC, maxC := dayCounts[0], dayCounts[0]
		for _, c := range dayCounts[1:] {
			if c < minC {
				minC = c
			}
			if c > maxC {
				maxC = c
			}
		}
		if maxC-minC > 4 {
			score += (maxC - minC) * 5
		}
	}

	return score
}

func (s *Solver) crossover(p1, p2 Individual) Individual {
	child := Individual{
		Timetables: make([]output.Days, len(p1.Timetables)),
	}
	copy(child.Timetables, p1.Timetables)
	if len(p1.Timetables) > 0 {
		dx := rand.Intn(len(p1.Timetables))
		for i := 0; i < 2; i++ {
			day := rand.Intn(5)
			child.Timetables[dx][day] = p2.Timetables[dx][day]
		}
	}
	return child
}

func (s *Solver) mutate(ind *Individual) {
	if rand.Float64() > s.MutationRate {
		return
	}
	// Randomly pick a division/day and swap two slots if possible
	dx := rand.Intn(len(ind.Timetables))
	day := rand.Intn(5)
	if len(ind.Timetables[dx][day]) > 1 {
		slot1 := rand.Intn(len(ind.Timetables[dx][day]))
		slot2 := rand.Intn(len(ind.Timetables[dx][day]))
		ind.Timetables[dx][day][slot1], ind.Timetables[dx][day][slot2] = ind.Timetables[dx][day][slot2], ind.Timetables[dx][day][slot1]
	}
}
