<template>
	<div class="timetable-container">
		<h2 align="center">Division {{ divisionIndex }}</h2>
		<table border="1" cellspacing="0" cellpadding="4" class="tabela">
			<thead>
				<tr>
					<th>Nr</th>
					<th v-for="day, dayIndex in days" :key="dayIndex">{{ day }}</th>
				</tr>
			</thead>
			<tbody>
				<!-- Rows for each lesson starting at 1 -->
				<tr v-for="rowIndex in actualRows" :key="rowIndex">
					<td>{{ rowIndex }}</td>
					<td v-for="dayIndex in 5" :key="dayIndex" class="l">
						<!-- Check for the existence of lessons -->
						<div v-if="schedule[dayIndex - 1]?.[rowIndex - 1]">
							<div v-for="(subject, groupIndex) in schedule[dayIndex - 1][rowIndex - 1]"
								:key="groupIndex">
								<div v-if="subject.global_subject">
									<span class="p">{{ subject.global_subject }}</span>
									<span class="n">&nbsp;{{ subject.teacher }}</span>
									<span class="s">&nbsp;{{ subject.classroom }}</span>
								</div>
							</div>
						</div>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
	divisionData: {
		type: Array, // Days Array for the division
		required: true,
	},
	divisionIndex: {
		type: Number,
		required: true,
	},
});

// Weekdays
const days = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday"];

// Actual rows based on lessons
const actualRows = computed(() => {
	return Math.max(
		...props.divisionData.map((day) =>
			day ? day.filter((group) => group.some((subject) => subject.global_subject)).length : 0
		)
	);
});

// Alias for clarity
const schedule = props.divisionData;
</script>

<style scoped>
.timetable-container {
	margin: 20px auto;
	border: 1px solid #ddd;
	padding: 10px;
	width: 95%;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	border-radius: 5px;
}

.tabela {
	width: 100%;
	border-collapse: collapse;
}

th,
td {
	border: 1px solid #000;
	text-align: center;
	padding: 4px;
}

.l {
	text-align: left;
	vertical-align: top;
}

.p {
	font-style: italic;
}

.n,
.s {
	font-weight: bold;
	margin-left: 5px;
}

.n::before,
.s::before {
	content: "\00a0";
}
</style>
